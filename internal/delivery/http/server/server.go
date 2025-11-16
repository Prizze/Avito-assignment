// Package server server.go композирует handlers для удовлетворения ServerInterface
package server

import (
	"net/http"
	"pr-reviewer/internal/api"
	pullrequest "pr-reviewer/internal/delivery/http/PullRequest"
	team "pr-reviewer/internal/delivery/http/Team"
	user "pr-reviewer/internal/delivery/http/User"
)

type Server struct {
	User *user.UserHandler
	Team *team.TeamHandler
	PR   *pullrequest.PRHandler
}

func NewServer(u *user.UserHandler, t *team.TeamHandler, pr *pullrequest.PRHandler) *Server {
	return &Server{
		User: u,
		Team: t,
		PR:   pr,
	}
}

func (s *Server) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	s.PR.PostPullRequestCreate(w, r)
}

func (s *Server) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	s.PR.PostPullRequestMerge(w, r)
}

func (s *Server) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	s.PR.PostPullRequestReassign(w, r)
}

func (s *Server) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	s.Team.PostTeamAdd(w, r)
}

func (s *Server) GetTeamGet(w http.ResponseWriter, r *http.Request, params api.GetTeamGetParams) {
	s.Team.GetTeamGet(w, r, params)
}

func (s *Server) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params api.GetUsersGetReviewParams) {
	s.User.GetUsersGetReview(w, r, params)
}

func (s *Server) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	s.User.PostUsersSetIsActive(w, r)
}
