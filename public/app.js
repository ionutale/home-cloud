document.addEventListener('DOMContentLoaded', () => {
    loadFiles();

    const uploadForm = document.getElementById('uploadForm');
    const fileInput = document.getElementById('fileInput');
    const uploadStatus = document.getElementById('uploadStatus');

    uploadForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        if (fileInput.files.length === 0) {
            return;
        }

        const formData = new FormData();
        formData.append('file', fileInput.files[0]);

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
    });
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
            
            const nameSpan = document.createElement('span');
            nameSpan.textContent = file.name;
            
            const downloadLink = document.createElement('a');
            downloadLink.href = `/download/${encodeURIComponent(file.name)}`;
            downloadLink.textContent = 'Download';
            downloadLink.className = 'download-btn';
            downloadLink.setAttribute('download', '');

            li.appendChild(nameSpan);
            li.appendChild(downloadLink);
            fileList.appendChild(li);
        });
    } catch (error) {
        console.error('Error loading files:', error);
    }
}