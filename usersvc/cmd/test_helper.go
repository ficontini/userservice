package main

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"testing"

	"github.com/ficontini/user-search/types"
	"github.com/ficontini/user-search/usersvc/pkg/userendpoint"
	"github.com/ficontini/user-search/usersvc/pkg/userservice"
	"github.com/ficontini/user-search/usersvc/pkg/usertransport"
	"github.com/ficontini/user-search/usersvc/proto"
	"github.com/ficontini/user-search/usersvc/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func setupGRPCServer(t *testing.T, mockStore store.Store) (*bufconn.Listener, *grpc.Server) {
	var (
		svc        = userservice.New(mockStore)
		endpoint   = userendpoint.New(svc)
		userServer = usertransport.NewGRPCServer(endpoint)
		ln         = bufconn.Listen(bufSize)
		server     = grpc.NewServer()
	)
	proto.RegisterUserServer(server, userServer)
	go func() {
		if err := server.Serve(ln); err != nil {
			log.Fatal(err)
		}
	}()
	return ln, server
}

func setupGRPCClient(t *testing.T, ln *bufconn.Listener) (*grpc.ClientConn, userservice.Service) {
	conn, err := grpc.NewClient("passthrough:whatever", grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return ln.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	client := usertransport.NewGRPClient(conn)
	return conn, client
}

type MockStore struct {
	mu    sync.RWMutex
	Users map[int64]*types.User
}

func NewMockStore() *MockStore {
	return &MockStore{Users: make(map[int64]*types.User)}
}
func (m *MockStore) Add(_ context.Context, user *types.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Users[user.ID] = user
	return nil
}

func (m *MockStore) GetByID(_ context.Context, id int64) (*types.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	user, ok := m.Users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}
func (m *MockStore) GetByIDs(_ context.Context, ids []int64) ([]*types.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var users []*types.User
	for _, id := range ids {
		if user := m.Users[id]; user != nil {
			users = append(users, user)
		}
	}
	return users, nil
}
func (m *MockStore) SearchByCriteria(_ context.Context, criteria types.SearchCriteria) ([]*types.User, error) {
	var users []*types.User
	for _, u := range m.Users {
		if criteria.Meets(u) {
			users = append(users, u)
		}
	}
	return users, nil
}
