package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/f6o/memoapp/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func main() {
	// gRPCサーバーへの接続を確立
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		getMemo(client)
	case "listMemos":
		listMemos(client)
	case "createMemo":
		createMemo(client, os.Args[2], os.Args[3])
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}

func getMemo(client pb.MemoServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.GetMemoRequest{MemoId: 1}
	res, err := client.GetMemo(ctx, req)
	if err != nil {
		log.Fatalf("could not get memo: %v", err)
	}
	log.Printf("Memo: %v", res.Memo)
}

func listMemos(client pb.MemoServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.ListMemosRequest{}
	res, err := client.ListMemos(ctx, req)
	if err != nil {
		log.Fatalf("could not list memos: %v", err)
	}
	log.Printf("Memos: %v", res.Memos)
}

func createMemo(client pb.MemoServiceClient, title, content string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.CreateMemoRequest{Title: title, Content: content}
	res, err := client.CreateMemo(ctx, req)
	if err != nil {
		log.Fatalf("could not create memo: %v", err)
	}
	log.Printf("Created Memo: %v", res.Memo)
}
