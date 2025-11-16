package validation

import (
	"pr-reviewer/internal/domain"
	"unicode"
)

func ValidateUserId(id string) error {
	if len(id) < 2 || id[0] != 'u' {
		return domain.ErrInvalidUser
	}
	for _, r := range id[1:] {
		if !unicode.IsDigit(r) {
			return domain.ErrInvalidUser
		}
	}
	return nil
}
