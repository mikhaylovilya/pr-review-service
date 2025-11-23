package endpoints

import (
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
		errResp := entities.ErrGeneric("failed to unmarshal JSON body")
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	if err := teamDto.Validate(); err != nil {
		errResp := entities.ErrGeneric("failed to validate TeamDto: " + err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	team, err := entities.NewTeam(teamDto.TeamName, usersFromTeamMemberDtos(teamDto.Members, teamDto.TeamName))
	if err != nil {
		errResp := entities.ErrGeneric("failed to create entity Team: " + err.Error())
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	if err := (*r.InMemory).AddTeam(*team); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, teamToTeamDto(team))
}

func (r *Repository) GetTeamHandler(c *gin.Context) {
	teamName := c.Param("teamName")
	team, err := (*r.InMemory).GetTeam(teamName)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, teamToTeamDto(&team))
}

func (r *Repository) SetUserStatusHandler(c *gin.Context) {
	var setUserStatusDto SetUserStatusDto
	if err := c.ShouldBindBodyWithJSON(&setUserStatusDto); err != nil {
		errResp := entities.ErrGeneric("failed to unmarshal JSON body: " + err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	if err := setUserStatusDto.Validate(); err != nil {
		errResp := entities.ErrGeneric("failed to validate SetUserStatusDto: " + err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	user, err := (*r.InMemory).SetUserStatus(setUserStatusDto.UserId, setUserStatusDto.IsActive)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, userToUserDto(&user))
}

func (r *Repository) CreatePullRequestHandler(c *gin.Context) {
	var createPRDto CreatePullRequestDto
	if err := c.ShouldBindBodyWithJSON(&createPRDto); err != nil {
		errResp := entities.ErrGeneric("failed to unmarshal JSON body: " + err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	if err := createPRDto.Validate(); err != nil {
		errResp := entities.ErrGeneric("failed to validate CreatePullRequestDto: " + err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	pr := *entities.NewPullRequest(createPRDto.PullRequestId, createPRDto.PullRequestName, createPRDto.AuthorId)
	pr, err := (*r.InMemory).CreatePullRequest(pr)
	if err != nil {
		httpStatus := http.StatusInternalServerError
		if errResp, ok := err.(*entities.ErrorResponse); ok {
			if errResp.ErrorBody.Code == "PR_EXISTS" {
				httpStatus = http.StatusConflict
			} else {
				httpStatus = http.StatusNotFound
			}
		}
		c.JSON(httpStatus, err)
		return
	}

	c.JSON(http.StatusCreated, pullRequstToPullRequestDto(&pr))
}

func (r *Repository) MergePullRequestHandler(c *gin.Context) {
	var mergePRDto MergePullRequestDto
	if err := c.ShouldBindBodyWithJSON(&mergePRDto); err != nil {
		errResp := entities.ErrGeneric("failed to unmarshal JSON body: " + err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	if err := mergePRDto.Validate(); err != nil {
		errResp := entities.ErrGeneric("failed to validate MergePullRequestDto: " + err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	pr, err := (*r.InMemory).MergePullRequest(mergePRDto.PullRequestId)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, pullRequstToPullRequestDto(&pr))
}

func (r *Repository) ReassignHandler(c *gin.Context) {
	var reassignDto ReassignDto
	if err := c.ShouldBindBodyWithJSON(&reassignDto); err != nil {
		errResp := entities.ErrGeneric("failed to unmarshal JSON body" + err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	if err := reassignDto.Validate(); err != nil {
		errResp := entities.ErrGeneric("failed to validate ReassignDto" + err.Error())
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	pr, err := (*r.InMemory).ReassignReviewer(reassignDto.PullRequestId, reassignDto.OldUserId)
	if err != nil {
		var httpStatus int = http.StatusInternalServerError
		if errResp, ok := err.(*entities.ErrorResponse); ok {
			if errResp.ErrorBody.Code == "NOT_FOUND" {
				httpStatus = http.StatusNotFound
			} else {
				httpStatus = http.StatusConflict
			}
		}
		c.JSON(httpStatus, err)
		return
	}

	c.JSON(http.StatusOK, pullRequstToPullRequestDto(&pr))
}

func (r *Repository) GetReviewHandler(c *gin.Context) {
	userId := c.Param("userId")
	prs, err := (*r.InMemory).GetReview(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, err)
		return
	}

	sprs := shortPullRequestsFromPullRequests(prs)
	c.JSON(http.StatusOK, sprs)
}

func shortPullRequestsFromPullRequests(prs []entities.PullRequest) []PullRequestShortDto {
	sprs := make([]PullRequestShortDto, 0, len(prs))
	for _, pr := range prs {
		sprs = append(sprs, PullRequestShortDto{
			PullRequestId:   pr.PullRequestId,
			PullRequestName: pr.PullRequestName,
			AuthorId:        pr.AuthorId,
			Status:          pr.Status,
		})
	}
	return sprs
}

func usersFromTeamMemberDtos(teamMember []TeamMemberDto, teamName string) []entities.User {
	users := make([]entities.User, 0, len(teamMember))
	for _, m := range teamMember {
		users = append(users, *entities.NewUser(m.Id, m.Name, teamName, m.IsActive))
	}
	return users
}

func teamToTeamDto(team *entities.Team) *TeamDto {
	teamDto := TeamDto{
		TeamName: team.TeamName,
		Members:  make([]TeamMemberDto, 0, len(team.Members)),
	}
	for _, u := range team.Members {
		teamDto.Members = append(teamDto.Members,
			TeamMemberDto{
				Id:       u.Id,
				Name:     u.Name,
				IsActive: u.IsActive,
			},
		)
	}

	return &teamDto
}

func userToUserDto(user *entities.User) *UserDto {
	return &UserDto{
		TeamMemberDto: TeamMemberDto{
			Id:       user.Id,
			Name:     user.Name,
			IsActive: user.IsActive,
		},
		TeamName: user.TeamName,
	}
}

func pullRequstToPullRequestDto(pr *entities.PullRequest) *PullRequestDto {
	return &PullRequestDto{
		PullRequestShortDto: PullRequestShortDto{
			PullRequestId:   pr.PullRequestId,
			PullRequestName: pr.PullRequestName,
			AuthorId:        pr.AuthorId,
			Status:          pr.Status,
		},
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}
