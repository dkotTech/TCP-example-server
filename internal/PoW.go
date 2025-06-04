package internal

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strings"
	"tcp_test_work/errors"

	"github.com/rs/zerolog/log"
)

func VerifyPoW(challenge, nonce []byte, difficulty int) error {
	hash := sha256.Sum256(append(challenge, nonce...))
	log.Info().Msgf("hash: %s", fmt.Sprintf("%x", hash))
	prefix := strings.Repeat("0", difficulty)
	if !strings.HasPrefix(fmt.Sprintf("%x", hash), prefix) {
		return errors.ErrConnectionFailed
	}

	return nil
}

func SolvePoW(challenge []byte, difficulty int) []byte {
	targetPrefix := strings.Repeat("0", difficulty)
	nonce := make([]byte, 8)
	for i := 0; ; i++ {
		binary.BigEndian.PutUint64(nonce, uint64(i))
		hash := sha256.Sum256(append(challenge, nonce...))
		hashStr := fmt.Sprintf("%x", hash)
		if strings.HasPrefix(hashStr, targetPrefix) {
			log.Info().Msgf("hash: %s", hashStr)
			return nonce
		}
	}
}
