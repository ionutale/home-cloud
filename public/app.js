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

    uploadForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        if (fileInput.files.length > 0) {
            handleFiles(fileInput.files);
        }
    });

    async function handleFiles(files) {
        if (files.length === 0) return;

        const formData = new FormData();
        formData.append('file', files[0]); // Currently handling single file upload

        uploadStatus.textContent = 'Uploading...';

        try {
            const response = await fetch('/upload', {
                method: 'POST',
                body: formData
            });

            if (response.ok) {
                uploadStatus.textContent = 'Upload successful!';
                fileInput.value = ''; // Clear input
                loadFiles(); // Refresh list
            } else {
                uploadStatus.textContent = 'Upload failed.';
            }
        } catch (error) {
            console.error('Error:', error);
            uploadStatus.textContent = 'An error occurred.';
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