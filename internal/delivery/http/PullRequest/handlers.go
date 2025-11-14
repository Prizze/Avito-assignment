package pullrequest

import "net/http"

type PRHandler struct {
}

func NewPRHandler() *PRHandler {
	return &PRHandler{}
}

func (h *PRHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {

}

func (h *PRHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {

}

func (h *PRHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {

}
