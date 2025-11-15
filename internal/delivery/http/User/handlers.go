package user

import (
	"net/http"
	"pr-reviewer/internal/api"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params api.GetUsersGetReviewParams) {

}

func (h *UserHandler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {

}
