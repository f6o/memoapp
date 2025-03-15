package server

import (
	"context"
	"testing"

	"github.com/f6o/memoapp/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type MockMemoRepository struct {
	mock.Mock
}

func (m *MockMemoRepository) CreateMemo(title, content string) (*proto.Memo, error) {
	args := m.Called(title, content)
	return args.Get(0).(*proto.Memo), args.Error(1)
}

func (m *MockMemoRepository) GetMemo(memoId int64) (*proto.Memo, error) {
	args := m.Called(memoId)
	return args.Get(0).(*proto.Memo), args.Error(1)
}

func (m *MockMemoRepository) ListMemos() ([]*proto.Memo, error) {
	args := m.Called()
	return args.Get(0).([]*proto.Memo), args.Error(1)
}

func TestCreateMemo(t *testing.T) {
	mockRepo := new(MockMemoRepository)
	s := NewServer(mockRepo)

	req := &proto.CreateMemoRequest{Title: "Test Title", Content: "Test Content"}
	expectedMemo := &proto.Memo{Id: 1, Title: "Test Title", Content: "Test Content"}
	mockRepo.On("CreateMemo", req.Title, req.Content).Return(expectedMemo, nil)

	resp, err := s.CreateMemo(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedMemo, resp.Memo)
	mockRepo.AssertExpectations(t)
}

func TestGetMemo(t *testing.T) {
	mockRepo := new(MockMemoRepository)
	s := NewServer(mockRepo)

	req := &proto.GetMemoRequest{MemoId: 1}
	expectedMemo := &proto.Memo{Id: 1, Title: "Test Title", Content: "Test Content"}
	mockRepo.On("GetMemo", req.MemoId).Return(expectedMemo, nil)

	resp, err := s.GetMemo(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, expectedMemo, resp.Memo)
	mockRepo.AssertExpectations(t)
}

func TestListMemos(t *testing.T) {
	mockRepo := new(MockMemoRepository)
	s := NewServer(mockRepo)

	expectedMemos := []*proto.Memo{
		{Id: 1, Title: "Test Title 1", Content: "Test Content 1"},
		{Id: 2, Title: "Test Title 2", Content: "Test Content 2"},
	}
	mockRepo.On("ListMemos").Return(expectedMemos, nil)

	resp, err := s.ListMemos(context.Background(), &proto.ListMemosRequest{})
	assert.NoError(t, err)
	assert.Equal(t, expectedMemos, resp.Memos)
	mockRepo.AssertExpectations(t)
}
