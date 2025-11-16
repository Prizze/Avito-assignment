package validation

import (
	"pr-reviewer/internal/api"
	"pr-reviewer/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePRId(t *testing.T) {
	tests := []struct {
		id        string
		wantError error
	}{
		{"", domain.ErrInvalidPullRequest},
		{"pr", domain.ErrInvalidPullRequest},
		{"px-123", domain.ErrInvalidPullRequest},
		{"pr-abc", domain.ErrInvalidPullRequest},
		{"pr-", domain.ErrInvalidPullRequest},
		{"pr-123", nil},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			err := ValidatePRId(tt.id)
			assert.Equal(t, tt.wantError, err)
		})
	}
}

func TestValidatePR(t *testing.T) {
	tests := []struct {
		pr        api.PostPullRequestCreateJSONRequestBody
		wantError error
	}{
		{api.PostPullRequestCreateJSONRequestBody{AuthorId: "invalid", PullRequestId: "pr-1", PullRequestName: "Name"}, domain.ErrInvalidUser},
		{api.PostPullRequestCreateJSONRequestBody{AuthorId: "u123", PullRequestId: "wrong", PullRequestName: "Name"}, domain.ErrInvalidPullRequest},
		{api.PostPullRequestCreateJSONRequestBody{AuthorId: "u123", PullRequestId: "pr-1", PullRequestName: ""}, domain.ErrInvalidPullRequest},
		{api.PostPullRequestCreateJSONRequestBody{AuthorId: "u123", PullRequestId: "pr-1", PullRequestName: "PR name"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.pr.PullRequestId, func(t *testing.T) {
			err := ValidatePR(tt.pr)
			assert.Equal(t, tt.wantError, err)
		})
	}
}

func TestValidateTeamName(t *testing.T) {
	tests := []struct {
		name      string
		wantError error
	}{
		{"", domain.ErrTeamNameEmpty},
		{"   ", domain.ErrTeamNameEmpty},
		{"Team A", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTeamName(tt.name)
			assert.Equal(t, tt.wantError, err)
		})
	}
}

func TestValidateTeam(t *testing.T) {
	tests := []struct {
		team      api.PostTeamAddJSONRequestBody
		wantError error
	}{
		{
			api.PostTeamAddJSONRequestBody{
				TeamName: "",
				Members:  []api.TeamMember{{UserId: "u1", Username: "Alice"}},
			},
			domain.ErrTeamNameEmpty,
		},
		{
			api.PostTeamAddJSONRequestBody{
				TeamName: "Team A",
				Members:  []api.TeamMember{},
			},
			domain.ErrTeamEmptyMembers,
		},
		{
			api.PostTeamAddJSONRequestBody{
				TeamName: "Team A",
				Members:  []api.TeamMember{{UserId: "invalid", Username: "Alice"}},
			},
			domain.ErrInvalidUser,
		},
		{
			api.PostTeamAddJSONRequestBody{
				TeamName: "Team A",
				Members:  []api.TeamMember{{UserId: "u1", Username: ""}},
			},
			domain.ErrInvalidUser,
		},
		{
			api.PostTeamAddJSONRequestBody{
				TeamName: "Team A",
				Members:  []api.TeamMember{{UserId: "u1", Username: "Alice"}},
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.team.TeamName, func(t *testing.T) {
			err := ValidateTeam(tt.team)
			assert.Equal(t, tt.wantError, err)
		})
	}
}

func TestValidateUserId(t *testing.T) {
	tests := []struct {
		id        string
		wantError error
	}{
		{"", domain.ErrInvalidUser},
		{"u", domain.ErrInvalidUser},
		{"x123", domain.ErrInvalidUser},
		{"u12a", domain.ErrInvalidUser},
		{"u123", nil},
		{"u0", nil},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			err := ValidateUserId(tt.id)
			assert.Equal(t, tt.wantError, err)
		})
	}
}
