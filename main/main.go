package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/omarmendozaaa/backend-go/server"
)

func main() {
	s := server.New()
	fmt.Println("Welcome . . .")
	log.Fatal(http.ListenAndServe(":8080", s.Router()))
}
