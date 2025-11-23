package repository

import (
	"sync"

	"github.com/mikhaylovilya/pr-review-service/core/entities"
)

type InMemoryRepository interface {
	AddTeam(team entities.Team) error
	// AddUsers(users []entities.User) error
	GetTeam(teamName string) (entities.Team, error)
	SetUserIsActive(userId string, isActive bool) (entities.User, error)
}

type InMemoryService struct {
	Teams map[string]entities.Team
	Users map[string]entities.User

	mtx sync.Mutex
}

func NewInMemory() *InMemoryService {
	return &InMemoryService{
		Teams: make(map[string]entities.Team),
		Users: make(map[string]entities.User),
	}
}

// type PersistentRepository interface {
// 	AddTeam(team entities.Team)
// }
