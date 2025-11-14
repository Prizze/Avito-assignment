package response

import (
	"encoding/json"
	"net/http"
	"pr-reviewer/internal/api"
	"pr-reviewer/internal/domain"
)

func SendResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(data)
}

func SendErrorResponse(w http.ResponseWriter, code api.ErrorResponseErrorCode, statusCode int) {
	var errResp api.ErrorResponse

	errResp.Error.Code = code

	msg, ok := domain.Messages[code]
	if !ok {
		msg = domain.UnknownError
	}
	errResp.Error.Message = msg

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(errResp)
}


