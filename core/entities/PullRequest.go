package entities

import (
	"time"
)

type PullRequest struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          string
	Status            string
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          time.Time
}

func NewPullRequest(prId string, prName string, authorId string) *PullRequest {
	return &PullRequest{
		PullRequestId:   prId,
		PullRequestName: prName,
		AuthorId:        authorId,
		Status:          "OPEN",
		CreatedAt:       time.Now(),
	}
}

func (pr *PullRequest) Merge() {
	if pr.Status == "OPEN" {
		pr.MergedAt = time.Now()
	}
	pr.Status = "MERGED"
}

func (pr *PullRequest) AssignReviewers(users []User) error {
	assignedReviewers := make([]string, 0, 2)
	for _, u := range users {
		if u.IsActive && u.Id != pr.AuthorId {
			assignedReviewers = append(assignedReviewers, u.Id)

			if len(assignedReviewers) == 2 {
				break
			}
		}
	}

	pr.AssignedReviewers = assignedReviewers
	return nil
}

func (pr *PullRequest) ReassignReviewer(reviewerId string, users []User) error {
	if pr.Status == "MERGED" {
		return ErrPRMerged(pr.PullRequestId)
	}

	reviewerAssigned := false
	for _, rId := range pr.AssignedReviewers {
		if rId == reviewerId {
			reviewerAssigned = true
		}
	}

	if !reviewerAssigned {
		return ErrNotAssigned(reviewerId)
	}

	var newReviewer *User
	for _, u := range users {
		if u.IsActive && u.Id != pr.AuthorId && u.Id != reviewerId {
			newReviewer = &u
			break
		}
	}

	if newReviewer == nil {
		return ErrNoCandidate(reviewerId)
	}

	//TODO
	// pr.AuthorId = newReviewer.Id

	return nil
}
