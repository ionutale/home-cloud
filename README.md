# ☁️ Home Cloud

A simple, self-hosted file sharing web application built with Go. Upload, view, and manage your files through a modern web interface or via FTP.

## ✨ Features

- **📂 File Management**: Upload and download files easily.
- **🖼️ Image Thumbnails**: Automatic thumbnail generation for uploaded images.
- **🖱️ Drag & Drop**: Modern upload interface with drag-and-drop support.
- **🔌 FTP Support**: Built-in FTP server for alternative file access.
- **🐳 Docker Ready**: Easy deployment with Docker and Docker Compose.

## 🚀 Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
- *Or* [Go](https://golang.org/) 1.25+ for local development

### 🐳 Running with Docker (Recommended)

1.  Clone the repository:
    ```bash
    git clone https://github.com/ionutale/home-cloud.git
    cd home-cloud
    ```

2.  Start the application:
    ```bash
    docker-compose up -d --build
    ```

3.  Access the application:
    -   **Web Interface**: [http://localhost:8080](http://localhost:8080)
    -   **FTP Server**: `localhost:2121`

### 🛠️ Running Locally

1.  Install dependencies:
    ```bash
    go mod tidy
    ```

2.  Run the application:
    ```bash
    go run main.go
    ```

## 📖 Usage

### Web Interface
Open your browser to `http://localhost:8080`. You can drag and drop files into the upload area or click to select files. Uploaded images will automatically display thumbnails.

### FTP Access
You can connect to the file server using any FTP client (like FileZilla) or the command line.

- **Host**: `localhost`
- **Port**: `2121`
- **Username**: `admin`
- **Password**: `admin`

## 🏗️ Project Structure

```
home-cloud/
├── main.go           # Go backend (HTTP & FTP servers)
├── Dockerfile        # Docker build configuration
├── docker-compose.yml # Docker Compose configuration
├── public/           # Frontend assets
│   ├── index.html
│   ├── style.css
│   └── app.js
├── uploads/          # Directory for stored files
└── thumbnails/       # Directory for generated thumbnails
```

## 💻 Technologies

- **Backend**: Go (Golang)
- **Frontend**: HTML5, CSS3, JavaScript
- **FTP Server**: [goftp.io/server](https://github.com/goftp/server)
- **Image Processing**: [disintegration/imaging](https://github.com/disintegration/imaging)

## 📄 License

This project is open source and available under the [MIT License](LICENSE).
