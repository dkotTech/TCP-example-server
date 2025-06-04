package main

import (
	"context"
	"flag"
	"net"
	"tcp_test_work/internal"

	"github.com/rs/zerolog/log"
)

var (
	challengeLen = flag.Int("challenge", 10, "a challenge length")
	address      = flag.String("address", ":8080", "the address to listen on")
	difficulty   = flag.Int("difficulty", 4, "PoW difficulty")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	// Initialize the random challenge generator
	rndService := internal.NewRandomChallenge()

	listener, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatal().Ctx(ctx).Err(err).Send()
	}
	defer listener.Close()

	log.Info().Ctx(ctx).Msgf("server address: %s", *address)

	mux := internal.NewMux(*challengeLen, *difficulty, rndService)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Warn().Ctx(ctx).Err(err).Msg("Accept")
			continue
		}

		// Handle the connection in a separate goroutine
		go func(conn net.Conn) {
			connCtx, can := context.WithCancel(ctx)
			defer can()

			if err := mux.HandleConnection(connCtx, conn); err != nil {
				log.Error().Ctx(ctx).Err(err).Msg("HandleConnection")
				return
			}
		}(conn)
	}
}
