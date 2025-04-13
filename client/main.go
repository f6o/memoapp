package main

import (
	"context"
	"flag" // Import the flag package
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/f6o/memoapp/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure" // Import for insecure connections
)

const (
	fallbackAddress = "localhost:50051"
)

func main() {
	// Define a boolean flag for development mode (insecure connection)
	devMode := flag.Bool("dev", false, "Enable insecure connection for development")
	flag.Parse() // Parse the command-line flags

	// gRPCサーバーへの接続を確立
	address := os.Getenv("MEMOAPP_SERVER_ADDRESS")
	if address == "" {
		address = fallbackAddress
	}

	var opts []grpc.DialOption
	if *devMode {
		// Use insecure connection if --dev flag is set
		log.Println("WARNING: Connecting insecurely (--dev mode)")
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// Use TLS credentials by default
		certFile := os.Getenv("MEMOAPP_SERVER_CERTFILE")
		if certFile == "" {
			log.Fatalf("MEMOAPP_SERVER_CERTFILE environment variable must be set for secure connection")
		}
		cred, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Fatalf("could not load tls cert '%s': %v", certFile, err)
		}
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	conn, err := grpc.NewClient(address, opts...) // Pass the options slice
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMemoServiceClient(conn)

	// Use flag.Args() to get non-flag arguments
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalf("usage: %s [--dev] [getMemo|listMemos|createMemo] [arguments...]", os.Args[0])
	}

	command := args[0]
	switch command {
	case "getMemo":
		if len(args) < 2 {
			log.Fatalf("usage: %s [--dev] getMemo <memoId>", os.Args[0])
		}
		memoId, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Fatalf("invalid memoId: %v", err)
		}
		getMemo(client, memoId)
	case "listMemos":
		listMemos(client)
	case "createMemo":
		if len(args) < 3 {
			log.Fatalf("usage: %s [--dev] createMemo <title> <content>", os.Args[0])
		}
		createMemo(client, args[1], args[2])
	default:
		log.Fatalf("unknown command: %s", command)
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
