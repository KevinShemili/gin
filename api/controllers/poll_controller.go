package controllers

import (
	"gin/api/requests"
	commands "gin/application/usecase/poll/commands/contracts"
	queries "gin/application/usecase/poll/queries/contracts"
	"gin/application/utility"
	"gin/domain/entities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PollController struct {
	CreatePollCommand  commands.ICreatePollCommand
	AddVoteCommand     commands.IAddVoteCommand
	DeletePollCommmand commands.IDeletePollCommand
	EndPollCommand     commands.IEndPollCommand
	UpdatePollCommand  commands.IUpdatePollCommand
	GetPollQuery       queries.IGetPollQuery
	GetPollsQuery      queries.IGetPollsQuery
	GetUserPollsQuery  queries.IGetUserPollsQuery
}

func NewPollController(
	CreatePollCommand commands.ICreatePollCommand,
	AddVoteCommand commands.IAddVoteCommand,
	DeletePollCommand commands.IDeletePollCommand,
	EndPollCommand commands.IEndPollCommand,
	GetPollQuery queries.IGetPollQuery,
	GetPollsQuery queries.IGetPollsQuery,
	GetUserPollsQuery queries.IGetUserPollsQuery,
	UpdatePollCommand commands.IUpdatePollCommand) *PollController {
	return &PollController{
		CreatePollCommand:  CreatePollCommand,
		AddVoteCommand:     AddVoteCommand,
		DeletePollCommmand: DeletePollCommand,
		EndPollCommand:     EndPollCommand,
		GetPollQuery:       GetPollQuery,
		GetPollsQuery:      GetPollsQuery,
		GetUserPollsQuery:  GetUserPollsQuery,
		UpdatePollCommand:  UpdatePollCommand}
}

// CreatePoll handles poll creation.
//
// @Summary Create a new poll
// @Description Create a new poll with a title, expiration time, and categories. The user must be authenticated.
// @Tags Polls
// @Accept json
// @Produce json
// @Param request body requests.CreatePollRequest true "Create Poll Request"
// @Success 200 {object} results.CreatePollResult "Poll created successfully"
// @Failure 400 {object} utility.ErrorCode "Bad Request - Invalid input"
// @Failure 401 {object} utility.ErrorCode "Unauthorized - Invalid or missing token"
// @Failure 500 {object} utility.ErrorCode "Internal server error"
// @Router /polls [post]
// @Security BearerAuth
func (uc *PollController) CreatePoll(c *gin.Context) {

	var request requests.CreatePollRequest

	if err := c.Bind(&request); err != nil {
		c.JSON(utility.BindFailure.StatusCode, utility.BindFailure)
		return
	}

	userAny, ok := c.Get("user")
	if !ok {
		c.JSON(utility.Unauthorized.StatusCode, utility.Unauthorized)
	}

	user, ok := userAny.(*entities.User)
	if !ok {
		c.JSON(utility.InternalServerError.StatusCode, utility.InternalServerError)
	}

	result, err := uc.CreatePollCommand.CreatePoll(&request, user)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// AddVote handles voting on a specific poll.
//
// @Summary Vote on a poll
// @Description Add a vote to a specific poll category by providing the poll ID in the route and the category ID in the request body. The user must be authenticated.
// @Tags Polls
// @Accept json
// @Produce json
// @Param id path int true "Poll ID"
// @Param request body requests.AddVoteRequest true "Add Vote Request"
// @Success 200 {object} bool "Vote added successfully"
// @Failure 400 {object} utility.ErrorCode "Bad Request - Invalid input"
// @Failure 401 {object} utility.ErrorCode "Unauthorized - Invalid or missing token"
// @Failure 404 {object} utility.ErrorCode "Poll or category not found"
// @Failure 500 {object} utility.ErrorCode "Internal server error"
// @Router /polls/{id}/vote [post]
// @Security BearerAuth
func (uc *PollController) AddVote(c *gin.Context) {

	var request requests.AddVoteRequest

	if err := c.Bind(&request); err != nil {
		c.JSON(utility.BindFailure.StatusCode, utility.BindFailure)
		return
	}

	pollIDString := c.Param("id")
	pollID, errParse := strconv.ParseUint(pollIDString, 10, 32)
	if errParse != nil {
		c.JSON(utility.RouteParameterCast.StatusCode, utility.RouteParameterCast)
		return
	}
	request.PollID = uint(pollID)

	userAny, ok := c.Get("user")
	if !ok {
		c.JSON(utility.Unauthorized.StatusCode, utility.Unauthorized)
	}

	user, ok := userAny.(*entities.User)
	if !ok {
		c.JSON(utility.InternalServerError.StatusCode, utility.InternalServerError)
	}

	result, err := uc.AddVoteCommand.AddVote(&request, user)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeletePoll handles deleting a specific poll.
//
// @Summary Delete a poll
// @Description Delete a poll by providing the poll ID in the route. The user must be authenticated.
// @Tags Polls
// @Accept json
// @Produce json
// @Param id path int true "Poll ID"
// @Success 200 {object} bool "Poll deleted successfully"
// @Failure 400 {object} utility.ErrorCode "Bad Request - Invalid input"
// @Failure 401 {object} utility.ErrorCode "Unauthorized - Invalid or missing token"
// @Failure 404 {object} utility.ErrorCode "Poll not found"
// @Failure 500 {object} utility.ErrorCode "Internal server error"
// @Router /polls/{id} [delete]
// @Security BearerAuth
func (uc *PollController) DeletePoll(c *gin.Context) {

	pollIDString := c.Param("id")
	pollID, errParse := strconv.ParseUint(pollIDString, 10, 32)
	if errParse != nil {
		c.JSON(utility.RouteParameterCast.StatusCode, utility.RouteParameterCast)
		return
	}

	result, err := uc.DeletePollCommmand.DeletePoll(uint(pollID))

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// EndPoll handles ending a specific poll.
//
// @Summary End a poll
// @Description End a poll by providing the poll ID in the route. The user must be authenticated.
// @Tags Polls
// @Accept json
// @Produce json
// @Param id path int true "Poll ID"
// @Success 200 {object} bool "Poll ended successfully"
// @Failure 400 {object} utility.ErrorCode "Bad Request - Invalid input"
// @Failure 401 {object} utility.ErrorCode "Unauthorized - Invalid or missing token"
// @Failure 404 {object} utility.ErrorCode "Poll not found"
// @Failure 500 {object} utility.ErrorCode "Internal server error"
// @Router /polls/{id}/end [patch]
// @Security BearerAuth
func (uc *PollController) EndPoll(c *gin.Context) {

	pollIDString := c.Param("id")
	pollID, errParse := strconv.ParseUint(pollIDString, 10, 32)
	if errParse != nil {
		c.JSON(utility.RouteParameterCast.StatusCode, utility.RouteParameterCast)
		return
	}

	result, err := uc.EndPollCommand.EndPoll(uint(pollID))

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetPoll godoc
// @Summary Get a specific poll
// @Description Retrieve a specific poll by its ID. Requires authentication.
// @Tags Polls
// @Accept json
// @Produce json
// @Param id path int true "Poll ID"
// @Success 200 {object} results.GetPollResult "Poll data"
// @Failure 400 {object} utility.ErrorCode "Bad Request - Invalid input"
// @Failure 401 {object} utility.ErrorCode "Unauthorized - Invalid or missing token"
// @Failure 404 {object} utility.ErrorCode "Poll not found"
// @Failure 500 {object} utility.ErrorCode "Internal server error"
// @Router /polls/{id} [get]
// @Security BearerAuth
func (uc *PollController) GetPoll(c *gin.Context) {

	pollIDString := c.Param("id")
	pollID, errParse := strconv.ParseUint(pollIDString, 10, 32)
	if errParse != nil {
		c.JSON(utility.RouteParameterCast.StatusCode, utility.RouteParameterCast)
		return
	}

	result, err := uc.GetPollQuery.GetPoll(uint(pollID))

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetPolls godoc
// @Summary Get polls with pagination and optional filter
// @Description Retrieves a paginated list of polls.
// @Tags Polls
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Items per page (default 10)"
// @Param filter query string false "Filter text (partial match against title or description)"
// @Success 200 {object} utility.PaginatedResponse[results.GetPollResult] "List of polls"
// @Failure 500 {object} utility.ErrorCode "Internal server error"
// @Router /polls [get]
func (uc *PollController) GetPolls(c *gin.Context) {

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	filter := c.Query("filter")

	params := utility.QueryParams{
		Page:     page,
		PageSize: pageSize,
		Filter:   filter,
	}

	result, err := uc.GetPollsQuery.GetPolls(params)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUserPolls godoc
// @Summary Get polls created by a specific user, with pagination/filter
// @Description Retrieves polls for the given user ID. Requires authentication.
// @Tags Polls
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Items per page (default 10)"
// @Param filter query string false "Filter text (partial match)"
// @Success 200 {object} utility.PaginatedResponse[results.GetPollResult] "List of user's polls"
// @Failure 400 {object} utility.ErrorCode "Bad Request - Invalid user ID"
// @Failure 401 {object} utility.ErrorCode "Unauthorized - Invalid or missing token"
// @Failure 404 {object} utility.ErrorCode "User not found"
// @Failure 500 {object} utility.ErrorCode "Internal server error"
// @Router /polls/users/{id} [get]
// @Security BearerAuth
func (uc *PollController) GetUserPolls(c *gin.Context) {

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	filter := c.Query("filter")

	params := utility.QueryParams{
		Page:     page,
		PageSize: pageSize,
		Filter:   filter,
	}

	userAny, ok := c.Get("user")
	if !ok {
		c.JSON(utility.Unauthorized.StatusCode, utility.Unauthorized)
	}

	user, ok := userAny.(*entities.User)
	if !ok {
		c.JSON(utility.InternalServerError.StatusCode, utility.InternalServerError)
	}

	result, err := uc.GetUserPollsQuery.GetPolls(user.ID, params)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdatePoll godoc
// @Summary Update a poll's details
// @Description Updates the specified poll's details, including title, expiration date, and categories.
// @Tags Polls
// @Accept json
// @Produce json
// @Param id path int true "Poll ID"
// @Param body body requests.UpdatePollRequest true "Poll update details"
// @Success 200 {object} bool "Poll updated successfully"
// @Failure 400 {object} utility.ErrorCode "Bad Request - Invalid input"
// @Failure 401 {object} utility.ErrorCode "Unauthorized - Invalid or missing authentication"
// @Failure 500 {object} utility.ErrorCode "Internal Server Error"
// @Router /polls/{id} [put]
// @Security BearerAuth
func (uc *PollController) UpdatePoll(c *gin.Context) {

	var request requests.UpdatePollRequest

	if err := c.Bind(&request); err != nil {
		c.JSON(utility.BindFailure.StatusCode, utility.BindFailure)
		return
	}

	pollIDString := c.Param("id")
	pollID, errParse := strconv.ParseUint(pollIDString, 10, 32)
	if errParse != nil {
		c.JSON(utility.RouteParameterCast.StatusCode, utility.RouteParameterCast)
		return
	}
	request.PollID = uint(pollID)

	userAny, ok := c.Get("user")
	if !ok {
		c.JSON(utility.Unauthorized.StatusCode, utility.Unauthorized)
	}

	user, ok := userAny.(*entities.User)
	if !ok {
		c.JSON(utility.InternalServerError.StatusCode, utility.InternalServerError)
	}

	result, err := uc.UpdatePollCommand.UpdatePoll(user.ID, &request)

	if err != nil {
		c.JSON(err.StatusCode, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
