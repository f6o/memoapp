package main

import (
	"log"
	"net"

	"github.com/f6o/memoapp/proto"
	"github.com/f6o/memoapp/server"
	"github.com/f6o/memoapp/server/repository"
	"google.golang.org/grpc"
)

func main() {
	repo, err := repository.NewMemoRepository("memos.db")
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterMemoServiceServer(s, server.NewServer(repo))

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
