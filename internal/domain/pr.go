package domain

import (
	"fmt"
	"pr-reviewer/internal/api"
	"strconv"
	"time"
)

// MaxReviewersNumber максимальное количество ревьюверов на PullRequest
const MaxReviewersNumber = 2

// PullRequestStatus тип для статуса PullRequest
type PullRequestStatus string

// Статусы PullRequest
const (
	PRStatusOpen   PullRequestStatus = "OPEN"
	PRStatusMerged PullRequestStatus = "MERGED"
)

// PullRequest domain модель для PullRequest
type PullRequest struct {
	ID                int
	Name              string
	AuthorID          int
	Status            PullRequestStatus
	AssignedReviewers []int
	CreatedAt         time.Time
	MergedAt          *time.Time
}

// CreatePullRequest domain модель для создания PullRequest
type CreatePullRequest struct {
	PullRequestId int
	Name          string
	AuthorId      int
}

// APIToDomainPullRequestCreate маппит API запрос в domain CreatePullRequest
func APIToDomainPullRequestCreate(pr api.PostPullRequestCreateJSONRequestBody) *CreatePullRequest {
	authorId, _ := strconv.Atoi(pr.AuthorId[1:])
	prId, _ := strconv.Atoi(pr.PullRequestId[3:])

	return &CreatePullRequest{
		Name:          pr.PullRequestName,
		AuthorId:      authorId,
		PullRequestId: prId,
	}
}

// MapDomainStatusToAPI маппинг domain PullRequestStatus в api PullRequestStatus
var MapDomainStatusToAPI = map[PullRequestStatus]api.PullRequestStatus{
	PRStatusOpen:   api.PullRequestStatusOPEN,
	PRStatusMerged: api.PullRequestStatusMERGED,
}

// MapStringToPullRequestStatusShort маппинг domain PullRequestStatus в api PullRequestStatusShort
var MapStringToPullRequestStatusShort = map[PullRequestStatus]api.PullRequestShortStatus{
	PRStatusOpen:   api.PullRequestShortStatusOPEN,
	PRStatusMerged: api.PullRequestShortStatusMERGED,
}

// MapStringToPullRequestStatus маппинг string Status в domain PullRequestStatus
var MapStringToPullRequestStatus = map[string]PullRequestStatus{
	"OPEN":   PRStatusOpen,
	"MERGED": PRStatusMerged,
}

// PullRequestResponse возвращаемое значение
type PullRequestResponse struct {
	PullRequest api.PullRequest `json:"pr"`
}

// DomainPRToAPI маппит domain PullRequest в api PullRequest
func DomainPRToAPI(pr *PullRequest) api.PullRequest {
	var reviewers []string
	for _, id := range pr.AssignedReviewers {
		reviewers = append(reviewers, fmt.Sprintf("u%d", id))
	}
	return api.PullRequest{
		PullRequestId:     fmt.Sprintf("pr-%d", pr.ID),
		AuthorId:          fmt.Sprintf("u%d", pr.AuthorID),
		CreatedAt:         &pr.CreatedAt,
		MergedAt:          pr.MergedAt,
		Status:            MapDomainStatusToAPI[pr.Status],
		PullRequestName:   pr.Name,
		AssignedReviewers: reviewers,
	}
}

// DomainPRToAPI маппит domain PullRequest в api PullRequestShort
func DomainPRToAPIShort(pr *PullRequest) api.PullRequestShort {
	return api.PullRequestShort{
		PullRequestId:   fmt.Sprintf("pr-%d", pr.ID),
		PullRequestName: pr.Name,
		AuthorId:        fmt.Sprintf("u%d", pr.AuthorID),
		Status:          MapStringToPullRequestStatusShort[pr.Status],
	}
}

// DomainPRToAPI маппит domain []PullRequest в api []PullRequestShort
func DomainPRsToAPIShort(prs []PullRequest) []api.PullRequestShort {
	prsAPI := make([]api.PullRequestShort, 0, len(prs))
	for _, pr := range prs {
		prAPI := DomainPRToAPIShort(&pr)
		prsAPI = append(prsAPI, prAPI)
	}

	return prsAPI
}

// ReassingReviewer domain запрос на переназначение ревьювера
type ReassingReviewer struct {
	PullRequestID int
	UserID        int
}

// ReassignResponse ответ на запрос переназаначения интервьювера
type ReassignResponse struct {
	PullRequest api.PullRequest `json:"pr"`
	ReplacedBy  string          `json:"replaced_by"`
}

// APIReassignToDomain маппит api PostPullRequestReassignJSONRequestBody в domain ReassingReviewer
func APIReassignToDomain(ras api.PostPullRequestReassignJSONRequestBody) *ReassingReviewer {
	prID, _ := strconv.Atoi(ras.PullRequestId[3:])
	userID, _ := strconv.Atoi(ras.OldUserId[1:])
	return &ReassingReviewer{
		PullRequestID: prID,
		UserID:        userID,
	}
}
