package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/ficontini/user-search/usersvc/pkg/userendpoint"
	"github.com/ficontini/user-search/usersvc/pkg/userservice"
	"github.com/ficontini/user-search/usersvc/pkg/usertransport"
	"github.com/ficontini/user-search/usersvc/proto"
	"github.com/ficontini/user-search/usersvc/store"
	"google.golang.org/grpc"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3004", "listen address")
	flag.Parse()
	var (
		store       = store.NewInMemoryStore()
		svc         = userservice.New(store)
		endpoint    = userendpoint.New(svc)
		grpcHandler = usertransport.NewGRPCServer(endpoint)
	)
	ln, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	server := grpc.NewServer(grpc.EmptyServerOption{})
	proto.RegisterUserServer(server, grpcHandler)
	if err := server.Serve(ln); err != nil {
		log.Fatal(err)
	}
}
