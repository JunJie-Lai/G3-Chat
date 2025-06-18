package middleware

import (
	"Backend/domain"
	"Backend/internal/session"
	"Backend/responses"
	"Backend/userContext"
	"Backend/utils"
	"Backend/validator"
	"errors"
	"fmt"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Middleware struct {
	er      *responses.ErrorResponses
	util    *utils.Utils
	session session.IService
}

func NewMiddleware(er *responses.ErrorResponses, util *utils.Utils, session session.IService) *Middleware {
	return &Middleware{
		er:      er,
		util:    util,
		session: session,
	}
}

func (m *Middleware) RecoverPanic(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				m.er.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func (m *Middleware) RateLimit(next http.Handler) http.HandlerFunc {
	type client struct {
		*rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENVIRONMENT") == "production" {
			ip := realip.FromRequest(r)

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					Limiter: rate.NewLimiter(rate.Limit(25), 100),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].Limiter.Allow() {
				mu.Unlock()
				m.er.RateLimitExceededResponse(w, r)
				return
			}

			// Very importantly, unlock the mutex before calling the next handler in the
			// chain. Notice that we DON'T use defer to unlock the mutex, as that would mean
			// that the mutex isn't unlocked until all the handlers downstream of this
			// middleware have also returned.
			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	}
}

func (m *Middleware) Authenticate(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = userContext.ContextSetUser(r, domain.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			m.er.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if session.ValidateTokenPlaintext(v, token); !v.Valid() {
			m.er.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := m.session.CheckSession(token)
		if err != nil {
			switch {
			case errors.Is(err, utils.ErrRecordNotFound):
				m.er.InvalidAuthenticationTokenResponse(w, r)
			default:
				m.er.ServerErrorResponse(w, r, err)
			}
			return
		}

		r = userContext.ContextSetUser(r, user)
		next.ServeHTTP(w, r)
	}
}

func (m *Middleware) RequireNonAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user := userContext.ContextGetUser(r); !user.IsAnonymous() {
			if err := m.util.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user}, nil); err != nil {
				m.er.ServerErrorResponse(w, r, err)
				return
			}
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (m *Middleware) RequireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if user := userContext.ContextGetUser(r); user.IsAnonymous() {
			m.er.AuthenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (m *Middleware) EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		trustedOriginsEnv := os.Getenv("TRUSTED_ORIGIN")
		trustedOrigins := strings.Split(trustedOriginsEnv, ",")
		if origin != "" {
			for i := range trustedOrigins {
				if origin == trustedOrigins[i] || trustedOrigins[i] == "*" {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "PATCH, OPTIONS, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Api-Key")

						w.WriteHeader(http.StatusOK)
						return
					}
					break
				}
			}
		}
		next.ServeHTTP(w, r)
	}
}
