package main

import (
	"fmt"
	"net/http"
	"os"

	"./app"
	"github.com/gorilla/mux"

	"./controllers"
)

func main() {
	router := mux.NewRouter()
	router.Use(app.JwtAuthentication)

	router.HandleFunc("/api/users/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/users/login", controllers.Authenticate).Methods("POST")

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
