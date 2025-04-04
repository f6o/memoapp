package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/f6o/memoapp/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	fallbackAddress = "localhost:50051"
)

func main() {
	// gRPCサーバーへの接続を確立
	address := os.Getenv("MEMOAPP_SERVER_ADDRESS")
	if address == "" {
		address = fallbackAddress
	}

	cred, err := credentials.NewClientTLSFromFile(os.Getenv("MEMOAPP_SERVER_CERTFILE"), "")
	if err != nil {
		log.Fatalf("could not load tls cert: %v", err)
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMemoServiceClient(conn)

	if len(os.Args) < 2 {
		log.Fatalf("usage: %s [getMemo|listMemos|createMemo]", os.Args[0])
	}

	switch os.Args[1] {
	case "getMemo":
		memoId, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			log.Fatalf("invalid memoId: %v", err)
		}
		getMemo(client, memoId)
	case "listMemos":
		listMemos(client)
	case "createMemo":
		createMemo(client, os.Args[2], os.Args[3])
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}

func getMemo(client pb.MemoServiceClient, memoId int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetMemoRequest{MemoId: memoId}
	res, err := client.GetMemo(ctx, req)
	if err != nil {
		log.Fatalf("could not get memo: %v", err)
	}
	log.Printf("Memo: %v", res.Memo)
}

func listMemos(client pb.MemoServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.ListMemosRequest{}
	res, err := client.ListMemos(ctx, req)
	if err != nil {
		log.Fatalf("could not list memos: %v", err)
	}
	log.Printf("Memos: %v", res.Memos)
}

func createMemo(client pb.MemoServiceClient, title, content string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.CreateMemoRequest{Title: title, Content: content}
	res, err := client.CreateMemo(ctx, req)
	if err != nil {
		log.Fatalf("could not create memo: %v", err)
	}
	log.Printf("Created Memo: %v", res.Memo)
}
