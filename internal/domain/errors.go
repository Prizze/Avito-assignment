// Package domain errors.go внутренние ошибки
package domain

import "errors"

// Ошибки для Team
var (
	ErrTeamNameEmpty    = errors.New("team name cannot be empty")
	ErrInvalidUser      = errors.New("invalid user data")
	ErrTeamEmptyMembers = errors.New("team members cannot be empty")
	ErrTeamExists       = errors.New("team_name already exists")
	ErrTeamNotFound     = errors.New("team not found")
)

// Ошибки для User
var (
	ErrUserNotFound = errors.New("user not found")
)

// Ошибки для PullRequest
var (
	ErrInvalidPullRequest   = errors.New("invalid pull_request data")
	ErrPullRequestExists    = errors.New("pull_request already exists")
	ErrPullRequestNotFound  = errors.New("pull_request not found")
	ErrNoAvailableCandidats = errors.New("no active replacement candidate in team")
	ErrPullRequestIsMerged  = errors.New("pull_request merged already")
	ErrNotAssigned          = errors.New("reviewer is not assigned to this PR")
)
