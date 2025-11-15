package domain

import (
	"fmt"
	"pr-reviewer/internal/api"
	"strconv"
)

type Team struct {
	ID      int
	Name    string
	Members []TeamMember
}

type TeamMember struct {
	IsActive bool
	UserID   int
	Username string
}

func APIToDomainTeam(ta api.Team) *Team {
	members := make([]TeamMember, 0, len(ta.Members))

	for _, m := range ta.Members {
		id, _ := strconv.Atoi(m.UserId[1:])

		members = append(members, TeamMember{
			IsActive: m.IsActive,
			UserID:   id,
			Username: m.Username,
		})
	}

	return &Team{
		Name:    ta.TeamName,
		Members: members,
	}
}

type TeamResponse struct {
	Team api.Team `json:"team"`
}

func DomainTeamToAPI(team *Team) api.Team {
	members := make([]api.TeamMember, 0, len(team.Members))

	for _, m := range team.Members {
		members = append(members, api.TeamMember{
			UserId:   fmt.Sprintf("u%d", m.UserID),
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	return api.Team{
		TeamName: team.Name,
		Members:  members,
	}
}
