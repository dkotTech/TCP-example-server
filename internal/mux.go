package internal

import (
	"context"
	"encoding/gob"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

// Mux handles incoming connections and manages communication
type Mux struct {
	challengeLen, difficulty int
	rnd                      Randomer
}

// NewMux creates a new Mux instance
func NewMux(challengeLen, difficulty int, rnd Randomer) *Mux {
	return &Mux{
		challengeLen: challengeLen,
		difficulty:   difficulty,
		rnd:          rnd,
	}
}

func (m *Mux) HandleConnection(ctx context.Context, conn net.Conn) error {
	enc := gob.NewEncoder(conn)
	dec := gob.NewDecoder(conn)

	err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return err
	}

	challenge, err := m.rnd.Generate(m.challengeLen)
	if err != nil {
		return err
	}

	// Send the challenge to the client
	challengeRequest := Msg{
		Kind: ChallengeKind,
		Challenge: &Challenge{
			Challenge:  challenge,
			Difficulty: m.difficulty,
		},
	}

	err = enc.Encode(&challengeRequest)
	if err != nil {
		return err
	}

	log.Info().Ctx(ctx).Msgf("challenge sent: %v", challengeRequest)

	// Receive the client's response
	var clietnResponse Msg
	err = dec.Decode(&clietnResponse)
	if err != nil {
		return err
	}

	// Verify the Proof-of-Work solution
	err = VerifyPoW(challengeRequest.Challenge.Challenge, clietnResponse.Response.Nonce, challengeRequest.Challenge.Difficulty)
	if err != nil {
		log.Info().Ctx(ctx).Msg("receive bad challenge sent")
		_ = conn.Close()
		return err
	}

	// Send a random quote to the client
	quote := GetRandomQuote()

	sendQuoteMsg := Msg{
		Kind: QuoteResponseKind,
		QuoteResponse: &QuoteResponse{
			Text: quote,
		},
	}

	log.Info().Ctx(ctx).Msgf("sent quote: %s", quote)

	err = enc.Encode(&sendQuoteMsg)
	if err != nil {
		return err
	}

	// Remove the read deadline
	err = conn.SetReadDeadline(time.Time{})
	if err != nil {
		return err
	}

	/*
		here we can handle connection to map with connection id, and do concurrency request
	*/
	defer conn.Close()

	// Handle further communication with the client
	var resp Msg
	for {
		err = dec.Decode(&resp)
		if err != nil {
			return err
		}

		switch resp.Kind {
		case ChallengeKind:
		case ClientResponseKind:
		case QuoteResponseKind:
		case EchoKind:
			log.Info().Ctx(ctx).Msgf("sent echo: %s", resp.Echo.Text)
			err = enc.Encode(&resp)
			if err != nil {
				return err
			}
		}
	}
}
