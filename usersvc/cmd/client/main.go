package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ficontini/user-search/usersvc/pkg/usertransport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const grpcAddr = "localhost:3004"

func main() {
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	svc := usertransport.NewGRPClient(conn)
	user, err := svc.GetUserByID(context.Background(), 6)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", user)

	users, err := svc.GetUsersByIDs(context.TODO(), []int64{2, 3, 5})
	if err != nil {
		log.Fatal(err)
	}
	for _, u := range users {
		fmt.Printf("%+v\n", u)
	}

}
