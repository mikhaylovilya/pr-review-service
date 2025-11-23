package endpoints

import "time"

type TeamDto struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDto `json:"members"`
}

type TeamMemberDto struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type UserDto struct {
	TeamMemberDto
	TeamName string `json:"team_name"`
}

type PullRequestShortDto struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
	Status          string `json:"status"`
}

type PullRequestDto struct {
	PullRequestShortDto
	AssignedReviewers []string  `json:"assigned_reviewers"`
	CreatedAt         time.Time `json:"created_at"`
	MergedAt          time.Time `json:"merged_at"`
}

type SetUserStatusDto struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type CreatePullRequestDto struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
}

type MergePullRequestDto struct {
	PullRequestId string `json:"pull_request_id"`
}

type ReassignDto struct {
	PullRequestId string `json:"pull_request_id"`
	OldUserId     string `json:"old_user_id"`
}
