package endpoints

import (
	"time"
)

type TeamDto struct {
	TeamName string
	Members  []TeamMemberDto
}

type TeamMemberDto struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	IsActive bool   `json:"IsActive"`
}

type UserDto struct {
	TeamMemberDto
	TeamName string `json:"TeamName"`
}

type PullRequestShortDto struct {
	PullRequestId   string
	PullRequestName string
	AuthorId        string
	Status          string
}

type PullRequestDto struct {
	PullRequestShortDto
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          time.Time
}

// const (
// 	ErrCodeTeamExists = iota
// 	ErrCodePRExists
// 	ErrCodePRMerged
// 	ErrCodeNotAssigned
// 	ErrCodeNoCandidate
// 	ErrCodeNotFound
// )
