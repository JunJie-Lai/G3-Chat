package config

import (
	"github.com/valkey-io/valkey-go"
	"os"
)

func NewValkeyDB() (valkey.Client, error) {
	return valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{os.Getenv("VALKEY_URL")},
	})
}
