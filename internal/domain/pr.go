package domain

import (
	"fmt"
	"pr-reviewer/internal/api"
	"strconv"
	"time"
)

const MaxReviewersNumber = 2

type PullRequestStatus string

const (
	PRStatusOpen   PullRequestStatus = "OPEN"
	PRStatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	ID                int
	Name              string
	AuthorID          int
	Status            PullRequestStatus
	AssignedReviewers []int
	CreatedAt         time.Time
	MergedAt          *time.Time
}

type CreatePullRequest struct {
	PullRequestId int
	Name          string
	AuthorId      int
}

func APIToDomainPullRequestCreate(pr api.PostPullRequestCreateJSONRequestBody) *CreatePullRequest {
	authorId, _ := strconv.Atoi(pr.AuthorId[1:])
	prId, _ := strconv.Atoi(pr.PullRequestId[3:])

	return &CreatePullRequest{
		Name:          pr.PullRequestName,
		AuthorId:      authorId,
		PullRequestId: prId,
	}
}

var MapDomainStatusToAPI = map[PullRequestStatus]api.PullRequestStatus{
	PRStatusOpen:   api.PullRequestStatusOPEN,
	PRStatusMerged: api.PullRequestStatusMERGED,
}

var MapStringToPullRequestStatus = map[string]PullRequestStatus{
	"OPEN":   PRStatusOpen,
	"MERGED": PRStatusMerged,
}

type PullRequestResponse struct {
	PullRequest api.PullRequest `json:"pr"`
}

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
