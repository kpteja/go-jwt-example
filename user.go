package main

import (
	"database/sql"
	"fmt"
	"time"
)

// User represents a user in DB.
type User struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedOn time.Time `json:"created_on"`
}

// GetUserByEmail fetches user info by email.
func (s *Server) GetUserByEmail(email string) (*User, error) {
	var user User
	selectSQL := `SELECT id, email, password, createdon FROM users WHERE email = $1`
	row := s.DB.QueryRow(selectSQL, email)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedOn)
	switch err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("Email [%v] not found", email)
	case nil:
		return &user, nil
	default:
		return nil, err
	}
}

// GetUserByRefreshToken fetches user info by refresh token.
func (s *Server) GetUserByRefreshToken(refreshToken string) (*User, error) {
	var user User
	selectSQL := `
		SELECT u.id, u.email, u.createdon FROM users u
		JOIN refreshtokens r ON u.id = r.userid
		WHERE r.refreshtoken = $1`
	row := s.DB.QueryRow(selectSQL, refreshToken)
	err := row.Scan(&user.ID, &user.Email, &user.CreatedOn)
	switch err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("Refresh token [%v] not found", refreshToken)
	case nil:
		return &user, nil
	default:
		return nil, err
	}
}
