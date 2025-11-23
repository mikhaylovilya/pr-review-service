package repository

import (
	"github.com/mikhaylovilya/pr-review-service/core/entities"
)

func (mem *InMemoryService) AddTeam(team entities.Team) error {
	mem.mtx.Lock()
	defer mem.mtx.Unlock()

	if _, ok := mem.Teams[team.TeamName]; ok {
		// return &endpoints.ErrTeamExists
		// return errors.New(team.TeamName + "already exists")
		return entities.ErrTeamExists(team.TeamName)
	}

	for _, u := range team.Members {
		if _, ok := mem.Users[u.Id]; ok {
			// return errors.New(u.Id + "already exists")
			return entities.ErrUserExists(u.Id)
		}
	}
	// mem.AddUsers(mem.Teams[team.TeamName].Members)

	for _, u := range team.Members {
		mem.Users[u.Id] = u
	}
	mem.Teams[team.TeamName] = team
	return nil
}

// func (mem *InMemory) AddUsers(users []entities.User) error {
// 	mem.mtx.Lock()
// 	defer mem.mtx.Unlock()

// 	for _, u := range users {
// 		if _, ok := mem.Users[u.Id]; ok {
// 			// return &endpoints.ErrUserExists
// 			return errors.New(u.Id + "already exists")
// 		}

// 		mem.Users[u.Id] = u
// 	}

// 	return nil
// }

func (mem *InMemoryService) GetTeam(teamName string) (entities.Team, error) {
	mem.mtx.Lock()
	defer mem.mtx.Unlock()

	if _, ok := mem.Teams[teamName]; !ok {
		// return entities.Team{}, &endpoints.ErrNotFound
		// return entities.Team{}, errors.New(teamName + "does not exists")
		return entities.Team{}, entities.ErrNotFound(teamName)
	}

	team := mem.Teams[teamName]
	return team, nil
}

func (mem *InMemoryService) SetUserIsActive(userId string, isActive bool) (entities.User, error) {
	mem.mtx.Lock()
	defer mem.mtx.Unlock()

	user, ok := mem.Users[userId]
	if !ok {
		// return entities.User{}, errors.New(userId + "does not exists")
		return entities.User{}, entities.ErrNotFound(userId)
	}

	user.SetStatus(isActive)
	mem.Users[userId] = user

	return user, nil
}

// func (mem *InMemory) CreatePullRequest(pr entit)
