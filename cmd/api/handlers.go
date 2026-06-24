package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (app *application) StoreUploadedFile(w http.ResponseWriter, r *http.Request) {
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


func (app *application) StreamHandler(w http.ResponseWriter, r *http.Request) {
	videoName := r.URL.Path[len("/videos/"):]
	videoPath := filepath.Join("static", "videos", videoName)

	file, err := os.Open(videoPath)
	if err != nil {
        if os.IsExist(err) {
			http.Error(w, "Video not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to open video", http.StatusInternalServerError)
		}
    }
	defer file.Close()

	modTime := time.Now()
	http.ServeContent(w, r, videoName, modTime, file)
}

func (app *application) HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}
