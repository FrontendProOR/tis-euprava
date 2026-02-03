package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("ZZS service started")

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal(err)
	}
}
