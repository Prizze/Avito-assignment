package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pr-reviewer/internal/api"
	"pr-reviewer/internal/delivery/http/User/mocks"
	"pr-reviewer/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUsersGetReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := mocks.NewMockuserUC(ctrl)
	handler := NewUserHandler(usecase)

	t.Run("ok get reviews", func(t *testing.T) {
		params := api.GetUsersGetReviewParams{UserId: "u1"}
		prs := []domain.PullRequest{
			{ID: 1, Name: "Fix Bug"},
			{ID: 2, Name: "Add Feature"},
		}

		usecase.EXPECT().GetUserPullRequests(gomock.Any(), 1).Return(prs, nil)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/getReview", nil)

		handler.GetUsersGetReview(rec, req, params)

		assert.Equal(t, rec.Code, http.StatusOK)
		var resp domain.UserReviews
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "u1", resp.UserID)
		assert.Len(t, resp.PullRequests, 2)
	})

	t.Run("bad userId: validation error", func(t *testing.T) {
		params := api.GetUsersGetReviewParams{
			UserId: "INVALID",
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/getReview", nil)

		handler.GetUsersGetReview(rec, req, params)

		assert.Equal(t, rec.Code, http.StatusBadRequest)
	})

	t.Run("uc returns error", func(t *testing.T) {
		params := api.GetUsersGetReviewParams{
			UserId: "u1234",
		}

		usecase.EXPECT().GetUserPullRequests(gomock.Any(), 1234).Return(nil, domain.ErrUserNotFound)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/u123/reviews", nil)

		handler.GetUsersGetReview(rec, req, params)

		assert.Equal(t, rec.Code, http.StatusNotFound)
	})
}

func TestPostUsersSetIsActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := mocks.NewMockuserUC(ctrl)
	handler := NewUserHandler(usecase)

	t.Run("set is_active ok", func(t *testing.T) {
		reqBody := api.PostUsersSetIsActiveJSONRequestBody{
			UserId:   "u1",
			IsActive: false,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		setActive := &domain.SetUserIsActive{
			ID:       1,
			IsActive: false,
		}

		expectedUser := &domain.User{
			ID:       1,
			Username: "john",
			IsActive: false,
		}

		usecase.EXPECT().SetUserIsActive(gomock.Any(), setActive).Return(expectedUser, nil)

		handler.PostUsersSetIsActive(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp domain.UserResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, "john", resp.User.Username)
		assert.Equal(t, false, resp.User.IsActive)
	})

	t.Run("bad json", func(t *testing.T) {
		body := `{invalid-json}`
		req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		handler.PostUsersSetIsActive(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid user_id", func(t *testing.T) {
		reqBody := api.PostUsersSetIsActiveJSONRequestBody{
			UserId:   "user1",
			IsActive: false,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.PostUsersSetIsActive(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("uc returns error: user not found", func(t *testing.T) {
		reqBody := api.PostUsersSetIsActiveJSONRequestBody{
			UserId:   "u123",
			IsActive: false,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		setActive := &domain.SetUserIsActive{
			ID:       123,
			IsActive: false,
		}

		usecase.EXPECT().SetUserIsActive(gomock.Any(), setActive).Return(nil, domain.ErrUserNotFound)

		handler.PostUsersSetIsActive(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("uc returns error: internal", func(t *testing.T) {
		reqBody := api.PostUsersSetIsActiveJSONRequestBody{
			UserId:   "u1",
			IsActive: false,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		setActive := &domain.SetUserIsActive{
			ID:       1,
			IsActive: false,
		}

		usecase.EXPECT().SetUserIsActive(gomock.Any(), setActive).Return(nil, errors.New("internal error"))

		handler.PostUsersSetIsActive(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

}
