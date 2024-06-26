package service

import (
	"context"
	"database/sql"
	"log"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	result, err := s.db.ExecContext(ctx, insert, subject, description)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	row := s.db.QueryRowContext(ctx, confirm, id)

	var todo model.TODO
	todo.ID = id
	err = row.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	
	return &todo, err
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	if size <= 0  {
		return []*model.TODO{}, nil
	}

	var rows *sql.Rows
	var err error

	if prevID == 0{
		rows,err = s.db.QueryContext(ctx, read, size)
	} else {
		rows,err = s.db.QueryContext(ctx, readWithID, prevID, size)
	}
	// Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	
	todos := []*model.TODO{}
	for rows.Next() {
		var todo model.TODO
		err = rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		todos = append(todos, &todo)
	}

	return todos, err
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	result, err := s.db.ExecContext(ctx, update, subject, description, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if rows == 0 {
		return nil, &model.ErrNotFound{}
	}
	row := s.db.QueryRowContext(ctx, confirm, id)

	var todo model.TODO
	todo.ID = id
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)

	return &todo, err
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
