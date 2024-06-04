package userendpoint

import (
	"context"

	"github.com/ficontini/user-search/types"
	"github.com/ficontini/user-search/usersvc/pkg/userservice"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	GetUserByIDEndpoint        endpoint.Endpoint
	GetUsersByIDsEndpoint      endpoint.Endpoint
	GetUsersByCriteriaEndpoint endpoint.Endpoint
}

func New(svc userservice.Service) Set {
	return Set{
		GetUserByIDEndpoint:        makeGetUserByIDEndpoint(svc),
		GetUsersByIDsEndpoint:      makeGetUsersByIDsEndpoint(svc),
		GetUsersByCriteriaEndpoint: makeGetUsersByCriteriaEndpoint(svc),
	}
}
func (s Set) GetUserByID(ctx context.Context, id int64) (*types.User, error) {
	resp, err := s.GetUserByIDEndpoint(ctx, &UserRequest{ID: id})
	if err != nil {
		return nil, err
	}
	response := resp.(*UserResponse)
	return &types.User{
		ID:      response.User.ID,
		FName:   response.User.FName,
		City:    response.User.City,
		Phone:   response.User.Phone,
		Height:  response.User.Height,
		Married: response.User.Married,
	}, nil
}
func (s Set) GetUsersByIDs(ctx context.Context, ids []int64) ([]*types.User, error) {
	resp, err := s.GetUsersByIDsEndpoint(ctx, &UsersRequest{IDs: ids})
	if err != nil {
		return nil, err
	}
	response := resp.(*UsersResponse)
	var users []*types.User
	for _, r := range response.Users {
		users = append(users, &types.User{
			ID:      r.ID,
			FName:   r.FName,
			City:    r.City,
			Phone:   r.Phone,
			Height:  r.Height,
			Married: r.Married,
		})
	}
	return users, nil
}
func (s Set) GetUsersByCriteria(ctx context.Context, req types.SearchCriteria) ([]*types.User, error) {
	resp, err := s.GetUsersByCriteriaEndpoint(ctx, &SearchRequest{City: req.City(), Name: req.Name()})
	if err != nil {
		return nil, err
	}
	response := resp.(*UsersResponse)
	var users []*types.User
	for _, r := range response.Users {
		users = append(users, &types.User{
			ID:      r.ID,
			FName:   r.FName,
			City:    r.City,
			Phone:   r.Phone,
			Height:  r.Height,
			Married: r.Married,
		})
	}
	return users, nil
}

func makeGetUserByIDEndpoint(svc userservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*UserRequest)
		res, err := svc.GetUserByID(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return &UserResponse{
			User: User{
				ID:      res.ID,
				FName:   res.FName,
				City:    res.City,
				Phone:   res.Phone,
				Height:  res.Height,
				Married: res.Married,
			}}, nil
	}
}
func makeGetUsersByIDsEndpoint(svc userservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*UsersRequest)
		res, err := svc.GetUsersByIDs(ctx, req.IDs)
		if err != nil {
			return nil, err
		}
		var users []User
		for _, user := range res {
			users = append(users,
				User{
					ID:      user.ID,
					FName:   user.FName,
					City:    user.City,
					Phone:   user.Phone,
					Height:  user.Height,
					Married: user.Married,
				})
		}
		return &UsersResponse{
			Users: users,
		}, nil
	}
}
func makeGetUsersByCriteriaEndpoint(svc userservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*SearchRequest)
		res, err := svc.GetUsersByCriteria(ctx, types.NewSearchCriteria(req.Name, req.City))
		if err != nil {
			return nil, err
		}
		var users []User
		for _, user := range res {
			users = append(users,
				User{
					ID:      user.ID,
					FName:   user.FName,
					City:    user.City,
					Phone:   user.Phone,
					Height:  user.Height,
					Married: user.Married,
				})
		}
		return &UsersResponse{
			Users: users,
		}, nil
	}
}

type UserRequest struct {
	ID int64
}
type UsersRequest struct {
	IDs []int64
}

type UserResponse struct {
	User User
}
type UsersResponse struct {
	Users []User
}
type User struct {
	ID      int64
	FName   string
	City    string
	Phone   string
	Height  float64
	Married bool
}

type SearchRequest struct {
	City string
	Name string
}
