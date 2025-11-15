package domain

import "errors"

var (
	ErrTeamNameEmpty    = errors.New("team name cannot be empty")
	ErrInvalidUser      = errors.New("invalid user data")
	ErrTeamEmptyMembers = errors.New("team members cannot be empty")
	ErrTeamExists       = errors.New("team_name already exists")
	ErrTeamNotFound     = errors.New("team not found")
)

var (
	ErrUserNotFound = errors.New("user not found")
)

var (
	ErrInvalidPullRequest  = errors.New("invalid pull_request data")
	ErrPullRequestExists   = errors.New("pull_request already exists")
	ErrPullRequestNotFound = errors.New("pull_request not found")
)
