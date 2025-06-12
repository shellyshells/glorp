let currentPostType = 'text';
let uploadedImageUrl = '';
let isUploading = false;

// Post type tab switching
document.querySelectorAll('.post-type-tab').forEach(tab => {
    tab.addEventListener('click', function() {
        // Remove active class from all tabs
        document.querySelectorAll('.post-type-tab').forEach(t => t.classList.remove('active'));
        
        // Add active class to clicked tab
        this.classList.add('active');
        
        // Hide all content areas
        document.getElementById('text-content').style.display = 'none';
        document.getElementById('link-content').style.display = 'none';
        document.getElementById('image-content').style.display = 'none';
        
        // Show corresponding content area
        const type = this.dataset.type;
        currentPostType = type;
        document.getElementById(type + '-content').style.display = 'block';
        
        // Reset upload state when switching away from image
        if (type !== 'image') {
            resetImageUpload();
        }
    });
});

// Tag selection
function toggleTag(tagElement) {
    const checkbox = tagElement.querySelector('input[type="checkbox"]');
    checkbox.checked = !checkbox.checked;
    tagElement.classList.toggle('selected', checkbox.checked);
}

// Initialize tag states
document.querySelectorAll('.tag-option input[type="checkbox"]').forEach(checkbox => {
    checkbox.addEventListener('change', function() {
        this.closest('.tag-option').classList.toggle('selected', this.checked);
    });
});

// Image upload functionality
const imageUpload = document.getElementById('image-upload');
const uploadArea = document.getElementById('upload-area');
const imagePreviewSection = document.getElementById('image-preview-section');
const imagePreview = document.getElementById('image-preview');

// Handle file selection
imageUpload.addEventListener('change', function(e) {
    const file = e.target.files[0];
    if (file) {
        handleImageUpload(file);
    }
});

// Drag and drop functionality
uploadArea.addEventListener('dragover', function(e) {
    e.preventDefault();
    this.classList.add('drag-over');
});

uploadArea.addEventListener('dragleave', function(e) {
    e.preventDefault();
    this.classList.remove('drag-over');
});

uploadArea.addEventListener('drop', function(e) {
    e.preventDefault();
    this.classList.remove('drag-over');
    
    const files = e.dataTransfer.files;
    if (files.length > 0) {
        const file = files[0];
        if (file.type.startsWith('image/')) {
            imageUpload.files = files;
            handleImageUpload(file);
        } else {
            showError('Please select an image file (JPG, PNG, GIF, WEBP)');
        }
    }
});

// Handle image upload
async function handleImageUpload(file) {
    // Validate file size (10MB)
    if (file.size > 10 * 1024 * 1024) {
        showError('File size must be less than 10MB');
        return;
    }

    // Validate file type
    const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
    if (!validTypes.includes(file.type)) {
        showError('Please select a valid image file (JPG, PNG, GIF, WEBP)');
        return;
    }

    isUploading = true;
    updateUploadArea('Uploading...', true);

    try {
        const formData = new FormData();
        formData.append('image', file);

        const response = await fetch('/api/upload/image', {
            method: 'POST',
            body: formData
        });

        if (response.ok) {
            const result = await response.json();
            uploadedImageUrl = result.url;
            showImagePreview(result.url);
            showSuccess('Image uploaded successfully!');
        } else {
            const error = await response.text();
            showError('Upload failed: ' + error);
            resetImageUpload();
        }
    } catch (error) {
        console.error('Upload error:', error);
        showError('Upload failed. Please try again.');
        resetImageUpload();
    } finally {
        isUploading = false;
    }
}

function updateUploadArea(text, loading = false) {
    const uploadArea = document.getElementById('upload-area');
    const icon = loading ? 'fas fa-spinner fa-spin' : 'fas fa-cloud-upload-alt';
    
    uploadArea.innerHTML = `
        <i class="${icon}"></i>
        <p>${text}</p>
        ${!loading ? '<button type="button" class="btn btn-outline" onclick="document.getElementById(\'image-upload\').click()">Choose Image</button>' : ''}
        <small>Max 10MB • JPG, PNG, GIF, WEBP</small>
    `;
}

function showImagePreview(imageUrl) {
    imagePreview.src = imageUrl;
    uploadArea.style.display = 'none';
    imagePreviewSection.style.display = 'block';
}

function removeImage() {
    uploadedImageUrl = '';
    resetImageUpload();
}

function resetImageUpload() {
    uploadArea.style.display = 'block';
    imagePreviewSection.style.display = 'none';
    imageUpload.value = '';
    updateUploadArea('Drag and drop images here or click to upload');
}

// Preview functionality
function previewPost() {
    const title = document.querySelector('.title-input').value;
    const selectedTags = Array.from(document.querySelectorAll('.tag-option.selected label')).map(label => label.textContent);
    
    let content = '';
    let previewData = {
        title: title || 'Untitled Post',
        type: currentPostType,
        tags: selectedTags
    };
    
    switch(currentPostType) {
        case 'text':
            content = document.querySelector('textarea[name="description"]').value;
            previewData.content = content || 'No content';
            break;
        case 'link':
            const linkUrl = document.querySelector('input[name="link_url"]').value;
            const linkDesc = document.querySelector('textarea[name="link_description"]').value;
            previewData.linkUrl = linkUrl;
            previewData.content = linkDesc || 'No description';
            break;
        case 'image':
            const imageDesc = document.querySelector('textarea[name="image_description"]').value;
            previewData.imageUrl = uploadedImageUrl;
            previewData.content = imageDesc || 'No description';
            break;
    }
    
    // Simple preview modal (you can enhance this)
    const previewModal = document.createElement('div');
    previewModal.className = 'preview-modal';
    previewModal.innerHTML = `
        <div class="preview-content">
            <div class="preview-header">
                <h3>Post Preview</h3>
                <button onclick="this.closest('.preview-modal').remove()">&times;</button>
            </div>
            <div class="preview-body">
                <h2>${previewData.title}</h2>
                <div class="preview-tags">
                    ${previewData.tags.map(tag => `<span class="tag">${tag}</span>`).join('')}
                </div>
                <div class="preview-type">${previewData.type.toUpperCase()} POST</div>
                ${previewData.linkUrl ? `<div class="preview-link"><strong>Link:</strong> ${previewData.linkUrl}</div>` : ''}
                ${previewData.imageUrl ? `<div class="preview-image"><img src="${previewData.imageUrl}" alt="Preview" style="max-width: 100%; height: auto; border-radius: 4px;"></div>` : ''}
                <div class="preview-content-text">${previewData.content}</div>
            </div>
        </div>
    `;
    
    document.body.appendChild(previewModal);
}

// Form submission
document.getElementById('create-post-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    if (isUploading) {
        showError('Please wait for image upload to complete');
        return;
    }
    
    const formData = new FormData(this);
    const selectedTags = Array.from(formData.getAll('tags')).map(id => parseInt(id));
    
    let description = '';
    let imageUrl = '';
    let linkUrl = '';
    
    switch(currentPostType) {
        case 'text':
            description = formData.get('description') || '';
            break;
        case 'link':
            linkUrl = formData.get('link_url') || '';
            description = formData.get('link_description') || '';
            break;
        case 'image':
            imageUrl = uploadedImageUrl;
            description = formData.get('image_description') || '';
            break;
    }
    
    const data = {
        title: formData.get('title'),
        description: description,
        tags: selectedTags,
        post_type: currentPostType,
        image_url: imageUrl,
        link_url: linkUrl
    };
    
    const submitBtn = document.querySelector('.post-submit-btn');
    const originalText = submitBtn.innerHTML;
    
    // Show loading state
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<span class="loading-spinner"></span> Posting...';
    
    try {
        const response = await fetch('/api/threads', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });
        
        if (response.ok) {
            const result = await response.json();
            showSuccess('Post created successfully! Redirecting...');
            
            // Redirect after a short delay
            setTimeout(() => {
                window.location.href = `/threads/${result.thread.id}`;
            }, 1500);
        } else {
            const result = await response.text();
            showError(result);
        }
    } catch (error) {
        showError('Failed to create post. Please try again.');
        console.error('Error creating post:', error);
    } finally {
        // Reset button state
        submitBtn.disabled = false;
        submitBtn.innerHTML = originalText;
    }
});

// Utility functions
function showError(message) {
    const errorDiv = document.getElementById('error-message');
    const successDiv = document.getElementById('success-message');
    
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
    successDiv.style.display = 'none';
    
    // Auto-hide after 5 seconds
    setTimeout(() => {
        errorDiv.style.display = 'none';
    }, 5000);
}

function showSuccess(message) {
    const errorDiv = document.getElementById('error-message');
    const successDiv = document.getElementById('success-message');
    
    successDiv.textContent = message;
    successDiv.style.display = 'block';
    errorDiv.style.display = 'none';
}

// Auto-resize textarea
document.querySelectorAll('.content-textarea').forEach(textarea => {
    textarea.addEventListener('input', function() {
        this.style.height = 'auto';
        this.style.height = Math.max(120, this.scrollHeight) + 'px';
    });
});

// Character counter for title
const titleInput = document.querySelector('.title-input');
titleInput.addEventListener('input', function() {
    const remaining = 200 - this.value.length;
    // You can add a character counter display here if needed
});

// Keyboard shortcuts
document.addEventListener('keydown', function(e) {
    // Ctrl+Enter to submit
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        e.preventDefault();
        document.getElementById('create-post-form').dispatchEvent(new Event('submit'));
    }
});