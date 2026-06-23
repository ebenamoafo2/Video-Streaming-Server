package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func storeUploadedFile(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the request body to prevent large uploads
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) 

	// Parse the multipart form data
	if err := r.ParseMultipartForm(5 << 20); err != nil { 
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.MultipartForm.RemoveAll() 

	// Retrieve the uploaded file from the form data
	uploadedFile, uploadedFileHeader, err := r.FormFile("media")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer uploadedFile.Close()

	// Create the destination file in the "static/videos" directory
	f, err := os.Create(filepath.Join("static", "videos", uploadedFileHeader.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Copy the uploaded file's content to the destination file
	_, err = io.Copy(f, uploadedFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


func streamHandler(w http.ResponseWriter, r *http.Request) {
	videoName := r.URL.Path[len("/videos/"):]
	videoPath := filepath.Join("static", "videos", videoName)

	file, _ := os.Open(videoPath)

	modTime := time.Now()
	http.ServeContent(w, r, videoName, modTime, file)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/videos/", streamHandler)
	http.HandleFunc("/upload", storeUploadedFile)
	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}