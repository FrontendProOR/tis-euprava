package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("SSO service started")

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}
