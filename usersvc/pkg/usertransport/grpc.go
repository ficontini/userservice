package usertransport

import (
	"context"

	"github.com/ficontini/user-search/usersvc/pkg/userendpoint"
	"github.com/ficontini/user-search/usersvc/pkg/userservice"
	"github.com/ficontini/user-search/usersvc/proto"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type grpcServer struct {
	getUserByID       grpctransport.Handler
	getUsersByIDs     grpctransport.Handler
	getUserByCriteria grpctransport.Handler
	proto.UnimplementedUserServer
}

func (s *grpcServer) GetUserByID(ctx context.Context, req *proto.UserRequest) (*proto.UserResponse, error) {
	_, rep, err := s.getUserByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.UserResponse), nil
}
func (s *grpcServer) GetUsersByIDs(ctx context.Context, req *proto.UsersRequest) (*proto.UsersResponse, error) {
	_, rep, err := s.getUsersByIDs.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.UsersResponse), nil
}
func (s *grpcServer) GetUsersByCriteria(ctx context.Context, req *proto.SearchRequest) (*proto.UsersResponse, error) {
	_, rep, err := s.getUserByCriteria.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*proto.UsersResponse), nil
}

func NewGRPCServer(endpoints userendpoint.Set) proto.UserServer {
	options := []grpctransport.ServerOption{}
	return &grpcServer{
		getUserByID: grpctransport.NewServer(
			endpoints.GetUserByIDEndpoint,
			decodeGRPCGetUserReq,
			encodeGRPCGetUserResp,
			options...,
		),
		getUsersByIDs: grpctransport.NewServer(
			endpoints.GetUsersByIDsEndpoint,
			decodeGRPCGetUsersReq,
			encodeGRPCGetUsersResp,
			options...,
		),
		getUserByCriteria: grpctransport.NewServer(
			endpoints.GetUsersByCriteriaEndpoint,
			decodeGRPCSearchRequest,
			encodeGRPCGetUsersResp,
			options...,
		),
	}
}

func decodeGRPCGetUserReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.UserRequest)
	return &userendpoint.UserRequest{ID: req.Id}, nil
}
func encodeGRPCGetUserResp(_ context.Context, resp interface{}) (interface{}, error) {
	response := resp.(*userendpoint.UserResponse)
	return &proto.UserResponse{
		Id:      response.User.ID,
		Fname:   response.User.FName,
		City:    response.User.City,
		Phone:   response.User.Phone,
		Height:  response.User.Height,
		Married: response.User.Married,
	}, nil
}
func decodeGRPCGetUsersReq(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.UsersRequest)
	return &userendpoint.UsersRequest{
		IDs: req.Ids,
	}, nil
}
func decodeGRPCSearchRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*proto.SearchRequest)
	return &userendpoint.SearchRequest{
		Name: req.Name,
		City: req.City,
	}, nil
}
func encodeGRPCGetUsersResp(_ context.Context, resp interface{}) (interface{}, error) {
	response := resp.(*userendpoint.UsersResponse)
	var users []*proto.UserResponse
	for _, r := range response.Users {
		users = append(users, &proto.UserResponse{
			Id:      r.ID,
			Fname:   r.FName,
			City:    r.City,
			Phone:   r.Phone,
			Height:  r.Height,
			Married: r.Married,
		})
	}
	return &proto.UsersResponse{
		Users: users,
	}, nil
}

func NewGRPClient(conn *grpc.ClientConn) userservice.Service {
	var (
		options                   = []grpctransport.ClientOption{}
		getUserByIDEndpoint       endpoint.Endpoint
		getUsersByIDsEndpoint     endpoint.Endpoint
		getUserByCriteriaEndpoint endpoint.Endpoint
	)
	getUserByIDEndpoint = grpctransport.NewClient(
		conn,
		"User",
		"GetUserByID",
		encodeGRPGetUserByIDRequest,
		decodeGRPGetUserByIDResponse,
		proto.UserResponse{},
		options...,
	).Endpoint()
	getUsersByIDsEndpoint = grpctransport.NewClient(
		conn,
		"User",
		"GetUsersByIDs",
		encodeGRPGetUsersByIDsRequest,
		decodeGRPGetUsersResponse,
		proto.UsersResponse{},
		options...,
	).Endpoint()
	getUserByCriteriaEndpoint = grpctransport.NewClient(
		conn,
		"User",
		"GetUsersByCriteria",
		encodeGRPGetUsersByCriteriaRequest,
		decodeGRPGetUsersResponse,
		proto.UsersResponse{},
		options...,
	).Endpoint()

	return userendpoint.Set{
		GetUserByIDEndpoint:        getUserByIDEndpoint,
		GetUsersByIDsEndpoint:      getUsersByIDsEndpoint,
		GetUsersByCriteriaEndpoint: getUserByCriteriaEndpoint,
	}

}
func encodeGRPGetUserByIDRequest(ctx context.Context, req interface{}) (interface{}, error) {
	request := req.(*userendpoint.UserRequest)
	return &proto.UserRequest{
		Id: request.ID,
	}, nil
}

func decodeGRPGetUserByIDResponse(_ context.Context, resp interface{}) (interface{}, error) {
	response := resp.(*proto.UserResponse)
	return &userendpoint.UserResponse{
		User: userendpoint.User{
			ID:      response.Id,
			FName:   response.Fname,
			City:    response.City,
			Phone:   response.Phone,
			Height:  response.Height,
			Married: response.Married,
		},
	}, nil
}
func encodeGRPGetUsersByIDsRequest(ctx context.Context, req interface{}) (interface{}, error) {
	request := req.(*userendpoint.UsersRequest)
	return &proto.UsersRequest{
		Ids: request.IDs,
	}, nil
}
func encodeGRPGetUsersByCriteriaRequest(ctx context.Context, req interface{}) (interface{}, error) {
	request := req.(*userendpoint.SearchRequest)
	return &proto.SearchRequest{
		Name: request.Name,
		City: request.City,
	}, nil
}

func decodeGRPGetUsersResponse(_ context.Context, resp interface{}) (interface{}, error) {
	response := resp.(*proto.UsersResponse)
	var users []userendpoint.User
	for _, r := range response.Users {
		users = append(users, userendpoint.User{
			ID:      r.Id,
			FName:   r.Fname,
			City:    r.City,
			Phone:   r.Phone,
			Height:  r.Height,
			Married: r.Married,
		})
	}
	return &userendpoint.UsersResponse{
		Users: users,
	}, nil
}
