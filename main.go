package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"goftp.io/server/v2"
	"goftp.io/server/v2/driver/file"
)

const (
	uploadDir = "./uploads"
	httpPort  = ":8080"
	ftpPort   = 2121
)

type FileInfo struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
}

func main() {
	// Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Start FTP Server in a goroutine
	go startFTPServer()

	// HTTP Server
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/files", handleListFiles)
	http.HandleFunc("/download/", handleDownload)

	fmt.Printf("HTTP Server starting on http://localhost%s\n", httpPort)
	fmt.Printf("FTP Server starting on localhost:%d\n", ftpPort)
	if err := http.ListenAndServe(httpPort, nil); err != nil {
		log.Fatal(err)
	}
}

func startFTPServer() {
	driver, err := file.NewDriver(uploadDir)
	if err != nil {
		log.Fatal(err)
	}

	// Simple authentication: admin/admin
	auth := &server.SimpleAuth{
		Name:     "admin",
		Password: "admin",
	}

	s, err := server.NewServer(&server.Options{
		Driver: driver,
		Auth:   auth,
		Port:   ftpPort,
		Perm:   server.NewSimplePerm("admin", "admin"),
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := s.ListenAndServe(); err != nil {
		log.Printf("FTP Server error: %v", err)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form, 10MB max
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too big", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create destination file
	dstPath := filepath.Join(uploadDir, filepath.Base(header.Filename))
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy content
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully")
}

func handleListFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entries, err := os.ReadDir(uploadDir)
	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	var files []FileInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			files = append(files, FileInfo{
				Name:    entry.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime().Format(time.RFC3339),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := r.URL.Path[len("/download/"):]
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	// Prevent directory traversal by only using the base filename
	filename = filepath.Base(filename)

	filePath := filepath.Join(uploadDir, filename)
	http.ServeFile(w, r, filePath)
}
