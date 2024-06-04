package main

import (
	"context"
	"math/rand"
	"testing"

	"github.com/ficontini/user-search/types"
	"github.com/ficontini/user-search/usersvc/store"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	var (
		mockStore = NewMockStore()
		user      = store.GenerateRandomUser(1)
	)
	mockStore.Add(context.Background(), user)

	ln, server := setupGRPCServer(t, mockStore)
	defer server.Stop()

	conn, client := setupGRPCClient(t, ln)
	defer conn.Close()

	res, err := client.GetUserByID(context.TODO(), user.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, user.ID, res.ID)
	assert.Equal(t, user.FName, res.FName)
}
func TestGetUserNotFound(t *testing.T) {
	var (
		mockStore = NewMockStore()
	)

	ln, server := setupGRPCServer(t, mockStore)
	defer server.Stop()

	conn, client := setupGRPCClient(t, ln)
	defer conn.Close()

	res, err := client.GetUserByID(context.TODO(), 1)
	if err != nil {
		assert.Contains(t, err.Error(), "user not found")
	}
	assert.Nil(t, res)

}
func TestGetUsersByIDs(t *testing.T) {
	var (
		ctx       = context.Background()
		mockStore = NewMockStore()
		user1     = store.GenerateRandomUser(1)
		user2     = store.GenerateRandomUser(2)
		ids       = []int64{user1.ID, user2.ID}
	)

	mockStore.Add(ctx, user1)
	mockStore.Add(ctx, user2)

	ln, server := setupGRPCServer(t, mockStore)
	defer server.Stop()

	conn, client := setupGRPCClient(t, ln)
	defer conn.Close()

	res, err := client.GetUsersByIDs(ctx, ids)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(res))

	for i, id := range ids {
		assert.Equal(t, id, res[i].ID)
	}

}

func TestGetUsersByCritera(t *testing.T) {
	var (
		ctx       = context.Background()
		mockStore = NewMockStore()
		users     = []*types.User{
			types.NewUser(1, "Steve", "LA", "123456990", rand.Float64()+5, true),
			types.NewUser(2, "Bob", "LA", "123453990", rand.Float64()+5, false),
			types.NewUser(3, "Fin", "San Antonio", "111456990", rand.Float64()+5, true),
			types.NewUser(4, "Paul", "San Diego", "123986990", rand.Float64()+5, false),
			types.NewUser(5, "Steve", "Dallas", "123234990", rand.Float64()+5, true),
		}
	)
	for _, u := range users {
		mockStore.Add(ctx, u)
	}

	ln, server := setupGRPCServer(t, mockStore)
	defer server.Stop()

	conn, client := setupGRPCClient(t, ln)
	defer conn.Close()

	res, err := client.GetUsersByCriteria(ctx, types.NewSearchCriteria("Steve", "LA"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(res))

	res, err = client.GetUsersByCriteria(ctx, types.NewSearchCriteria("", "Dallas"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(res))

	res, err = client.GetUsersByCriteria(ctx, types.NewSearchCriteria("", ""))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 5, len(res))
}
