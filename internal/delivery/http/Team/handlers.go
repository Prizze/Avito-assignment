package team

import (
	"encoding/json"
	"errors"
	"net/http"
	"pr-reviewer/internal/api"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/pkg/response"
	"pr-reviewer/internal/pkg/validation"
)

type TeamHandler struct {
	uc teamUC
}

func NewTeamHandler(uc teamUC) *TeamHandler {
	return &TeamHandler{
		uc: uc,
	}
}

func (h *TeamHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	var req api.PostTeamAddJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, api.BADREQUEST, http.StatusBadRequest)
		return
	}

	if err := validation.ValidateTeam(req); err != nil {
		response.SendErrorResponse(w, api.BADREQUEST, http.StatusBadRequest)
		return
	}

	team := domain.APIToDomainTeam(req)

	createdTeam, err := h.uc.CreateTeam(r.Context(), team)
	if err != nil {
		code, status := h.mapDomainErrorToAPI(err)
		response.SendErrorResponse(w, code, status)
		return
	}

	teamAPI := domain.DomainTeamToAPI(createdTeam)
	resp := domain.TeamResponse{Team: teamAPI}

	response.SendResponse(w, http.StatusCreated, resp)
}

func (h *TeamHandler) GetTeamGet(w http.ResponseWriter, r *http.Request, params api.GetTeamGetParams) {
	teamName := params.TeamName
	if err := validation.ValidateTeamName(teamName); err != nil {
		response.SendErrorResponse(w, api.BADREQUEST, http.StatusBadRequest)
		return
	}

	team, err := h.uc.GetTeamByName(r.Context(), teamName)
	if err != nil {
		code, status := h.mapDomainErrorToAPI(err)
		response.SendErrorResponse(w, code, status)
		return
	}

	teamAPI := domain.DomainTeamToAPI(team)
	response.SendResponse(w, http.StatusOK, teamAPI)
}

func (h *TeamHandler) mapDomainErrorToAPI(err error) (api.ErrorResponseErrorCode, int) {
	switch {
	case errors.Is(err, domain.ErrTeamExists):
		return api.TEAMEXISTS, http.StatusBadRequest
	case errors.Is(err, domain.ErrTeamNotFound):
		return api.NOTFOUND, http.StatusNotFound
	default:
		return api.INTERNAL, http.StatusInternalServerError
	}
}
