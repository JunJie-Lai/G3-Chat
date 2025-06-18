package user

import (
	"Backend/domain"
	"Backend/utils"
	"context"
	"database/sql"
	"github.com/valkey-io/valkey-go"
	"time"
)

type repo interface {
	upsert(*domain.User) error
	delete(string) error
	getRefreshToken(string) (string, error)
	getStateToken(string) (bool, error)
	setStateToken(string) error
}

type Model struct {
	db *sql.DB
	vk valkey.Client
}

func NewRepo(db *sql.DB, vk valkey.Client) *Model {
	return &Model{
		db: db,
		vk: vk,
	}
}

func (m *Model) upsert(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx,
		"INSERT INTO users (id, name, email, picture, refresh_token) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET name = excluded.name, picture = excluded.picture, refresh_token = excluded.refresh_token WHERE excluded.refresh_token <> ''",
		user.ID, user.Name, user.Email, user.Picture, user.RefreshToken)
	return err
}

func (m *Model) delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return utils.ErrRecordNotFound
	}
	return nil
}

func (m *Model) getRefreshToken(userID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var token string
	if err := m.db.QueryRowContext(ctx, "SELECT refresh_token FROM users WHERE id = $1", userID).Scan(&token); err != nil {
		return "", err
	}
	return token, nil
}

func (m *Model) getStateToken(stateToken string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.vk.Do(ctx, m.vk.B().Exists().Key(stateToken).Build()).AsBool()
}

func (m *Model) setStateToken(stateToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.vk.Do(ctx, m.vk.B().Setex().Key(stateToken).Seconds(300).Value("").Build()).Error()
}
