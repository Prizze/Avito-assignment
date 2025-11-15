package pullrequest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"pr-reviewer/internal/api"
	"pr-reviewer/internal/domain"
	"pr-reviewer/internal/pkg/response"
	"pr-reviewer/internal/pkg/validation"
	"strconv"
)

type PRHandler struct {
	uc prUC
}

func NewPRHandler(uc prUC) *PRHandler {
	return &PRHandler{
		uc: uc,
	}
}

func (h *PRHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	var req api.PostPullRequestCreateJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, api.BADREQUEST, http.StatusBadRequest)
		return
	}

	if err := validation.ValidatePR(req); err != nil {
		response.SendErrorResponse(w, api.BADREQUEST, http.StatusBadRequest)
		return
	}

	cr := domain.APIToDomainPullRequestCreate(req)

	createdPR, err := h.uc.CreatePullRequest(r.Context(), cr)
	if err != nil {
		code, status := h.mapDomainErrorToAPI(err)
		response.SendErrorResponse(w, code, status)
		return
	}

	prAPI := domain.DomainPRToAPI(createdPR)
	resp := domain.PullRequestResponse{PullRequest: prAPI}

	response.SendResponse(w, http.StatusCreated, resp)
}

func (h *PRHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	var req api.PostPullRequestMergeJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, api.BADREQUEST, http.StatusBadRequest)
		return
	}

	if err := validation.ValidatePRId(req.PullRequestId); err != nil {
		response.SendErrorResponse(w, api.BADREQUEST, http.StatusBadRequest)
		return
	}

	prID, _ := strconv.Atoi(req.PullRequestId[3:])

	mergedPR, err := h.uc.MergePullRequest(r.Context(), prID)
	if err != nil {
		code, status := h.mapDomainErrorToAPI(err)
		response.SendErrorResponse(w, code, status)
		return
	}

	prAPI := domain.DomainPRToAPI(mergedPR)
	resp := domain.PullRequestResponse{PullRequest: prAPI}

	response.SendResponse(w, http.StatusOK, resp)
}

func (h *PRHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {

}

func (h *PRHandler) mapDomainErrorToAPI(err error) (api.ErrorResponseErrorCode, int) {
	log.Println(err)
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return api.NOTFOUND, http.StatusNotFound
	case errors.Is(err, domain.ErrPullRequestExists):
		return api.PREXISTS, http.StatusConflict
	case errors.Is(err, domain.ErrPullRequestNotFound):
		return api.NOTFOUND, http.StatusNotFound
	default:
		return api.INTERNAL, http.StatusInternalServerError
	}
}
