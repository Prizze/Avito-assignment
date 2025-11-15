package validation

import (
	"pr-reviewer/internal/api"
	"pr-reviewer/internal/domain"
	"strings"
	"unicode"
)

func ValidatePR(pr api.PostPullRequestCreateJSONRequestBody) error {
	if err := ValidateUserId(pr.AuthorId); err != nil {
		return domain.ErrInvalidUser
	}

	if err := ValidatePRId(pr.PullRequestId); err != nil {
		return domain.ErrInvalidPullRequest
	}

	if strings.TrimSpace(pr.PullRequestName) == "" {
		return domain.ErrInvalidPullRequest
	}

	return nil
}

func ValidatePRId(id string) error {
	if len(id) < 4 || id[:3] != "pr-" {
		return domain.ErrInvalidPullRequest
	}
	nums := id[3:]
	if len(nums) == 0 {
		return domain.ErrInvalidPullRequest
	}

	for _, r := range nums {
		if !unicode.IsDigit(r) {
			return domain.ErrInvalidPullRequest
		}
	}

	return nil
}
