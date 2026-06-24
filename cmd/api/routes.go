package main

import (
	"net/http"
)


func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	
	router.HandleFunc("/hello", app.HelloHandler)
	router.HandleFunc("/videos/", app.StreamHandler)
	router.HandleFunc("/upload", app.StoreUploadedFile)
	return router
}