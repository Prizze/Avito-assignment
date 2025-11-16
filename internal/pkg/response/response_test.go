package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pr-reviewer/internal/api"
	"pr-reviewer/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendResponse(t *testing.T) {
	rec := httptest.NewRecorder()

	data := map[string]string{"foo": "bar"}
	SendResponse(rec, http.StatusOK, data)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "bar", resp["foo"])
}

func TestSendErrorResponse(t *testing.T) {
	rec := httptest.NewRecorder()

	SendErrorResponse(rec, api.NOTFOUND, http.StatusNotFound)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusNotFound, rec.Code)

	var resp api.ErrorResponse
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, api.NOTFOUND, resp.Error.Code)

	// Проверка, что сообщение соответствует domain.Messages
	expectedMsg, ok := domain.Messages[api.NOTFOUND]
	if !ok {
		expectedMsg = domain.UnknownError
	}
	assert.Equal(t, expectedMsg, resp.Error.Message)
}
