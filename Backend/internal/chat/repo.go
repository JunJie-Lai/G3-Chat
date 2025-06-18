package chat

import (
	"context"
	"database/sql"
	"github.com/tmc/langchaingo/llms"
	"github.com/valkey-io/valkey-go"
	"time"
)

type Chat struct {
	ID      int32     `json:"id"`
	Title   string    `json:"title,omitempty"`
	Message []Message `json:"message,omitempty"`
}

type Message struct {
	Text string `json:"text"`
}

type repo interface {
	getMessageHistory(int32) ([]llms.MessageContent, error)
	insertLatestMessage(int32, string, string) error
	insertTitle(string, string) (int32, string, error)
	getTitles(string) ([]Chat, error)
	deleteChat(string, int32) error
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

func (m *Model) getMessageHistory(chatID int32) ([]llms.MessageContent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, "SELECT text, type FROM message WHERE title_id = $1 ORDER BY timestamp", chatID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			return
		}
	}(rows)

	var results []llms.MessageContent
	for rows.Next() {
		var text string
		var messageType string
		if err := rows.Scan(&text, &messageType); err != nil {
			return nil, err
		}
		results = append(results, llms.TextParts(llms.ChatMessageType(messageType), text))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *Model) insertLatestMessage(chatID int32, prompt string, text string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := m.db.ExecContext(ctx, "INSERT INTO message (title_id, text, type) VALUES ($1, $2, $3), ($1, $4, $5)", chatID, prompt, llms.ChatMessageTypeHuman, text, llms.ChatMessageTypeAI); err != nil {
		return err
	}
	return nil
}

func (m *Model) insertTitle(userID string, title string) (int32, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var chatID int32
	if err := m.db.QueryRowContext(ctx, "INSERT INTO title (user_id, title) VALUES ($1, $2) RETURNING id", userID, title).Scan(&chatID); err != nil {
		return 0, "", err
	}

	return chatID, title, nil
}

func (m *Model) getTitles(userID string) ([]Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.db.QueryContext(ctx, "SELECT title.id, title FROM title LEFT JOIN message ON title.id = title_id WHERE user_id = $1 GROUP BY title.id ORDER BY MAX(timestamp) DESC", userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			return
		}
	}(rows)

	var chats []Chat
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ID, &chat.Title); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}

func (m *Model) deleteChat(userID string, chatID int32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := m.db.ExecContext(ctx, "DELETE FROM title WHERE id = $1 AND user_id = $2", chatID, userID); err != nil {
		return err
	}

	return nil
}
