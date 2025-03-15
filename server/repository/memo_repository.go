package repository

import (
	"database/sql"
	"time"

	"github.com/f6o/memoapp/proto"
	_ "github.com/mattn/go-sqlite3"
)

type MemoRepository struct {
	db *sql.DB
}

type MemoRepositoryInterface interface {
	CreateMemo(title, content string) (*proto.Memo, error)
	GetMemo(id int64) (*proto.Memo, error)
	ListMemos() ([]*proto.Memo, error)
}

func NewMemoRepository(dbPath string) (*MemoRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// テーブルが存在しない場合に作成する
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS memos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &MemoRepository{db: db}, nil
}

func (r *MemoRepository) CreateMemo(title, content string) (*proto.Memo, error) {
	now := time.Now().Unix()
	result, err := r.db.Exec("INSERT INTO memos (title, content, created_at, updated_at) VALUES (?, ?, ?, ?)", title, content, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &proto.Memo{
		Id:        id,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r *MemoRepository) GetMemo(id int64) (*proto.Memo, error) {
	row := r.db.QueryRow("SELECT id, title, content, created_at, updated_at FROM memos WHERE id = ?", id)

	var memo proto.Memo
	err := row.Scan(&memo.Id, &memo.Title, &memo.Content, &memo.CreatedAt, &memo.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &memo, nil
}

func (r *MemoRepository) ListMemos() ([]*proto.Memo, error) {
	rows, err := r.db.Query("SELECT id, title, content, created_at, updated_at FROM memos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memos []*proto.Memo
	for rows.Next() {
		var memo proto.Memo
		if err := rows.Scan(&memo.Id, &memo.Title, &memo.Content, &memo.CreatedAt, &memo.UpdatedAt); err != nil {
			return nil, err
		}
		memos = append(memos, &memo)
	}

	return memos, nil
}
