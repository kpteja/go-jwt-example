package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Claims represent claims to be embedded into JWT.
type Claims struct {
	UserID    uint64 `json:"user_id"`
	UserEmail string `json:"user_email"`
	jwt.StandardClaims
}

// RefreshToken represents a refresh token in DB.
type RefreshToken struct {
	ID        uint64
	Token     string
	UserID    uint64
	Expiry    time.Time
	CreatedOn time.Time
}

// CreateAccessToken creates a new JWT access token.
func CreateAccessToken(data map[string]interface{}) (string, error) {
	claims := &Claims{
		UserID:    data["userID"].(uint64),
		UserEmail: data["userEmail"].(string),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

// CreateRefreshToken creates a refresh token.
func CreateRefreshToken() string {
	uniqueID, err := uuid.NewRandom()
	if err != nil {
		log.Println(err)
	}

	return uniqueID.String()
}

// ExtractAccessToken extracts token from HTTP Authorization header.
func ExtractAccessToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

// ValidToken verifies if token is valid.
func ValidToken(r *http.Request) bool {
	tokenString := ExtractAccessToken(r)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Does this token conform to "SigningMethodHMAC"?
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JWTSecret), nil
	})
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// StoreRefreshToken stores refresh token in DB.
func (s *Server) StoreRefreshToken(refreshToken string, userID uint64) (int, error) {
	insertSQL := `
		INSERT INTO refreshTokens (refreshToken, userid, expiry)
		VALUES ($1, $2, $3)
		RETURNING id`
	id := 0
	expiry := time.Now().Local().Add(time.Minute * time.Duration(30))
	err := s.DB.QueryRow(insertSQL, refreshToken, userID, expiry).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// GetRefreshTokenByToken fetches refresh token.
func (s *Server) GetRefreshTokenByToken(token string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	selectSQL := `
		SELECT id, refreshToken, userid, expiry createdon FROM refreshTokens
		WHERE r.refreshToken = $1`
	row := s.DB.QueryRow(selectSQL, token)
	err := row.Scan(
		&refreshToken.ID,
		&refreshToken.Token,
		&refreshToken.UserID,
		&refreshToken.Expiry,
		&refreshToken.CreatedOn,
	)
	switch err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("Refresh token [%v] not found", token)
	case nil:
		return &refreshToken, nil
	default:
		return nil, err
	}
}
