package team

import (
	"net/http"
	"pr-reviewer/internal/api"
)

type TeamHandler struct {
}

func NewTeamHandler() *TeamHandler {
	return &TeamHandler{}
}

func (h *TeamHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {

}

func (h *TeamHandler) GetTeamGet(w http.ResponseWriter, r *http.Request, params api.GetTeamGetParams) {
	
}
