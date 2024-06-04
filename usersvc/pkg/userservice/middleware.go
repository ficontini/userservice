package userservice

import (
	"context"

	"github.com/ficontini/user-search/types"
	"github.com/go-kit/log"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func NewLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}
func (mw loggingMiddleware) GetUserByID(ctx context.Context, id int64) (user *types.User, err error) {
	defer func() {
		mw.logger.Log("method", "GetUserbyID", "id", id, "err", err)
	}()
	user, err = mw.next.GetUserByID(ctx, id)
	return user, err
}
func (mw loggingMiddleware) GetUsersByIDs(ctx context.Context, ids []int64) (users []*types.User, err error) {
	defer func() {
		mw.logger.Log("method", "GetUsersbyIDs", "err", err)
	}()
	users, err = mw.next.GetUsersByIDs(ctx, ids)
	return users, err
}
func (mw loggingMiddleware) GetUsersByCriteria(ctx context.Context, criteria types.SearchCriteria) (users []*types.User, err error) {
	defer func() {
		mw.logger.Log("method", "GetUsersByCriteria", "criteria:", "city", criteria.City(), "name", criteria.Name(), "err", err)
	}()
	users, err = mw.next.GetUsersByCriteria(ctx, criteria)
	return users, err
}
