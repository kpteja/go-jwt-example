package main

// Database config.
const (
	DBDriver   string = "postgres"
	DBHost     string = "localhost"
	DBPort     string = "5432"
	DBUsername string = "postgres"
	DBPassword string = "cadbury"
	DBName     string = "todos"
)

// JWT config.
const (
	JWTSecret             string = "OneOfTheManyOpenSecrets"
	JWTAccessTokenExpiry  string = "1 day"
	JWTRefreshTokenExpiry string = "1 week"
)
