package team

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pr-reviewer/internal/api"
	"pr-reviewer/internal/delivery/http/Team/mocks"
	"pr-reviewer/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPostTeamAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := mocks.NewMockteamUC(ctrl)
	handler := NewTeamHandler(usecase)

	t.Run("team created ok", func(t *testing.T) {
		reqBody := api.PostTeamAddJSONRequestBody{
			TeamName: "NewTeam",
			Members: []api.TeamMember{
				{IsActive: true, UserId: "u1", Username: "john"},
				{IsActive: false, UserId: "u2", Username: "garry"},
			},
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		team := &domain.Team{
			Name: "NewTeam",
			Members: []domain.TeamMember{
				{UserID: 1, Username: "john", IsActive: true},
				{UserID: 2, Username: "garry", IsActive: false},
			},
		}

		createdTeam := &domain.Team{ID: 10, Name: "NewTeam", Members: team.Members}

		usecase.EXPECT().CreateTeam(gomock.Any(), team).Return(createdTeam, nil)

		handler.PostTeamAdd(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp domain.TeamResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, "NewTeam", resp.Team.TeamName)
		assert.Len(t, resp.Team.Members, 2)
		assert.Equal(t, "u1", resp.Team.Members[0].UserId)
	})

	t.Run("bad json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBufferString(`{invalid json}`))
		rec := httptest.NewRecorder()

		handler.PostTeamAdd(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation failed: team name empty", func(t *testing.T) {
		reqBody := api.PostTeamAddJSONRequestBody{
			TeamName: "",
			Members: []api.TeamMember{
				{IsActive: true, UserId: "u1", Username: "john"},
				{IsActive: false, UserId: "u2", Username: "garry"},
			},
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.PostTeamAdd(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("validation failed: empty members", func(t *testing.T) {
		reqBody := api.PostTeamAddJSONRequestBody{
			TeamName: "NormalName",
			Members:  []api.TeamMember{},
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.PostTeamAdd(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("uc returns error", func(t *testing.T) {
		reqBody := api.PostTeamAddJSONRequestBody{
			TeamName: "TeamA",
			Members: []api.TeamMember{
				{IsActive: true, UserId: "u1", Username: "john"},
				{IsActive: false, UserId: "u2", Username: "garry"},
			},
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		apiTeam := &domain.Team{
			Name: "TeamA",
			Members: []domain.TeamMember{
				{UserID: 1, Username: "john", IsActive: true},
				{UserID: 2, Username: "garry", IsActive: false},
			},
		}

		usecase.EXPECT().CreateTeam(gomock.Any(), apiTeam).Return(nil, domain.ErrTeamExists)

		handler.PostTeamAdd(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestGetTeamGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usecase := mocks.NewMockteamUC(ctrl)
	handler := NewTeamHandler(usecase)

	t.Run("get team ok", func(t *testing.T) {
		params := api.GetTeamGetParams{
			TeamName: "MyTeam",
		}

		req := httptest.NewRequest(http.MethodGet, "/team/get", nil)
		rec := httptest.NewRecorder()

		team := &domain.Team{
			ID:   1,
			Name: "MyTeam",
			Members: []domain.TeamMember{
				{UserID: 1, Username: "john", IsActive: true},
				{UserID: 2, Username: "alex", IsActive: false},
			},
		}

		usecase.EXPECT().GetTeamByName(gomock.Any(), "MyTeam").Return(team, nil)

		handler.GetTeamGet(rec, req, params)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp api.Team
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, "MyTeam", resp.TeamName)
		assert.Len(t, resp.Members, 2)
		assert.Equal(t, "u1", resp.Members[0].UserId)
	})

	t.Run("invalid team name", func(t *testing.T) {
		params := api.GetTeamGetParams{
			TeamName: "",
		}

		req := httptest.NewRequest(http.MethodGet, "/team/get", nil)
		rec := httptest.NewRecorder()

		handler.GetTeamGet(rec, req, params)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("team not found", func(t *testing.T) {
		params := api.GetTeamGetParams{
			TeamName: "UnknownTeam",
		}

		req := httptest.NewRequest(http.MethodGet, "/team/get", nil)
		rec := httptest.NewRecorder()

		usecase.EXPECT().GetTeamByName(gomock.Any(), "UnknownTeam").Return(nil, domain.ErrTeamNotFound)

		handler.GetTeamGet(rec, req, params)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		params := api.GetTeamGetParams{
			TeamName: "MyTeam",
		}

		req := httptest.NewRequest(http.MethodGet, "/team/get", nil)
		rec := httptest.NewRecorder()

		usecase.EXPECT().GetTeamByName(gomock.Any(), "MyTeam").Return(nil, errors.New("internal error"))

		handler.GetTeamGet(rec, req, params)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
