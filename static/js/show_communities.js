// Include all the image modal functions from index_threads.js
let currentZoom = 1;
let isDragging = false;
let startX, startY, scrollLeft, scrollTop;
let currentShareThreadId = null;

// Advanced Image modal functions with zoom and pan
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
        currentZoom = Math.min(currentZoom + zoomSpeed, 5);
    } else {
        currentZoom = Math.max(currentZoom - zoomSpeed, 0.1);
    }
    
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
    
    container.style.cursor = currentZoom > 1 ? 'grab' : 'default';
    
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

// Touch support
let touchStartX, touchStartY, touchDistance = 0;

function handleTouchStart(e) {
    if (e.touches.length === 1) {
        const touch = e.touches[0];
        startDrag({
            pageX: touch.pageX,
            pageY: touch.pageY
        });
    } else if (e.touches.length === 2) {
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
        const touch = e.touches[0];
        drag({
            pageX: touch.pageX,
            pageY: touch.pageY,
            preventDefault: () => {}
        });
    } else if (e.touches.length === 2) {
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

// Filter tabs functionality
document.querySelectorAll('.filter-tab').forEach(tab => {
    tab.addEventListener('click', function() {
        document.querySelectorAll('.filter-tab').forEach(t => t.classList.remove('active'));
        this.classList.add('active');
        
        const sort = this.dataset.sort;
        const url = new URL(window.location);
        url.searchParams.set('sort', sort);
        window.location.href = url.toString();
    });
});

// Join/Leave community functionality
async function joinCommunity(communityId) {
    let message = '';
    
    if ('{{.Community.JoinApproval}}' === 'approval_required') {
        message = prompt('Please provide a message with your join request (optional):') || '';
    }
    
    try {
        const response = await fetch(`/api/communities/${communityId}/join`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ message: message })
        });
        
        if (response.ok) {
            location.reload();
        } else {
            const result = await response.text();
            if (window.Glorp) {
                Glorp.showNotification(result, 'error');
            } else {
                alert('Failed to join community: ' + result);
            }
        }
    } catch (error) {
        if (window.Glorp) {
            Glorp.showNotification('Failed to join community. Please try again.', 'error');
        } else {
            alert('Failed to join community. Please try again.');
        }
    }
}

async function leaveCommunity(communityId) {
    if (!confirm('Are you sure you want to leave this community?')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/communities/${communityId}/leave`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            }
        });
        
        if (response.ok) {
            location.reload();
        } else {
            const result = await response.text();
            if (window.Glorp) {
                Glorp.showNotification(result, 'error');
            } else {
                alert('Failed to leave community: ' + result);
            }
        }
    } catch (error) {
        if (window.Glorp) {
            Glorp.showNotification('Failed to leave community. Please try again.', 'error');
        } else {
            alert('Failed to leave community. Please try again.');
        }
    }
}

async function processRequest(requestId, approved) {
    try {
        const response = await fetch(`/api/communities/join-requests/${requestId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ approved: approved })
        });
        
        if (response.ok) {
            location.reload();
        } else {
            const result = await response.text();
            if (window.Glorp) {
                Glorp.showNotification(result, 'error');
            } else {
                alert('Failed to process request: ' + result);
            }
        }
    } catch (error) {
        if (window.Glorp) {
            Glorp.showNotification('Failed to process request. Please try again.', 'error');
        } else {
            alert('Failed to process request. Please try again.');
        }
    }
}

// Vote thread functionality
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
            
            const threadCard = document.querySelector(`[data-thread-id="${threadId}"]`);
            if (threadCard) {
                const scoreElement = threadCard.querySelector('.vote-score');
                const upvoteBtn = threadCard.querySelector('.vote-arrow:first-child');
                const downvoteBtn = threadCard.querySelector('.vote-arrow:last-child');
                
                scoreElement.textContent = result.score;
                scoreElement.dataset.score = result.score;
                
                upvoteBtn.classList.toggle('upvoted', result.user_vote === 1);
                downvoteBtn.classList.toggle('downvoted', result.user_vote === -1);
                
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

// Share thread functionality
function shareThread(threadId) {
    currentShareThreadId = threadId;
    const shareUrl = `${window.location.origin}/threads/${threadId}`;
    document.getElementById('share-link-input').value = shareUrl;
    document.getElementById('share-modal').style.display = 'flex';
}

function copyShareLink() {
    const input = document.getElementById('share-link-input');
    input.select();
    input.setSelectionRange(0, 99999);
    
    navigator.clipboard.writeText(input.value).then(() => {
        if (window.Glorp) {
            Glorp.showNotification('Link copied to clipboard!', 'success');
        }
        closeModal('share-modal');
    }).catch(() => {
        document.execCommand('copy');
        if (window.Glorp) {
            Glorp.showNotification('Link copied to clipboard!', 'success');
        }
        closeModal('share-modal');
    });
}

function closeModal(modalId) {
    document.getElementById(modalId).style.display = 'none';
}

// Close image modal with escape key
document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape') {
        closeImageModal();
    }
});

// Close modal when clicking outside
document.addEventListener('click', function(e) {
    const shareModal = document.getElementById('share-modal');
    const imageModal = document.getElementById('image-modal');
    
    if (e.target === shareModal || e.target.classList.contains('share-modal-overlay')) {
        closeModal('share-modal');
    }
    
    if (e.target === imageModal || e.target.classList.contains('image-modal-overlay')) {
        closeImageModal();
    }
});