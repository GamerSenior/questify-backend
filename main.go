package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/GamerSenior/questify-backend/app"
	_ "github.com/GamerSenior/questify-backend/controllers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Use(JwtAuthentication)

	router.HandleFunc("/api/users/new", CreateAccount).Methods("POST")
	router.HandleFunc("/api/users/login", Authenticate).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		fmt.Print(err)
	}
}
