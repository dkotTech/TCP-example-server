package internal

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"

	"tcp_test_work/errors"
)

func TestVerifyPoW(t *testing.T) {
	challenge := []byte("testChallenge")
	difficulty := 4
	nonce := SolvePoW(challenge, difficulty)

	err := VerifyPoW(challenge, nonce, difficulty)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test with incorrect nonce
	invalidNonce := []byte("invalidNonce")
	err = VerifyPoW(challenge, invalidNonce, difficulty)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != errors.ErrConnectionFailed {
		t.Fatalf("expected ErrConnectionFailed, got %v", err)
	}
}

func TestSolvePoW(t *testing.T) {
	challenge := []byte("testChallenge")
	difficulty := 4
	nonce := SolvePoW(challenge, difficulty)

	hash := sha256.Sum256(append(challenge, nonce...))
	hashStr := fmt.Sprintf("%x", hash)
	expectedPrefix := strings.Repeat("0", difficulty)

	if !strings.HasPrefix(hashStr, expectedPrefix) {
		t.Fatalf("expected hash to start with %s, got %s", expectedPrefix, hashStr)
	}
}
