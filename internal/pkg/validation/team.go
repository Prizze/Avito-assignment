package validation

import (
	"pr-reviewer/internal/api"
	"pr-reviewer/internal/domain"
	"strings"
)

func ValidateTeam(team api.PostTeamAddJSONRequestBody) error {
	if err := ValidateTeamName(team.TeamName); err != nil {
		return err
	}

	if len(team.Members) == 0 {
		return domain.ErrTeamEmptyMembers
	}

	for _, m := range team.Members {
		// Проверка правильности написания user_id
		if err := ValidateUserId(m.UserId); err != nil {
			return err
		}

		if strings.TrimSpace(m.Username) == "" {
			return domain.ErrInvalidUser
		}
	}

	return nil
}

func ValidateTeamName(name string) error {
	if strings.TrimSpace(name) == "" {
		return domain.ErrTeamNameEmpty
	}
	return nil
}
