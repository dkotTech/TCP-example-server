package internal

import (
	"crypto/rand"
)

type Randomer interface {
	Generate(length int) ([]byte, error)
}

type RandomChallenge struct{}

func NewRandomChallenge() *RandomChallenge {
	return &RandomChallenge{}
}

func (r *RandomChallenge) Generate(length int) ([]byte, error) {
	buf := make([]byte, length)

	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
