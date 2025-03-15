package server

import (
	"context"

	"github.com/f6o/memoapp/proto"
	"github.com/f6o/memoapp/repository"
)

type server struct {
	proto.UnimplementedMemoServiceServer
	repo *repository.MemoRepository
}

func NewServer(repo *repository.MemoRepository) *server {
	return &server{repo: repo}
}

func (s *server) CreateMemo(ctx context.Context, req *proto.CreateMemoRequest) (*proto.CreateMemoResponse, error) {
	memo, err := s.repo.CreateMemo(req.Title, req.Content)
	if err != nil {
		return nil, err
	}
	return &proto.CreateMemoResponse{Memo: memo}, nil
}

func (s *server) GetMemo(ctx context.Context, req *proto.GetMemoRequest) (*proto.GetMemoResponse, error) {
	memo, err := s.repo.GetMemo(req.MemoId)
	if err != nil {
		return nil, err
	}
	return &proto.GetMemoResponse{Memo: memo}, nil
}

func (s *server) ListMemos(ctx context.Context, req *proto.ListMemosRequest) (*proto.ListMemosResponse, error) {
	memos, err := s.repo.ListMemos()
	if err != nil {
		return nil, err
	}
	return &proto.ListMemosResponse{Memos: memos}, nil
}
