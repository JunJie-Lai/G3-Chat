package responses

import (
	"Backend/utils"
	"fmt"
	"log/slog"
	"net/http"
)

type ErrorResponses struct {
	logger *slog.Logger
	utils  *utils.Utils
}

func NewErrorResponses(logger *slog.Logger, utils *utils.Utils) *ErrorResponses {
	return &ErrorResponses{
		logger: logger,
		utils:  utils,
	}
}

func (er *ErrorResponses) logError(r *http.Request, err error) {
	er.logger.Error(err.Error(), "Method", r.Method, "URL", r.URL.RequestURI())
}

func (er *ErrorResponses) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := utils.Envelope{
		"error": message,
	}

	if err := er.utils.WriteJSON(w, status, env, nil); err != nil {
		er.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (er *ErrorResponses) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	er.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	er.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (er *ErrorResponses) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	er.errorResponse(w, r, http.StatusNotFound, message)
}

func (er *ErrorResponses) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	er.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (er *ErrorResponses) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	er.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (er *ErrorResponses) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	er.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (er *ErrorResponses) InvalidStateTokenResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid state token"
	er.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (er *ErrorResponses) FailedCodeExchangeResponse(w http.ResponseWriter, r *http.Request) {
	message := "failed to exchange the code for the request"
	er.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (er *ErrorResponses) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	er.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (er *ErrorResponses) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	er.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (er *ErrorResponses) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	er.errorResponse(w, r, http.StatusTooManyRequests, message)
}
