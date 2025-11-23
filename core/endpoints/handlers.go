package endpoints

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikhaylovilya/pr-review-service/core/entities"
	"github.com/mikhaylovilya/pr-review-service/core/repository"
)

type Repository struct {
	Cache  *repository.InMemoryRepository
	logger *slog.Logger
}

func NewRepository(cache *repository.InMemoryRepository) *Repository {
	return &Repository{
		Cache:  cache,
		logger: slog.Default(),
	}
}

func (r *Repository) AddTeamHandler(c *gin.Context) {
	var teamDto TeamDto
	if err := c.ShouldBindBodyWithJSON(&teamDto); err != nil {
		errResp := errors.New("bad request. Use /help to see API specification")
		c.AbortWithError(http.StatusBadRequest, errResp)
		return
	}
	if err := teamDto.ValidateTeamDto(); err != nil {
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

	if err := (*r.Cache).AddTeam(*team); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, *team)
}

func (r *Repository) GetTeamHandler(c *gin.Context) {
	teamName := c.Param("teamName")
	team, err := (*r.Cache).GetTeam(teamName)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, &team)
}

func usersFromTeamMemberDtos(teamMember []TeamMemberDto, teamName string) []entities.User {
	users := make([]entities.User, 0, len(teamMember))
	for _, m := range teamMember {
		users = append(users, *entities.NewUser(m.Id, m.Name, teamName, m.IsActive))
	}
	return users
}
