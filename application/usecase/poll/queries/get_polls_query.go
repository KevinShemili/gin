package queries

import (
	"gin/api/requests"
	repo "gin/application/repository/contracts"
	"gin/application/usecase/poll/queries/contracts"
	"gin/application/usecase/poll/results"
	"gin/application/utility"

	"github.com/go-playground/validator/v10"
)

type GetPollsQuery struct {
	UnitOfWork repo.IUnitOfWork
	Validator  *validator.Validate
}

func NewGetPollsQuery(UnitOfWork repo.IUnitOfWork, Validator *validator.Validate) contracts.IGetPollsQuery {
	return &GetPollsQuery{UnitOfWork: UnitOfWork, Validator: Validator}
}

func (r GetPollsQuery) GetPolls(request *requests.GetPollsRequest) (utility.PaginatedResponse[results.GetPollResult], *utility.ErrorCode) {

	// validate request
	if err := r.Validator.Struct(request); err != nil {
		return utility.PaginatedResponse[results.GetPollResult]{}, utility.ValidationError.WithDescription(err.Error())
	}

	// get polls
	polls, err := r.UnitOfWork.IPollRepository().GetPollsPaginated(request.QueryParams, request.ShowActiveOnly)
	if err != nil {
		return utility.PaginatedResponse[results.GetPollResult]{}, utility.InternalServerError.WithDescription(err.Error())
	}
	if len(polls.Data) == 0 {
		return utility.PaginatedResponse[results.GetPollResult]{}, nil
	}

	// map results
	results := utility.MapPoll(polls)

	return results, nil
}
