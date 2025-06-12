// Tab switching
document.querySelectorAll('.profile-tab').forEach(tab => {
    tab.addEventListener('click', function() {
        switchTab(this.dataset.tab);
    });
});

function switchTab(tabName) {
    // Remove active class from all tabs and content
    document.querySelectorAll('.profile-tab').forEach(t => t.classList.remove('active'));
    document.querySelectorAll('.tab-content').forEach(c => {
        c.classList.remove('active');
        c.style.display = 'none';
    });
    
    // Add active class to selected tab and content
    document.querySelector(`[data-tab="${tabName}"]`).classList.add('active');
    const content = document.getElementById(tabName + '-tab');
    content.classList.add('active');
    content.style.display = 'block';
}

// Avatar modal functions
function showAvatarModal() {
    document.getElementById('avatar-modal').style.display = 'flex';
}

function closeAvatarModal() {
    document.getElementById('avatar-modal').style.display = 'none';
}

function handleAvatarUpload(input) {
    if (input.files && input.files[0]) {
        const file = input.files[0];
        const formData = new FormData();
        formData.append('avatar', file);
        
        // Show loading
        const uploadArea = document.querySelector('.upload-area');
        uploadArea.innerHTML = '<i class="fas fa-spinner fa-spin"></i><p>Uploading...</p>';
        
        fetch('/api/profile/avatar', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                location.reload(); // Refresh to show new avatar
            } else {
                alert('Failed to upload avatar: ' + data.message);
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to upload avatar');
        })
        .finally(() => {
            closeAvatarModal();
        });
    }
}

function selectAvatarStyle(style) {
    fetch('/api/profile/avatar-style', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ style: style })
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            location.reload(); // Refresh to show new avatar
        } else {
            alert('Failed to update avatar style');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Failed to update avatar style');
    })
    .finally(() => {
        closeAvatarModal();
    });
}

// Placeholder functions for future features
function sendMessage() {
    alert('Private messaging feature coming soon!');
}

function followUser() {
    alert('Follow feature coming soon!');
}

// Helper function for array slicing in templates
function slice(arr, start, end) {
    return arr.slice(start, end);
}

// Close modal when clicking outside
document.addEventListener('click', function(e) {
    const modal = document.getElementById('avatar-modal');
    if (e.target === modal) {
        closeAvatarModal();
    }
});