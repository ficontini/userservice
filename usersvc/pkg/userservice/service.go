package userservice

import (
	"context"
	"os"

	"github.com/go-kit/log"

	"github.com/ficontini/user-search/types"
	"github.com/ficontini/user-search/usersvc/store"
)

type Service interface {
	GetUserByID(context.Context, int64) (*types.User, error)
	GetUsersByIDs(context.Context, []int64) ([]*types.User, error)
	GetUsersByCriteria(context.Context, types.SearchCriteria) ([]*types.User, error)
}

type basicService struct {
	store store.Store
}

func newBasicService(store store.Store) Service {
	return &basicService{
		store: store,
	}
}

func New(store store.Store) Service {
	var (
		logger log.Logger
		svc    Service
	)
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "service", "user")
		svc = newBasicService(store)
		svc = NewLoggingMiddleware(logger)(svc)
	}
	return svc

}
func (svc *basicService) GetUserByID(ctx context.Context, id int64) (*types.User, error) {
	return svc.store.GetByID(ctx, id)
}
func (svc *basicService) GetUsersByIDs(ctx context.Context, ids []int64) ([]*types.User, error) {
	return svc.store.GetByIDs(ctx, ids)
}
func (svc *basicService) GetUsersByCriteria(ctx context.Context, criteria types.SearchCriteria) ([]*types.User, error) {
	return svc.store.SearchByCriteria(ctx, criteria)
}
