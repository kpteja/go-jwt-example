package main

import (
	"log"
	"time"
)

// Todo represents a todo in DB.
type Todo struct {
	ID        uint64    `json:"id"`
	Todo      string    `json:"todo"`
	UserID    uint64    `json:"user_id"`
	CreatedOn time.Time `json:"created_on"`
}

// GetAllTodos fetches all todos.
func (s *Server) GetAllTodos() ([]Todo, error) {
	todos := []Todo{}
	selectSQL := `SELECT id, todo, userid, createdon FROM todos ORDER BY id DESC`
	rows, err := s.DB.Query(selectSQL)
	if err != nil {
		return todos, err
	}
	defer rows.Close()

	for rows.Next() {
		todo := Todo{}
		err = rows.Scan(&todo.ID, &todo.Todo, &todo.UserID, &todo.CreatedOn)
		if err != nil {
			log.Println(err)
		}
		todos = append(todos, todo)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
	}

	return todos, err
}
