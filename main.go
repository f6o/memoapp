package main

import (
	"log"
	"net"
	"os"

	"github.com/f6o/memoapp/proto"
	"github.com/f6o/memoapp/server"
	"github.com/f6o/memoapp/server/repository"
	"google.golang.org/grpc"
)

const (
	defaultDbFile        = "memos.db"
	defaultListenAddress = ":50051"
)

func main() {
	dbFile := os.Getenv("DB_FILE")
	if dbFile == "" {
		dbFile = defaultDbFile
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = defaultListenAddress
	}
	if _, _, err := net.SplitHostPort(listenAddr); err != nil {
		log.Fatalf("invalid listen address: %v", err)
	}

	repo, err := repository.NewMemoRepository(dbFile)
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}

	lis, err := net.Listen("tcp", listenAddr)
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
