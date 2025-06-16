let currentShareThreadId = null;

// Advanced Image modal functions with zoom and pan
let currentZoom = 1;
let isDragging = false;
let startX, startY, scrollLeft, scrollTop;

function openImageModal(imageSrc) {
    const modal = document.getElementById('image-modal');
    const modalImage = document.getElementById('modal-image');
    const container = document.getElementById('image-container');
    
    modalImage.src = imageSrc;
    modal.style.display = 'flex';
    document.body.style.overflow = 'hidden';
    
    // Reset zoom and position
    currentZoom = 1;
    updateZoom();
    container.scrollLeft = 0;
    container.scrollTop = 0;
    
    // Add event listeners
    addImageEventListeners();
}

function closeImageModal() {
    const modal = document.getElementById('image-modal');
    modal.style.display = 'none';
    document.body.style.overflow = 'auto';
    removeImageEventListeners();
}

function addImageEventListeners() {
    const container = document.getElementById('image-container');
    const modalImage = document.getElementById('modal-image');
    
    // Scroll wheel zoom
    container.addEventListener('wheel', handleZoom, { passive: false });
    
    // Drag to pan
    container.addEventListener('mousedown', startDrag);
    container.addEventListener('mousemove', drag);
    container.addEventListener('mouseup', endDrag);
    container.addEventListener('mouseleave', endDrag);
    
    // Touch support for mobile
    container.addEventListener('touchstart', handleTouchStart, { passive: false });
    container.addEventListener('touchmove', handleTouchMove, { passive: false });
    container.addEventListener('touchend', handleTouchEnd);
    
    // Double click to zoom
    modalImage.addEventListener('dblclick', handleDoubleClick);
}

function removeImageEventListeners() {
    const container = document.getElementById('image-container');
    const modalImage = document.getElementById('modal-image');
    
    container.removeEventListener('wheel', handleZoom);
    container.removeEventListener('mousedown', startDrag);
    container.removeEventListener('mousemove', drag);
    container.removeEventListener('mouseup', endDrag);
    container.removeEventListener('mouseleave', endDrag);
    container.removeEventListener('touchstart', handleTouchStart);
    container.removeEventListener('touchmove', handleTouchMove);
    container.removeEventListener('touchend', handleTouchEnd);
    modalImage.removeEventListener('dblclick', handleDoubleClick);
}

function handleZoom(e) {
    e.preventDefault();
    
    const container = document.getElementById('image-container');
    const rect = container.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    
    const oldZoom = currentZoom;
    const zoomSpeed = 0.1;
    
    if (e.deltaY < 0) {
        // Zoom in
        currentZoom = Math.min(currentZoom + zoomSpeed, 5);
    } else {
        // Zoom out
        currentZoom = Math.max(currentZoom - zoomSpeed, 0.1);
    }
    
    // Calculate new scroll position to zoom towards cursor
    const zoomRatio = currentZoom / oldZoom;
    const newScrollLeft = (container.scrollLeft + x) * zoomRatio - x;
    const newScrollTop = (container.scrollTop + y) * zoomRatio - y;
    
    updateZoom();
    
    container.scrollLeft = newScrollLeft;
    container.scrollTop = newScrollTop;
}

function startDrag(e) {
    const container = document.getElementById('image-container');
    isDragging = true;
    container.style.cursor = 'grabbing';
    startX = e.pageX - container.offsetLeft;
    startY = e.pageY - container.offsetTop;
    scrollLeft = container.scrollLeft;
    scrollTop = container.scrollTop;
}

function drag(e) {
    if (!isDragging) return;
    e.preventDefault();
    
    const container = document.getElementById('image-container');
    const x = e.pageX - container.offsetLeft;
    const y = e.pageY - container.offsetTop;
    const walkX = (x - startX) * 1;
    const walkY = (y - startY) * 1;
    
    container.scrollLeft = scrollLeft - walkX;
    container.scrollTop = scrollTop - walkY;
}

function endDrag() {
    isDragging = false;
    const container = document.getElementById('image-container');
    container.style.cursor = currentZoom > 1 ? 'grab' : 'default';
}

function handleDoubleClick(e) {
    if (currentZoom === 1) {
        currentZoom = 2;
    } else {
        currentZoom = 1;
    }
    updateZoom();
    centerImage();
}

function zoomIn() {
    currentZoom = Math.min(currentZoom + 0.25, 5);
    updateZoom();
}

function zoomOut() {
    currentZoom = Math.max(currentZoom - 0.25, 0.1);
    updateZoom();
}

function resetZoom() {
    currentZoom = 1;
    updateZoom();
    centerImage();
}

function updateZoom() {
    const modalImage = document.getElementById('modal-image');
    const zoomLevel = document.getElementById('zoom-level');
    const zoomIndicator = document.getElementById('zoom-indicator');
    const container = document.getElementById('image-container');
    
    modalImage.style.transform = `scale(${currentZoom})`;
    const percentage = Math.round(currentZoom * 100);
    zoomLevel.textContent = `${percentage}%`;
    zoomIndicator.textContent = `${percentage}%`;
    
    // Update cursor
    container.style.cursor = currentZoom > 1 ? 'grab' : 'default';
    
    // Show/hide zoom indicator
    if (currentZoom !== 1) {
        zoomIndicator.style.opacity = '1';
        setTimeout(() => {
            zoomIndicator.style.opacity = '0';
        }, 1000);
    }
}

function centerImage() {
    const container = document.getElementById('image-container');
    const modalImage = document.getElementById('modal-image');
    
    setTimeout(() => {
        const containerRect = container.getBoundingClientRect();
        const imageRect = modalImage.getBoundingClientRect();
        
        const scrollLeft = (imageRect.width * currentZoom - containerRect.width) / 2;
        const scrollTop = (imageRect.height * currentZoom - containerRect.height) / 2;
        
        container.scrollLeft = Math.max(0, scrollLeft);
        container.scrollTop = Math.max(0, scrollTop);
    }, 10);
}

function downloadImage() {
    const modalImage = document.getElementById('modal-image');
    const link = document.createElement('a');
    link.href = modalImage.src;
    link.download = 'image.jpg';
    link.click();
}

// Touch support for mobile
let touchStartX, touchStartY, touchDistance = 0;

function handleTouchStart(e) {
    if (e.touches.length === 1) {
        // Single touch - start dragging
        const touch = e.touches[0];
        startDrag({
            pageX: touch.pageX,
            pageY: touch.pageY
        });
    } else if (e.touches.length === 2) {
        // Two touches - prepare for pinch zoom
        const touch1 = e.touches[0];
        const touch2 = e.touches[1];
        touchDistance = Math.sqrt(
            Math.pow(touch2.pageX - touch1.pageX, 2) +
            Math.pow(touch2.pageY - touch1.pageY, 2)
        );
    }
}

function handleTouchMove(e) {
    e.preventDefault();
    
    if (e.touches.length === 1 && isDragging) {
        // Single touch - dragging
        const touch = e.touches[0];
        drag({
            pageX: touch.pageX,
            pageY: touch.pageY,
            preventDefault: () => {}
        });
    } else if (e.touches.length === 2) {
        // Two touches - pinch zoom
        const touch1 = e.touches[0];
        const touch2 = e.touches[1];
        const newDistance = Math.sqrt(
            Math.pow(touch2.pageX - touch1.pageX, 2) +
            Math.pow(touch2.pageY - touch1.pageY, 2)
        );
        
        if (touchDistance > 0) {
            const scale = newDistance / touchDistance;
            currentZoom = Math.max(0.1, Math.min(5, currentZoom * scale));
            updateZoom();
        }
        
        touchDistance = newDistance;
    }
}

function handleTouchEnd(e) {
    if (e.touches.length === 0) {
        endDrag();
        touchDistance = 0;
    }
}

// Close image modal with escape key
document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape') {
        closeImageModal();
    }
});

// Helper function for time ago
function timeAgo(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffInSeconds = Math.floor((now - date) / 1000);
    
    if (diffInSeconds < 60) return 'just now';
    if (diffInSeconds < 3600) return Math.floor(diffInSeconds / 60) + 'm ago';
    if (diffInSeconds < 86400) return Math.floor(diffInSeconds / 3600) + 'h ago';
    if (diffInSeconds < 2592000) return Math.floor(diffInSeconds / 86400) + 'd ago';
    return date.toLocaleDateString();
}

// Truncate text
function truncate(text, length) {
    return text.length > length ? text.substring(0, length) + '...' : text;
}

// Update sort
function updateSort(sortBy) {
    const url = new URL(window.location);
    url.searchParams.set('sort', sortBy);
    url.searchParams.delete('page');
    window.location.href = url.toString();
}

// Update limit
function updateLimit(limit) {
    const url = new URL(window.location);
    url.searchParams.set('limit', limit);
    url.searchParams.delete('page');
    window.location.href = url.toString();
}

// Vote on thread
async function voteThread(threadId, voteType) {
    try {
        const response = await fetch(`/api/threads/${threadId}/vote`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ vote_type: voteType })
        });
        
        if (response.ok) {
            const result = await response.json();
            
            // Update the vote display
            const threadCard = document.querySelector(`[data-thread-id="${threadId}"]`);
            if (threadCard) {
                const scoreElement = threadCard.querySelector('.vote-score');
                const upvoteBtn = threadCard.querySelector('.vote-arrow:first-child');
                const downvoteBtn = threadCard.querySelector('.vote-arrow:last-child');
                
                // Update score
                scoreElement.textContent = result.score;
                scoreElement.dataset.score = result.score;
                
                // Update button states
                upvoteBtn.classList.toggle('upvoted', result.user_vote === 1);
                downvoteBtn.classList.toggle('downvoted', result.user_vote === -1);
                
                // Add visual feedback
                const feedbackClass = voteType === 1 ? 'vote-feedback-up' : 'vote-feedback-down';
                scoreElement.classList.add(feedbackClass);
                setTimeout(() => {
                    scoreElement.classList.remove(feedbackClass);
                }, 300);
            }
            
            if (window.Glorp) {
                const action = result.user_vote === voteType ? 
                    (voteType === 1 ? 'upvoted' : 'downvoted') : 'removed vote';
                Glorp.showNotification(`Post ${action}!`, 'success');
            }
        } else {
            const result = await response.text();
            if (window.Glorp) {
                Glorp.showNotification(result, 'error');
            }
        }
    } catch (error) {
        console.error('Error voting:', error);
        if (window.Glorp) {
            Glorp.showNotification('Failed to vote. Please try again.', 'error');
        }
    }
}

// Enhanced Share Functions
function shareThread(threadId) {
    currentShareThreadId = threadId;
    const shareUrl = `${window.location.origin}/threads/${threadId}`;
    document.getElementById('share-link-input').value = shareUrl;
    document.getElementById('share-modal').style.display = 'flex';
    generateQRCode(shareUrl);
}

function copyShareLink() {
    const input = document.getElementById('share-link-input');
    input.select();
    input.setSelectionRange(0, 99999);
    
    navigator.clipboard.writeText(input.value).then(() => {
        showCopyFeedback();
        if (window.Glorp) {
            Glorp.showNotification('Link copied to clipboard!', 'success');
        }
    }).catch(() => {
        // Fallback for older browsers
        document.execCommand('copy');
        showCopyFeedback();
        if (window.Glorp) {
            Glorp.showNotification('Link copied to clipboard!', 'success');
        }
    });
}

function showCopyFeedback() {
    const copyBtn = document.querySelector('.copy-btn');
    const originalHTML = copyBtn.innerHTML;
    copyBtn.innerHTML = '<i class="fas fa-check"></i>';
    copyBtn.style.background = '#28a745';
    
    setTimeout(() => {
        copyBtn.innerHTML = originalHTML;
        copyBtn.style.background = '';
    }, 1500);
}

// Enhanced Social Sharing Functions
function shareToReddit() {
    const url = document.getElementById('share-link-input').value;
    const title = `Check out this post on Glorp!`;
    window.open(`https://reddit.com/submit?url=${encodeURIComponent(url)}&title=${encodeURIComponent(title)}`, '_blank');
}

function shareToTwitter() {
    const url = document.getElementById('share-link-input').value;
    const text = 'Check out this interesting post on Glorp! 🚀';
    window.open(`https://twitter.com/intent/tweet?text=${encodeURIComponent(text)}&url=${encodeURIComponent(url)}`, '_blank');
}

function shareToFacebook() {
    const url = document.getElementById('share-link-input').value;
    window.open(`https://www.facebook.com/sharer/sharer.php?u=${encodeURIComponent(url)}`, '_blank');
}

function shareToDiscord() {
    const url = document.getElementById('share-link-input').value;
    // Discord doesn't have direct sharing URL, so we copy and show instruction
    navigator.clipboard.writeText(url);
    if (window.Glorp) {
        Glorp.showNotification('Link copied! Paste it in Discord 💬', 'info');
    }
}

function shareToTelegram() {
    const url = document.getElementById('share-link-input').value;
    const text = 'Check out this post on Glorp!';
    window.open(`https://t.me/share/url?url=${encodeURIComponent(url)}&text=${encodeURIComponent(text)}`, '_blank');
}

function shareToWhatsApp() {
    const url = document.getElementById('share-link-input').value;
    const text = `Check out this interesting post: ${url}`;
    window.open(`https://wa.me/?text=${encodeURIComponent(text)}`, '_blank');
}

// Cross-post functionality
async function createCrossPost() {
    const community = document.getElementById('crosspost-community').value;
    if (!community) {
        if (window.Glorp) {
            Glorp.showNotification('Please select a community first', 'warning');
        }
        return;
    }
    
    if (!currentShareThreadId) {
        if (window.Glorp) {
            Glorp.showNotification('No post selected for cross-posting', 'error');
        }
        return;
    }
    
    // For now, redirect to create page with pre-filled data
    // In a real implementation, you'd make an API call to create the cross-post
    const originalUrl = `${window.location.origin}/threads/${currentShareThreadId}`;
    const crossPostUrl = `/threads/create?crosspost=${currentShareThreadId}&community=${community}`;
    
    if (window.Glorp) {
        Glorp.showNotification(`Cross-posting to z/${community}...`, 'info');
    }
    
    // Close modal and redirect
    setTimeout(() => {
        closeModal('share-modal');
        window.location.href = crossPostUrl;
    }, 1000);
}

// QR Code generation (simple text-based for now)
function generateQRCode(url) {
    const qrContainer = document.getElementById('qr-code');
    // In a real implementation, you'd use a QR code library like qrcode.js
    // For now, we'll show a placeholder that links to a QR generator
    qrContainer.innerHTML = `
        <div class="qr-placeholder-content">
            <i class="fas fa-qrcode"></i>
            <span>QR Code</span>
            <small>Click to generate</small>
        </div>
    `;
    
    qrContainer.onclick = () => {
        window.open(`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(url)}`, '_blank');
    };
}

function downloadQR() {
    const url = document.getElementById('share-link-input').value;
    const qrUrl = `https://api.qrserver.com/v1/create-qr-code/?size=400x400&format=png&data=${encodeURIComponent(url)}`;
    
    const link = document.createElement('a');
    link.href = qrUrl;
    link.download = 'glorp-post-qr.png';
    link.click();
}

function closeModal(modalId) {
    document.getElementById(modalId).style.display = 'none';
}

// Template helper functions
function sub(a, b) { return a - b; }
function add(a, b) { return a + b; }

// Close modal when clicking outside, enhanced for new modal
document.addEventListener('click', function(e) {
    const shareModal = document.getElementById('share-modal');
    const imageModal = document.getElementById('image-modal');
    
    if (e.target === shareModal || e.target.classList.contains('share-modal-overlay')) {
        closeModal('share-modal');
    }
    
    if (e.target === imageModal || e.target.classList.contains('image-modal-overlay')) {
        closeImageModal();
    }
    
    // Prevent image clicks from bubbling up
    if (e.target.classList.contains('preview-image')) {
        e.stopPropagation();
    }
});

// Initialize time ago updates
setInterval(() => {
    document.querySelectorAll('[title*="2024"]').forEach(el => {
        const date = new Date(el.getAttribute('title'));
        el.textContent = timeAgo(date);
    });
}, 60000); // Update every minute

// Stop propagation on preview images so clicking them doesn't navigate to thread
document.querySelectorAll('.preview-image').forEach(img => {
    img.addEventListener('click', function(e) {
        e.stopPropagation();
        e.preventDefault();
    });
});