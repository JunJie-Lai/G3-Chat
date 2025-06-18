package chat

import (
	"Backend/middleware"
	"Backend/responses"
	"Backend/userContext"
	"Backend/utils"
	"net/http"
)

type Handler struct {
	chatService IService
	er          *responses.ErrorResponses
	utils       *utils.Utils
}

func NewHandler(chatService IService, er *responses.ErrorResponses, utils *utils.Utils) *Handler {
	return &Handler{
		chatService: chatService,
		er:          er,
		utils:       utils,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, middle *middleware.Middleware) {
	mux.HandleFunc("GET /v1/chat", middle.RequireAuthenticatedUser(h.getTitlesHandler))
	mux.HandleFunc("GET /v1/chat/{id}", middle.RequireAuthenticatedUser(h.getCurrentChatHistoryHandler))
	mux.HandleFunc("POST /v1/chat", h.sendMessageHandler)
	mux.HandleFunc("DELETE /v1/chat", middle.RequireAuthenticatedUser(h.deleteChatHandler))
}

func (h *Handler) getTitlesHandler(w http.ResponseWriter, r *http.Request) {
	user := userContext.ContextGetUser(r)

	titles, err := h.chatService.getTitles(user.ID)
	if err != nil {
		h.er.ServerErrorResponse(w, r, err)
		return
	}

	if err := h.utils.WriteJSON(w, http.StatusOK, utils.Envelope{"titles": titles}, nil); err != nil {
		h.er.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) getCurrentChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	chatID, err := h.utils.ReadIDParam(r)
	if err != nil {
		h.er.ServerErrorResponse(w, r, err)
		return
	}

	chatHistory, err := h.chatService.getChatHistory(int32(chatID))
	if err != nil {
		h.er.ServerErrorResponse(w, r, err)
		return
	}

	if err := h.utils.WriteJSON(w, http.StatusOK, utils.Envelope{"chatHistory": chatHistory}, nil); err != nil {
		h.er.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Api-Key")
	var input struct {
		ID        int32  `json:"id"` //0 for anon, use -1 once for title generation for anon user
		ModelType string `json:"model_type"`
		Model     string `json:"model"` //optional
		Prompt    string `json:"prompt"`
	}

	if err := h.utils.ReadJSON(w, r, &input); err != nil {
		h.er.ServerErrorResponse(w, r, err)
		return
	}

	if validInput, err := h.chatService.checkInput(input.ModelType, input.Model, input.Prompt); !validInput {
		h.er.FailedValidationResponse(w, r, err)
		return
	}

	user := userContext.ContextGetUser(r)
	var chat Chat
	chat.ID = input.ID
	if input.ID == -1 || (input.ID < 1 && !user.IsAnonymous()) {
		chatID, title, err := h.chatService.generateTitle(user.ID, input.ModelType, input.Model, apiKey, input.Prompt)
		if err != nil {
			h.er.ServerErrorResponse(w, r, err)
			return
		}
		chat.ID = chatID
		chat.Title = title
	}

	text, err := h.chatService.processOutput(user.ID, chat.ID, input.ModelType, input.Model, apiKey, input.Prompt)
	if err != nil {
		h.er.ServerErrorResponse(w, r, err)
		return
	}
	chat.Message = []Message{
		{
			Text: text,
		},
	}

	if err := h.utils.WriteJSON(w, http.StatusOK, utils.Envelope{"chat": chat}, nil); err != nil {
		h.er.ServerErrorResponse(w, r, err)
	}
}

func (h *Handler) deleteChatHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID int32 `json:"id"`
	}

	user := userContext.ContextGetUser(r)
	if err := h.chatService.deleteChat(user.ID, input.ID); err != nil {
		h.er.ServerErrorResponse(w, r, err)
		return
	}

	if err := h.utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Chat Deletion Successful!"}, nil); err != nil {
		h.er.ServerErrorResponse(w, r, err)
	}
}
