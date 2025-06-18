package session

import "time"

type Session struct {
	UserID    string        `json:"-"`
	Plaintext string        `json:"token"`
	Hash      []byte        `json:"-"`
	Expiry    time.Duration `json:"expiry"`
}
