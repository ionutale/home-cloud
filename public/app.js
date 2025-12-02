document.addEventListener('DOMContentLoaded', () => {
    loadFiles();

    const uploadForm = document.getElementById('uploadForm');
    const fileInput = document.getElementById('fileInput');
    const uploadStatus = document.getElementById('uploadStatus');
    const dropZone = document.getElementById('dropZone');

    // Drag and Drop events
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, preventDefaults, false);
    });

    function preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    ['dragenter', 'dragover'].forEach(eventName => {
        dropZone.addEventListener(eventName, highlight, false);
    });

    ['dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, unhighlight, false);
    });

    function highlight(e) {
        dropZone.classList.add('drag-over');
    }

    function unhighlight(e) {
        dropZone.classList.remove('drag-over');
    }

    dropZone.addEventListener('drop', handleDrop, false);

    function handleDrop(e) {
        const dt = e.dataTransfer;
        const files = dt.files;
        handleFiles(files);
    }

    let uploadQueue = [];
    let isUploading = false;

    uploadForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        if (fileInput.files.length > 0) {
            handleFiles(fileInput.files);
        }
    });

    function handleFiles(files) {
        if (files.length === 0) return;
        
        // Add files to queue
        for (let i = 0; i < files.length; i++) {
            uploadQueue.push(files[i]);
        }
        
        processQueue();
    }

    async function processQueue() {
        if (isUploading || uploadQueue.length === 0) return;

        isUploading = true;
        const file = uploadQueue.shift();
        
        uploadStatus.textContent = `Uploading ${file.name}... (${uploadQueue.length} more in queue)`;
        
        const progressContainer = document.getElementById('progressContainer');
        const progressBar = document.getElementById('progressBar');
        progressContainer.style.display = 'block';
        progressBar.style.width = '0%';

        const formData = new FormData();
        formData.append('file', file);

        try {
            await new Promise((resolve, reject) => {
                const xhr = new XMLHttpRequest();
                xhr.open('POST', '/upload', true);

                xhr.upload.onprogress = (e) => {
                    if (e.lengthComputable) {
                        const percentComplete = (e.loaded / e.total) * 100;
                        progressBar.style.width = percentComplete + '%';
                    }
                };

                xhr.onload = () => {
                    if (xhr.status === 200) {
                        resolve(xhr.response);
                    } else {
                        reject(new Error('Upload failed'));
                    }
                };

                xhr.onerror = () => {
                    reject(new Error('Network error'));
                };

                xhr.send(formData);
            });

            loadFiles(); // Refresh list after each success
        } catch (error) {
            console.error(`Upload failed for ${file.name}`, error);
            uploadStatus.textContent = `Upload failed for ${file.name}`;
        } finally {
            isUploading = false;
            progressBar.style.width = '0%';
            progressContainer.style.display = 'none';
            
            if (uploadQueue.length > 0) {
                processQueue();
            } else {
                uploadStatus.textContent = 'All uploads complete!';
                fileInput.value = ''; // Clear input
            }
        }
    }
});

async function loadFiles() {
    try {
        const response = await fetch('/files');
        const files = await response.json();
        
        const fileList = document.getElementById('fileList');
        fileList.innerHTML = '';

        files.forEach(file => {
            const li = document.createElement('li');
            li.className = 'file-item';
            
            // Thumbnail
            const thumbDiv = document.createElement('div');
            thumbDiv.className = 'file-thumbnail';
            if (file.thumbnailUrl) {
                const img = document.createElement('img');
                img.src = file.thumbnailUrl;
                img.className = 'file-thumbnail';
                thumbDiv.innerHTML = '';
                thumbDiv.appendChild(img);
            } else {
                thumbDiv.textContent = file.name.split('.').pop().toUpperCase();
            }

            const nameSpan = document.createElement('div');
            nameSpan.className = 'file-name';
            nameSpan.textContent = file.name;
            
            const downloadLink = document.createElement('a');
            downloadLink.href = `/download/${encodeURIComponent(file.name)}`;
            downloadLink.textContent = 'Download';
            downloadLink.className = 'download-btn';
            downloadLink.setAttribute('download', '');

            li.appendChild(thumbDiv);
            li.appendChild(nameSpan);
            li.appendChild(downloadLink);
            fileList.appendChild(li);
        });
    } catch (error) {
        console.error('Error loading files:', error);
    }
}