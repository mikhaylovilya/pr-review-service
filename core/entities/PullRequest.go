package entities

import "time"

type PullRequest struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          string
	Status            string
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          time.Time
}

func NewPullRequest(prId string, prName string, authorId string, status string) *PullRequest {
	return &PullRequest{
		PullRequestId:   prId,
		PullRequestName: prName,
		AuthorId:        authorId,
		Status:          status,
		CreatedAt:       time.Now(),
	}
}
