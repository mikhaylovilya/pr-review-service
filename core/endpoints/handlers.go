package endpoints

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikhaylovilya/pr-review-service/core/entities"
	"github.com/mikhaylovilya/pr-review-service/core/storage"
)

type Repository struct {
	InMemory *storage.InMemoryRepository
	logger   *slog.Logger
}

func NewRepository(cache *storage.InMemoryRepository) *Repository {
	return &Repository{
		InMemory: cache,
		logger:   slog.Default(),
	}
}

func (r *Repository) AddTeamHandler(c *gin.Context) {
	var teamDto TeamDto
	if err := c.ShouldBindBodyWithJSON(&teamDto); err != nil {
		errResp := errors.New("bad request. Use /help to see API specification")
		c.AbortWithError(http.StatusBadRequest, errResp)
		return
	}
	if err := teamDto.Validate(); err != nil {
		errResp := errors.New("failed to validate Team: " + err.Error())
		c.AbortWithError(http.StatusBadRequest, errResp)
		return
	}

	team, err := entities.NewTeam(teamDto.TeamName, usersFromTeamMemberDtos(teamDto.Members, teamDto.TeamName))
	if err != nil {
		errResp := errors.New("failed to create entity Team: " + err.Error())
		c.AbortWithError(http.StatusInternalServerError, errResp)
		return
	}

	if err := (*r.InMemory).AddTeam(*team); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, *team)
}

func (r *Repository) GetTeamHandler(c *gin.Context) {
	teamName := c.Param("teamName")
	team, err := (*r.InMemory).GetTeam(teamName)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, &team)
}

func (r *Repository) SetUserStatusHandler(c *gin.Context) {
	var setUserStatusDto SetUserStatusDto
	if err := c.ShouldBindBodyWithJSON(&setUserStatusDto); err != nil {
		errResp := errors.New("failed to unmarshall JSON body: " + err.Error())
		c.AbortWithError(http.StatusBadRequest, errResp)
		return
	}
	if err := setUserStatusDto.Validate(); err != nil {
		errResp := errors.New("failed to validate SetUserStatusDto: " + err.Error())
		c.AbortWithError(http.StatusBadRequest, errResp)
		return
	}

	user, err := (*r.InMemory).SetUserStatus(setUserStatusDto.UserId, setUserStatusDto.IsActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (r *Repository) CreatePullRequestHandler(c *gin.Context) {
	var createPRDto CreatePullRequestDto
	if err := c.ShouldBindBodyWithJSON(&createPRDto); err != nil {
		errResp := errors.New("failed to unmarshall JSON body: " + err.Error())
		c.AbortWithError(http.StatusBadRequest, errResp)
		return
	}

	if err := createPRDto.Validate(); err != nil {
		errResp := errors.New("failed to validate PullRequest: " + err.Error())
		c.AbortWithError(http.StatusBadRequest, errResp)
		return
	}

	pr := *entities.NewPullRequest(createPRDto.PullRequestId, createPRDto.PullRequestName, createPRDto.AuthorId)
	pr, err := (*r.InMemory).CreatePullRequest(pr)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, pr)
}

func (r *Repository) MergePullRequestHandler(c *gin.Context) {
	var mergePRDto MergePullRequestDto
	if err := c.ShouldBindBodyWithJSON(&mergePRDto); err != nil {
		errResp := errors.New("failed to unmarshall JSON body: " + err.Error())
		c.AbortWithError(http.StatusBadRequest, errResp)
		return
	}

	if err := mergePRDto.Validate(); err != nil {
		errResp := errors.New("failed to validate MergePullRequestDto: " + err.Error())
		c.AbortWithError(http.StatusBadRequest, errResp)
		return
	}

	pr, err := (*r.InMemory).MergePullRequest(mergePRDto.PullRequestId)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, pr)
}

func (r *Repository) ReassignHandler(c *gin.Context) {
	var reassignDto ReassignDto
	if err := c.ShouldBindBodyWithJSON(&reassignDto); err != nil {
		errResp := errors.New("failder to unmarshall JSON body" + err.Error())
		c.AbortWithError(http.StatusBadRequest, errResp)
		return
	}

	if err := reassignDto.Validate(); err != nil {
		errResp := errors.New("failder to validate ReassignDto" + err.Error())
		c.AbortWithError(http.StatusBadRequest, errResp)
		return
	}

	pr, err := (*r.InMemory).ReassignReviewer(reassignDto.PullRequestId, reassignDto.OldUserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, pr)
}

func usersFromTeamMemberDtos(teamMember []TeamMemberDto, teamName string) []entities.User {
	users := make([]entities.User, 0, len(teamMember))
	for _, m := range teamMember {
		users = append(users, *entities.NewUser(m.Id, m.Name, teamName, m.IsActive))
	}
	return users
}
