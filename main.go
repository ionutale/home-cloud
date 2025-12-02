package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"goftp.io/server/v2"
	"goftp.io/server/v2/driver/file"
)

const (
	uploadDir    = "./uploads"
	thumbnailDir = "./thumbnails"
	httpPort     = ":8080"
	ftpPort      = 2121
)

type FileInfo struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	ModTime      string `json:"modTime"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
}

func main() {
	// Ensure upload and thumbnail directories exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(thumbnailDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Start FTP Server in a goroutine
	go startFTPServer()

	// HTTP Server
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/thumbnails/", http.StripPrefix("/thumbnails/", http.FileServer(http.Dir(thumbnailDir))))
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

	// Generate thumbnail if it's an image
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
		go generateThumbnail(dstPath, header.Filename)
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

			fileInfo := FileInfo{
				Name:    entry.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime().Format(time.RFC3339),
			}

			// Check if thumbnail exists
			thumbPath := filepath.Join(thumbnailDir, entry.Name())
			if _, err := os.Stat(thumbPath); err == nil {
				fileInfo.ThumbnailUrl = "/thumbnails/" + entry.Name()
			}

			files = append(files, fileInfo)
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

func generateThumbnail(srcPath, filename string) {
	src, err := imaging.Open(srcPath)
	if err != nil {
		log.Printf("Error opening image for thumbnail: %v", err)
		return
	}

	// Resize to 100x100, preserving aspect ratio
	dst := imaging.Thumbnail(src, 100, 100, imaging.Lanczos)

	dstPath := filepath.Join(thumbnailDir, filename)
	if err := imaging.Save(dst, dstPath); err != nil {
		log.Printf("Error saving thumbnail: %v", err)
	}
}
