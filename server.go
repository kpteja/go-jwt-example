package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Postgres library
)

// Server contains necessary modules.
type Server struct {
	DB *sql.DB
}

// InitializeDB connects to a DB.
func (s *Server) InitializeDB(driver, host, port, user, password, database string) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		database,
	)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}

	s.DB = db
	return nil
}
