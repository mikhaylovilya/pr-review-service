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
	// if pr.Status == "MERGED" {
	// 	return ErrPRMerged(pr.PullRequestId)
	// }
	pr.Status = "MERGED"
	pr.MergedAt = time.Now()
	// return nil
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

// func (pr *PullRequest) ReassignReviewer(reassignee *User, users []User) error {

// }
