package session

import (
	"Backend/domain"
	"Backend/utils"
	"Backend/validator"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"time"
)

type IService interface {
	NewSessionToken(*domain.User, time.Duration) (*Session, error)
	CheckSession(string) (*domain.User, error)
}

type service struct {
	repo repo
}

func NewService(repo repo) IService {
	return &service{
		repo: repo,
	}
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

func generateToken(userID string, ttl time.Duration) (*Session, error) {
	// ttl: time-to-live
	sessionToken := &Session{
		UserID: userID,
		Expiry: ttl,
	}

	randomBytes := make([]byte, 16)

	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}

	sessionToken.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(sessionToken.Plaintext))
	sessionToken.Hash = hash[:]

	return sessionToken, nil
}

func (s *service) NewSessionToken(user *domain.User, ttl time.Duration) (*Session, error) {
	token, err := generateToken(user.ID, ttl)
	if err != nil {
		return nil, err
	}

	if err := s.repo.insert(user, token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *service) CheckSession(tokenPlaintext string) (*domain.User, error) {
	userStr, err := s.repo.get(tokenPlaintext)
	if err != nil {
		return nil, err
	}

	if userStr == "" {
		return nil, utils.ErrRecordNotFound
	}

	var user *domain.User
	if err := json.Unmarshal([]byte(userStr), &user); err != nil {
		return nil, err
	}

	return user, nil
}
