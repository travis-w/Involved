package main

import (
	"net/http"
	"./api"
)

func main() {
	api.Init()
	http.HandleFunc("/api/", api.Request)
	http.Handle("/", http.FileServer(http.Dir("frontend")))
    http.ListenAndServe(":8080", nil)
}