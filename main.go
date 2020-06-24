package main

import (
	"log"
	"net/http"
)

var server *Server = &Server{}

func main() {
	err := server.InitializeDB(
		DBDriver,
		DBHost,
		DBPort,
		DBUsername,
		DBPassword,
		DBName,
	)
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/login", Login)
	http.HandleFunc("/todos", Welcome)
	http.HandleFunc("/refresh", RenewJWT)
	http.HandleFunc("/logout", Logout)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
