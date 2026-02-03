package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("MUP service started")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
