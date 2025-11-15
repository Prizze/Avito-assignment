package domain

import "errors"

var (
	ErrTeamNameEmpty    = errors.New("team name cannot be empty")
	ErrInvalidUser      = errors.New("invalid user data")
	ErrTeamEmptyMembers = errors.New("team members cannot be empty")
	ErrTeamExists       = errors.New("team_name already exists")
	ErrTeamNotFound     = errors.New("team not found")
)
