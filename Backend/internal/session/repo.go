package session

import (
	"Backend/domain"
	"context"
	"crypto/sha256"
	"encoding/json"
	"github.com/valkey-io/valkey-go"
	"time"
)

type repo interface {
	insert(*domain.User, *Session) error
	get(string) (string, error)
}

type Model struct {
	vk valkey.Client
}

func NewRepo(vk valkey.Client) *Model {
	return &Model{vk: vk}
}

func (m *Model) insert(user *domain.User, session *Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userStr, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return m.vk.Do(ctx, m.vk.B().Set().Key(string(session.Hash)).Value(string(userStr)).Ex(session.Expiry).Build()).Error()
}

func (m *Model) get(tokenPlaintext string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	return m.vk.Do(ctx, m.vk.B().Get().Key(string(tokenHash[:])).Build()).ToString()
}
