package storage

import (
	"sync"

	"github.com/mikhaylovilya/pr-review-service/core/entities"
)

type InMemoryRepository interface {
	AddTeam(team entities.Team) error
	GetTeam(teamName string) (entities.Team, error)
	SetUserStatus(userId string, isActive bool) (entities.User, error)
	CreatePullRequest(pr entities.PullRequest) (entities.PullRequest, error)
	MergePullRequest(prId string) (entities.PullRequest, error)
	ReassignReviewer(prId string, reviewerId string) (entities.PullRequest, error)
	GetReview(userId string, prId string) ([]entities.PullRequest, error)
}

type InMemoryService struct {
	Teams        map[string]entities.Team
	Users        map[string]entities.User
	PullRequests map[string]entities.PullRequest

	mtx sync.Mutex
}

func NewInMemory() *InMemoryService {
	return &InMemoryService{
		Teams:        make(map[string]entities.Team),
		Users:        make(map[string]entities.User),
		PullRequests: make(map[string]entities.PullRequest),
	}
}

// type PersistentRepository interface {
// 	AddTeam(team entities.Team)
// }
