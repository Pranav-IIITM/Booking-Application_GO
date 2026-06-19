package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)
	log.Println("Frontend server listening on http://localhost:5502")
	log.Fatal(http.ListenAndServe(":5502", nil))
}
