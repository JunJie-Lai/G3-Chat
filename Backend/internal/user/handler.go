package user

import (
	"Backend/domain"
	"Backend/internal/session"
	"Backend/middleware"
	"Backend/responses"
	"Backend/userContext"
	"Backend/utils"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Handler struct {
	userService    IService
	sessionService session.IService
	er             *responses.ErrorResponses
	utils          *utils.Utils
}

func NewHandler(userService IService, sessionService session.IService, er *responses.ErrorResponses, utils *utils.Utils) *Handler {
	return &Handler{
		userService:    userService,
		sessionService: sessionService,
		er:             er,
		utils:          utils,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, middle *middleware.Middleware) {
	mux.HandleFunc("GET /v1/auth/google/login", middle.RequireNonAuthenticatedUser(h.handleGoogleLogin))
	mux.HandleFunc("GET /v1/auth/google/callback", middle.RequireNonAuthenticatedUser(h.handleGoogleCallback))
	mux.HandleFunc("GET /user", middle.RequireAuthenticatedUser(h.handlerUser))
	mux.HandleFunc("DELETE /v1/auth/google/revoke", middle.RequireAuthenticatedUser(h.handleGoogleRevoke))
}

func (h *Handler) handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	authURL, err := h.userService.getAuthURL()
	if err != nil {
		h.er.ServerErrorResponse(w, r, err)
		return
	}

	if err := h.utils.WriteJSON(w, http.StatusOK, utils.Envelope{"auth_url": authURL}, nil); err != nil {
		h.er.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	receivedStateToken := r.URL.Query().Get("state")
	validStateToken, err := h.userService.checkStateToken(receivedStateToken)
	if err != nil {
		h.er.ServerErrorResponse(w, r, err)
		return
	}

	if !validStateToken {
		h.er.InvalidStateTokenResponse(w, r)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := h.userService.getExchangeToken(r.Context(), code)
	if err != nil {
		h.er.FailedCodeExchangeResponse(w, r)
		return
	}

	client := h.userService.getOAuthClient(r.Context(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		h.er.ServerErrorResponse(w, r, err)
		return
	}

	r.Body = response.Body
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			h.er.ServerErrorResponse(w, r, err)
			return
		}
	}(response.Body)

	var userInfo struct {
		ID            string `json:"sub"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		GivenName     string `json:"given_name"`
		EmailVerified bool   `json:"email_verified"`
		FamilyName    string `json:"family_name"`
	}
	if err := h.utils.ReadJSON(w, r, &userInfo); err != nil {
		h.er.BadRequestResponse(w, r, err)
		return
	}

	user := &domain.User{
		ID:           userInfo.ID,
		Name:         userInfo.Name,
		Email:        userInfo.Email,
		Picture:      userInfo.Picture,
		RefreshToken: token.RefreshToken,
	}
	if err := h.userService.loginUser(user); err != nil {
		h.er.ServerErrorResponse(w, r, err)
		return
	}

	sessionToken, err := h.sessionService.NewSessionToken(user, 30*24*time.Hour)
	if err != nil {
		h.er.ServerErrorResponse(w, r, err)
		return
	}

	redirectURL := fmt.Sprintf("http://localhost:3000/auth?token=%s", sessionToken.Plaintext)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)

	//if err := h.utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user, "session_token": sessionToken}, nil); err != nil {
	//	h.er.ServerErrorResponse(w, r, err)
	//}
}

func (h *Handler) handlerUser(w http.ResponseWriter, r *http.Request) {
	user := userContext.ContextGetUser(r)
	if err := h.utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user}, nil); err != nil {
		h.er.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) handleGoogleRevoke(w http.ResponseWriter, r *http.Request) {
	user := userContext.ContextGetUser(r)

	token, err := h.userService.getRefreshToken(user.ID)

	if err != nil {
		h.er.ServerErrorResponse(w, r, err)
		return
	}

	response, err := h.userService.getOAuthClient(r.Context(), token).
		PostForm("https://oauth2.googleapis.com/revoke", url.Values{"token": {token.RefreshToken}})
	if err != nil {
		h.er.ServerErrorResponse(w, r, err)
		return
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			return
		}
	}(response.Body)

	if response.StatusCode == http.StatusOK {
		if err := h.userService.deleteUser(user.ID); err != nil {
			switch {
			case errors.Is(err, utils.ErrRecordNotFound):
				h.er.NotFoundResponse(w, r)
			default:
				h.er.ServerErrorResponse(w, r, err)
			}
			return
		}

		if err := h.utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Account Deletion Successful"}, nil); err != nil {
			h.er.ServerErrorResponse(w, r, err)
		}
		return
	}

	h.er.BadRequestResponse(w, r, errors.New("account deletion failed"))
}
