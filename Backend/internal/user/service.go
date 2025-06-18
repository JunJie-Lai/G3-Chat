package user

import (
	"Backend/domain"
	"context"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/oauth2"
	"net/http"
)

type IService interface {
	loginUser(*domain.User) error
	getRefreshToken(string) (*oauth2.Token, error)
	deleteUser(string) error
	checkStateToken(string) (bool, error)

	getAuthURL() (string, error)
	getExchangeToken(context.Context, string) (*oauth2.Token, error)
	getOAuthClient(context.Context, *oauth2.Token) *http.Client
}

type service struct {
	userRepo repo
	oauth    *oauth2.Config
}

func NewService(userRepo repo, oauth *oauth2.Config) IService {
	return &service{
		userRepo: userRepo,
		oauth:    oauth,
	}
}

func (s *service) loginUser(user *domain.User) error {
	return s.userRepo.upsert(user)
}

func (s *service) getRefreshToken(userID string) (*oauth2.Token, error) {
	refreshToken, err := s.userRepo.getRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	return token, nil
}

func (s *service) deleteUser(userID string) error {
	return s.userRepo.delete(userID)
}

func (s *service) checkStateToken(receivedStateToken string) (bool, error) {
	stateToken, err := s.userRepo.getStateToken(receivedStateToken)
	if err != nil {
		return false, err
	}

	return stateToken, nil
}

func generateStateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *service) getAuthURL() (string, error) {
	stateToken, err := generateStateToken()
	if err != nil {
		return "", err
	}

	if err := s.userRepo.setStateToken(stateToken); err != nil {
		return "", err
	}
	return s.oauth.AuthCodeURL(stateToken, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("include_granted_scopes", "true")), nil
}

func (s *service) getExchangeToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.oauth.Exchange(ctx, code)
}

func (s *service) getOAuthClient(ctx context.Context, token *oauth2.Token) *http.Client {
	return s.oauth.Client(ctx, token)
}
