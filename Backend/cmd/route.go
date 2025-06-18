package main

import (
	"Backend/internal/chat"
	"Backend/internal/session"
	"Backend/internal/user"
	"Backend/middleware"
	"net/http"
)

func (app *application) route() http.Handler {
	mux := http.NewServeMux()

	path := [3]string{"/v1/auth/google/login", "/v1/auth/google/callback", "/v1/auth/google/revoke"}
	for _, route := range path {
		mux.HandleFunc(route, app.responses.MethodNotAllowedResponse)
	}

	sessionRepo := session.NewRepo(app.vkDB)
	sessionService := session.NewService(sessionRepo)

	middle := middleware.NewMiddleware(app.responses, app.util, sessionService)

	userRepo := user.NewRepo(app.db, app.vkDB)
	userService := user.NewService(userRepo, app.oauth)
	userHandler := user.NewHandler(userService, sessionService, app.responses, app.util)
	userHandler.RegisterRoutes(mux, middle)

	chatRepo := chat.NewRepo(app.db, app.vkDB)
	chatService := chat.NewService(chatRepo, app.multiLLM)
	chatHandler := chat.NewHandler(chatService, app.responses, app.util)
	chatHandler.RegisterRoutes(mux, middle)

	return middle.RecoverPanic(middle.EnableCORS(middle.RateLimit(middle.Authenticate(mux))))
}
