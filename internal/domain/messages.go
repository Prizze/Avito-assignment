// Package domain messages.go содержит сообщения об ошибке для статуса из api
package domain

import "pr-reviewer/internal/api"

// Неизвестная ошибка
const (
	UnknownError = "Unknown error"
)

// Messages Сообщения об ошибке, соответствующие ErrorResponseErrorCode
var Messages = map[api.ErrorResponseErrorCode]string{
	api.NOCANDIDATE: "no active replacement candidate in team",
	api.NOTASSIGNED: "reviewer is not assigned to this PR",
	api.TEAMEXISTS:  "team_name already exists",
	api.NOTFOUND:    "resource not found",
	api.PREXISTS:    "PR id already exists",
	api.PRMERGED:    "cannot reassign on merged PR",
	api.BADREQUEST:  "invalid body request",
	api.INTERNAL:    "internal server error",
}
