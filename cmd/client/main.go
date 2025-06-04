package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"tcp_test_work/internal"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	address = flag.String("address", ":8080", "connection address")
)

var (
	// Channel used to signal when an echo response is received
	waitForEcho = make(chan bool)
)

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", *address)
	if err != nil {
		log.Fatal().Err(err).Msg("net.Dial")
	}
	defer conn.Close()

	// Initialize encoders and decoders for communication
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	// Function to read input from the console and send it to the server
	consoleReaderRunner := sync.OnceFunc(func() {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Println("Start reading lines. Press Ctrl+C to exit.")
	loop:
		for {
			resp := internal.Msg{Kind: internal.EchoKind, Echo: &internal.Echo{}}

			fmt.Print("> ")
			if scanner.Scan() {
				line := scanner.Text()

				resp.Echo.Text = line

				err = enc.Encode(&resp)
				if err != nil {
					log.Fatal().Err(err).Msg("enc.Encode")
				}

				select {
				case <-time.After(time.Second * 10):
					log.Info().Msg("> echo: timeouted")
					continue loop
				case <-waitForEcho:
				}
			} else {
				break loop
			}
		}
	})

	var msg internal.Msg

	for {
		err = dec.Decode(&msg)
		if err != nil {
			log.Fatal().Err(err).Msg("dec.Decode")
			return
		}

		switch msg.Kind {
		case internal.ChallengeKind:
			// Solve the Proof-of-Work challenge
			log.Printf("challenge: %s, difficulty: %d", msg.Challenge.Challenge, msg.Challenge.Difficulty)

			nonce := internal.SolvePoW(msg.Challenge.Challenge, msg.Challenge.Difficulty)
			log.Printf("found: %s", nonce)

			// Send the solution back to the server
			resp := internal.Msg{Response: &internal.Response{Nonce: nonce}}
			err = enc.Encode(&resp)
			if err != nil {
				log.Fatal().Err(err).Msg("resp.Encode")
			}
		case internal.QuoteResponseKind:
			// Print the quote received from the server
			fmt.Println("response:", msg.QuoteResponse.Text)

			// Start reading console input
			go func() {
				consoleReaderRunner()
			}()
		case internal.EchoKind:
			// Print the echo response and signal completion
			fmt.Printf("< %s\n", msg.Echo.Text)
			waitForEcho <- true
		default:
			continue
		}
	}
}
