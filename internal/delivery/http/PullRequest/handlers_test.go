package pullrequest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pr-reviewer/internal/api"
	"pr-reviewer/internal/delivery/http/PullRequest/mocks"
	"pr-reviewer/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPostPullRequestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := mocks.NewMockprUC(ctrl)
	handler := NewPRHandler(usecase)

	t.Run("pull request created ok", func(t *testing.T) {
		reqBody := api.PostPullRequestCreateJSONRequestBody{
			AuthorId:        "u123",
			PullRequestId:   "pr-1001",
			PullRequestName: "Improve coverage",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		apiPR := &domain.CreatePullRequest{
			AuthorId:      123,
			PullRequestId: 1001,
			Name:          "Improve coverage",
		}

		createdPR := &domain.PullRequest{
			ID:                1001,
			Name:              "Improve coverage",
			AuthorID:          123,
			Status:            domain.PRStatusOpen,
			AssignedReviewers: []int{1, 2},
		}

		usecase.EXPECT().CreatePullRequest(gomock.Any(), apiPR).Return(createdPR, nil)

		handler.PostPullRequestCreate(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp domain.PullRequestResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, "pr-1001", resp.PullRequest.PullRequestId)
		assert.Equal(t, "Improve coverage", resp.PullRequest.PullRequestName)
		assert.Equal(t, "u123", resp.PullRequest.AuthorId)
		assert.Len(t, resp.PullRequest.AssignedReviewers, 2)
	})

	t.Run("bad json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBufferString("{invalid"))
		rec := httptest.NewRecorder()

		handler.PostPullRequestCreate(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid author id", func(t *testing.T) {
		reqBody := api.PostPullRequestCreateJSONRequestBody{
			AuthorId:        "user123",
			PullRequestId:   "pr-1001",
			PullRequestName: "Improve coverage",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.PostPullRequestCreate(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("uc returns error", func(t *testing.T) {
		reqBody := api.PostPullRequestCreateJSONRequestBody{
			AuthorId:        "u123",
			PullRequestId:   "pr-99",
			PullRequestName: "Fix bug",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		apiPR := &domain.CreatePullRequest{
			AuthorId:      123,
			PullRequestId: 99,
			Name:          "Fix bug",
		}

		usecase.EXPECT().CreatePullRequest(gomock.Any(), apiPR).Return(nil, domain.ErrPullRequestExists)

		handler.PostPullRequestCreate(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}

func TestPostPullRequestCreateValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := mocks.NewMockprUC(ctrl)
	handler := NewPRHandler(usecase)

	tests := []struct {
		name           string
		body           api.PostPullRequestCreateJSONRequestBody
		expectedStatus int
	}{
		{
			name: "invalid author_id",
			body: api.PostPullRequestCreateJSONRequestBody{
				AuthorId:        "user123",
				PullRequestId:   "pr-1001",
				PullRequestName: "Improve coverage",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty pull_request_name",
			body: api.PostPullRequestCreateJSONRequestBody{
				AuthorId:        "u123",
				PullRequestId:   "pr-1002",
				PullRequestName: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid pull_request_id",
			body: api.PostPullRequestCreateJSONRequestBody{
				AuthorId:        "u123",
				PullRequestId:   "prwrong-1002",
				PullRequestName: "PR name",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(b))
			rec := httptest.NewRecorder()

			handler.PostPullRequestCreate(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}

}

func TestPostPullRequestMerge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := mocks.NewMockprUC(ctrl)
	handler := NewPRHandler(usecase)

	t.Run("merge success ok", func(t *testing.T) {
		bodyStruct := api.PostPullRequestMergeJSONRequestBody{
			PullRequestId: "pr-42",
		}
		body, _ := json.Marshal(bodyStruct)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		mergedPR := &domain.PullRequest{
			ID:       42,
			Name:     "Some PR",
			AuthorID: 1,
			Status:   domain.PRStatusMerged,
		}

		usecase.EXPECT().MergePullRequest(gomock.Any(), 42).Return(mergedPR, nil)

		handler.PostPullRequestMerge(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp domain.PullRequestResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "pr-42", resp.PullRequest.PullRequestId)
		assert.Equal(t, "Some PR", resp.PullRequest.PullRequestName)
		assert.Equal(t, api.PullRequestStatusMERGED, resp.PullRequest.Status)
	})

	t.Run("bad json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewBufferString("{invalid"))
		rec := httptest.NewRecorder()

		handler.PostPullRequestMerge(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid pull_request_id", func(t *testing.T) {
		bodyStruct := api.PostPullRequestMergeJSONRequestBody{
			PullRequestId: "wrong-format",
		}
		body, _ := json.Marshal(bodyStruct)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.PostPullRequestMerge(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("uc returns error: pr not found", func(t *testing.T) {
		bodyStruct := api.PostPullRequestMergeJSONRequestBody{
			PullRequestId: "pr-99",
		}
		body, _ := json.Marshal(bodyStruct)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		usecase.EXPECT().MergePullRequest(gomock.Any(), 99).Return(nil, domain.ErrPullRequestNotFound)

		handler.PostPullRequestMerge(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestPostPullRequestReassign(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := mocks.NewMockprUC(ctrl)
	handler := NewPRHandler(usecase)

	t.Run("reassigned reviewer ok", func(t *testing.T) {
		bodyStruct := api.PostPullRequestReassignJSONRequestBody{
			PullRequestId: "pr-42",
			OldUserId:     "u123",
		}
		body, _ := json.Marshal(bodyStruct)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		reasDomain := &domain.ReassingReviewer{
			PullRequestID: 42,
			UserID:        123,
		}

		reassignedPR := &domain.PullRequest{
			ID:                42,
			Name:              "Fix bug",
			AuthorID:          999,
			AssignedReviewers: []int{1, 456},
		}

		replacedBy := 456 // новый ревьювер

		usecase.EXPECT().ReassignReviewer(gomock.Any(), reasDomain).Return(reassignedPR, replacedBy, nil)

		handler.PostPullRequestReassign(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp domain.ReassignResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "pr-42", resp.PullRequest.PullRequestId)
		assert.Equal(t, "u456", resp.ReplacedBy)
	})

	t.Run("bad json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewBufferString("{invalid"))
		rec := httptest.NewRecorder()

		handler.PostPullRequestReassign(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("pr already merged", func(t *testing.T) {
		bodyStruct := api.PostPullRequestReassignJSONRequestBody{
			PullRequestId: "pr-42",
			OldUserId:     "u123",
		}
		body, _ := json.Marshal(bodyStruct)

		req := httptest.NewRequest(http.MethodPost, "/pr/reassign", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		reasDomain := &domain.ReassingReviewer{
			PullRequestID: 42,
			UserID:        123,
		}

		usecase.EXPECT().ReassignReviewer(gomock.Any(), reasDomain).Return(nil, 0, domain.ErrPullRequestIsMerged)

		handler.PostPullRequestReassign(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("old user not in team", func(t *testing.T) {
		bodyStruct := api.PostPullRequestReassignJSONRequestBody{
			PullRequestId: "pr-42",
			OldUserId:     "u999",
		}
		body, _ := json.Marshal(bodyStruct)

		req := httptest.NewRequest(http.MethodPost, "/pr/reassign", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		reasDomain := &domain.ReassingReviewer{
			PullRequestID: 42,
			UserID:        999,
		}

		usecase.EXPECT().ReassignReviewer(gomock.Any(), reasDomain).Return(nil, 0, domain.ErrNotAssigned)

		handler.PostPullRequestReassign(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("no available reviewers", func(t *testing.T) {
		bodyStruct := api.PostPullRequestReassignJSONRequestBody{
			PullRequestId: "pr-42",
			OldUserId:     "u123",
		}
		body, _ := json.Marshal(bodyStruct)

		req := httptest.NewRequest(http.MethodPost, "/pr/reassign", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		reasDomain := &domain.ReassingReviewer{
			PullRequestID: 42,
			UserID:        123,
		}

		usecase.EXPECT().ReassignReviewer(gomock.Any(), reasDomain).Return(nil, 0, domain.ErrNoAvailableCandidats)

		handler.PostPullRequestReassign(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}
