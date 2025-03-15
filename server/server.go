package server

import (
	"context"

	"github.com/f6o/memoapp/proto"
)

type server struct {
	proto.UnimplementedMemoServiceServer
}

func (s *server) GetMemo(context.Context, *proto.GetMemoRequest) (*proto.GetMemoResponse, error) {
	return nil, nil
}

func (s *server) ListMemos(context.Context, *proto.ListMemosRequest) (*proto.ListMemosResponse, error) {
	return nil, nil
}

func (s *server) CreateMemo(context.Context, *proto.CreateMemoRequest) (*proto.CreateMemoResponse, error) {
	return nil, nil
}
