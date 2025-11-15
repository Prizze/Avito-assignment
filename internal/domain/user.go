package domain

import (
	"fmt"
	"pr-reviewer/internal/api"
	"strconv"
)

type User struct {
	ID       int
	Username string
	TeamName string
	IsActive bool
}

type SetUserIsActive struct {
	ID       int
	IsActive bool
}

func APIToDomainSetIsActive(set api.PostUsersSetIsActiveJSONRequestBody) *SetUserIsActive {
	id, _ := strconv.Atoi(set.UserId[1:])

	return &SetUserIsActive{
		ID:       id,
		IsActive: set.IsActive,
	}
}

type UserResponse struct {
	User api.User `json:"user"`
}

func DomainUserToAPI(u *User) api.User {
	return api.User{
		IsActive: u.IsActive,
		TeamName: u.TeamName,
		UserId:   fmt.Sprintf("u%d", u.ID),
		Username: u.Username,
	}
}

type UserReviews struct {
	UserID       string                 `json:"user_id"`
	PullRequests []api.PullRequestShort `json:"pull_requests"`
}
