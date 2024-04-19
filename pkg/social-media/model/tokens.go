package model

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string
	Hash      []byte
	MemberID  int64
	Expiry    time.Time
	Scope     string
}

func generateToken(memberID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		MemberID: memberID,
		Expiry:   time.Now().Add(ttl),
		Scope:    scope,
	}
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}