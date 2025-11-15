package user

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

type UserHandler struct {
	uc userUC
}

func NewUserHandler(uc userUC) *UserHandler {
	return &UserHandler{
		uc: uc,
	}
}

func (h *UserHandler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params api.GetUsersGetReviewParams) {
	userIdAPI := params.UserId
	if err := validation.ValidateUserId(userIdAPI); err != nil {
		response.SendErrorResponse(w, api.BADREQUEST, http.StatusBadRequest)
		return
	}

	userID, _ := strconv.Atoi(userIdAPI[1:])
	userPRs, err := h.uc.GetUserPullRequests(r.Context(), userID)
	if err != nil {
		code, status := h.mapDomainErrorToAPI(err)
		response.SendErrorResponse(w, code, status)
		return
	}

	userPRsAPI := domain.DomainPRsToAPIShort(userPRs)
	resp := domain.UserReviews{UserID: userIdAPI, PullRequests: userPRsAPI}

	response.SendResponse(w, http.StatusOK, resp)
}

func (h *UserHandler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	var req api.PostUsersSetIsActiveJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, api.BADREQUEST, http.StatusBadRequest)
		return
	}

	if err := validation.ValidateUserId(req.UserId); err != nil {
		response.SendErrorResponse(w, api.BADREQUEST, http.StatusBadRequest)
		return
	}

	setIsActive := domain.APIToDomainSetIsActive(req)

	user, err := h.uc.SetUserIsActive(r.Context(), setIsActive)
	if err != nil {
		code, status := h.mapDomainErrorToAPI(err)
		response.SendErrorResponse(w, code, status)
		return
	}

	userAPI := domain.DomainUserToAPI(user)
	resp := domain.UserResponse{User: userAPI}

	response.SendResponse(w, http.StatusOK, resp)
}

func (h *UserHandler) mapDomainErrorToAPI(err error) (api.ErrorResponseErrorCode, int) {
	log.Println(err)
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return api.NOTFOUND, http.StatusNotFound
	default:
		return api.INTERNAL, http.StatusInternalServerError
	}
}
