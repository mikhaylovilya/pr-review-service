package storage

import (
	"maps"
	"slices"

	"github.com/mikhaylovilya/pr-review-service/core/entities"
)

func (mem *InMemoryService) AddTeam(team entities.Team) error {
	mem.mtx.Lock()
	defer mem.mtx.Unlock()

	if _, ok := mem.Teams[team.TeamName]; ok {
		return entities.ErrTeamExists(team.TeamName)
	}

	for _, u := range team.Members {
		if _, ok := mem.Users[u.Id]; ok {
			return entities.ErrUserExists(u.Id)
		}
	}

	for _, u := range team.Members {
		mem.Users[u.Id] = u
	}
	mem.Teams[team.TeamName] = team
	return nil
}

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

func (mem *InMemoryService) SetUserStatus(userId string, isActive bool) (entities.User, error) {
	mem.mtx.Lock()
	defer mem.mtx.Unlock()

	user, ok := mem.Users[userId]
	if !ok {
		// return entities.User{}, errors.New(userId + "does not exists")
		return entities.User{}, entities.ErrNotFound(userId)
	}

	user.SetStatus(isActive)
	mem.Users[userId] = user
	for i, m := range mem.Teams[user.TeamName].Members {
		if m.Id == userId {
			m.SetStatus(isActive)
			mem.Teams[user.TeamName].Members[i] = user
		}
	}

	return user, nil
}

func (mem *InMemoryService) CreatePullRequest(pr entities.PullRequest) (entities.PullRequest, error) {
	mem.mtx.Lock()
	defer mem.mtx.Unlock()

	if _, ok := mem.PullRequests[pr.PullRequestId]; ok {
		return entities.PullRequest{}, entities.ErrPRExists(pr.PullRequestId)
	}

	if _, ok := mem.Users[pr.AuthorId]; !ok {
		return entities.PullRequest{}, entities.ErrNotFound(pr.AuthorId)
	}

	teamName := mem.Users[pr.AuthorId].TeamName
	if err := pr.AssignReviewers(mem.Teams[teamName].Members); err != nil {
		return entities.PullRequest{}, err
	}

	mem.PullRequests[pr.PullRequestId] = pr
	return pr, nil
}

func (mem *InMemoryService) MergePullRequest(prId string) (entities.PullRequest, error) {
	mem.mtx.Lock()
	defer mem.mtx.Unlock()

	pr, ok := mem.PullRequests[prId]
	if !ok {
		return entities.PullRequest{}, entities.ErrNotFound(pr.PullRequestId)
	}

	pr.Merge()
	mem.PullRequests[pr.PullRequestId] = pr
	return pr, nil
}

func (mem *InMemoryService) ReassignReviewer(prId string, reviewerId string) (entities.PullRequest, error) {
	mem.mtx.Lock()
	defer mem.mtx.Unlock()

	if _, ok := mem.PullRequests[prId]; !ok {
		return entities.PullRequest{}, entities.ErrNotFound(prId)
	}

	if _, ok := mem.Users[reviewerId]; !ok {
		return entities.PullRequest{}, entities.ErrNotFound(reviewerId)
	}

	pr := mem.PullRequests[prId]
	if err := pr.ReassignReviewer(reviewerId, mem.Teams[mem.Users[reviewerId].TeamName].Members); err != nil {
		return entities.PullRequest{}, err
	}

	mem.PullRequests[prId] = pr
	return pr, nil
}

func (mem *InMemoryService) GetReview(userId string) ([]entities.PullRequest, error) {
	mem.mtx.Lock()
	defer mem.mtx.Unlock()

	if _, ok := mem.Users[userId]; !ok {
		return []entities.PullRequest{}, entities.ErrNotFound(userId)
	}

	prs := make([]entities.PullRequest, 0, len(mem.PullRequests))
	for _, pr := range slices.Collect(maps.Values(mem.PullRequests)) {
		for _, rev := range pr.AssignedReviewers {
			if userId == rev {
				prs = append(prs, pr)
			}
		}
	}

	return prs, nil
}
