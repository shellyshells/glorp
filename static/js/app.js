// Enhanced JavaScript functionality for glorp

document.addEventListener('DOMContentLoaded', function() {
    // Initialize all functionality
    initializeNavigation();
    initializeSearch();
    initializeLogout();
    initializeNotifications();
    initializeKeyboardShortcuts();
    
    // Initialize any vote buttons
    initializeVoting();
});

// Navigation functionality
function initializeNavigation() {
    // Mobile menu toggle (if needed in future)
    const navToggle = document.querySelector('.nav-toggle');
    const navMenu = document.querySelector('.nav-menu');
    
    if (navToggle && navMenu) {
        navToggle.addEventListener('click', function() {
            navMenu.classList.toggle('active');
        });
    }

    // Active page highlighting
    const currentPath = window.location.pathname;
    const navLinks = document.querySelectorAll('.nav-link');
    
    navLinks.forEach(link => {
        if (link.getAttribute('href') === currentPath) {
            link.classList.add('active');
        }
    });
}

// Enhanced search functionality
function initializeSearch() {
    const searchForm = document.getElementById('search-form');
    const searchInput = document.getElementById('search-input');
    
    if (searchForm) {
        searchForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const query = searchInput.value.trim();
            
            if (query) {
                // Redirect to home page with search query
                window.location.href = `/?search=${encodeURIComponent(query)}`;
            }
        });
    }
    
    // Auto-submit search on Enter key
    if (searchInput) {
        searchInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                e.preventDefault();
                searchForm.dispatchEvent(new Event('submit'));
            }
        });

        // Search suggestions (placeholder for future enhancement)
        let searchTimeout;
        searchInput.addEventListener('input', function() {
            clearTimeout(searchTimeout);
            const query = this.value.trim();
            
            if (query.length >= 2) {
                searchTimeout = setTimeout(() => {
                    // You could implement search suggestions here
                    // fetchSearchSuggestions(query);
                }, 300);
            }
        });
    }
}

// Enhanced logout functionality
function initializeLogout() {
    const logoutBtn = document.getElementById('logout-btn');
    
    if (logoutBtn) {
        logoutBtn.addEventListener('click', async function(e) {
            e.preventDefault();
            
            // Show loading state
            const originalText = this.innerHTML;
            this.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Logging out...';
            this.disabled = true;
            
            try {
                const response = await fetch('/api/logout', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });
                
                if (response.ok) {
                    showNotification('Logged out successfully!', 'success');
                    setTimeout(() => {
                        window.location.href = '/';
                    }, 1000);
                } else {
                    console.error('Logout failed');
                    showNotification('Logout failed, redirecting anyway...', 'warning');
                    setTimeout(() => {
                        window.location.href = '/';
                    }, 1500);
                }
            } catch (error) {
                console.error('Logout error:', error);
                showNotification('Network error, redirecting...', 'warning');
                setTimeout(() => {
                    window.location.href = '/';
                }, 1500);
            } finally {
                // Reset button state
                this.innerHTML = originalText;
                this.disabled = false;
            }
        });
    }
}

// Initialize voting functionality
function initializeVoting() {
    // Add event listeners to vote buttons
    document.querySelectorAll('.vote-arrow, .comment-vote-btn').forEach(button => {
        if (!button.onclick && !button.href) {
            button.addEventListener('click', function(e) {
                e.preventDefault();
                showNotification('Please log in to vote', 'info');
            });
        }
    });
}

// Keyboard shortcuts
function initializeKeyboardShortcuts() {
    document.addEventListener('keydown', function(e) {
        // Only trigger shortcuts when not typing in input fields
        if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
            return;
        }

        switch(e.key) {
            case 'c':
                // 'c' to create new post
                if (document.querySelector('.btn-create')) {
                    window.location.href = '/threads/create';
                }
                break;
            case 'h':
                // 'h' to go home
                window.location.href = '/';
                break;
            case '/':
                // '/' to focus search
                e.preventDefault();
                const searchInput = document.getElementById('search-input');
                if (searchInput) {
                    searchInput.focus();
                }
                break;
            case 'Escape':
                // Escape to close modals
                document.querySelectorAll('.modal').forEach(modal => {
                    modal.style.display = 'none';
                });
                break;
        }
    });
}

// Enhanced notification system
function initializeNotifications() {
    // Check for URL parameters that might indicate messages
    const urlParams = new URLSearchParams(window.location.search);
    const error = urlParams.get('error');
    const success = urlParams.get('success');
    
    if (error === 'banned') {
        showNotification('Your account has been banned', 'error');
    } else if (success === 'registered') {
        showNotification('Registration successful! Welcome to Glorp!', 'success');
    }
}

// Enhanced notification function
function showNotification(message, type = 'info', duration = 5000) {
    // Remove existing notifications
    document.querySelectorAll('.notification').forEach(n => n.remove());
    
    // Create notification element
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.innerHTML = `
        <div class="notification-content">
            <i class="notification-icon ${getNotificationIcon(type)}"></i>
            <span class="notification-message">${message}</span>
            <button class="notification-close" onclick="removeNotification(this.parentElement)">&times;</button>
        </div>
    `;
    
    // Add styles
    notification.style.cssText = `
        position: fixed;
        top: 80px;
        right: 20px;
        padding: 16px 20px;
        border-radius: 8px;
        color: white;
        font-weight: 500;
        z-index: 10000;
        display: flex;
        align-items: center;
        gap: 12px;
        max-width: 400px;
        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        animation: slideIn 0.3s ease-out;
        backdrop-filter: blur(10px);
    `;
    
    // Set background color based on type
    switch (type) {
        case 'success':
            notification.style.background = 'linear-gradient(135deg, #28a745, #20c997)';
            break;
        case 'error':
            notification.style.background = 'linear-gradient(135deg, #dc3545, #fd7e14)';
            break;
        case 'warning':
            notification.style.background = 'linear-gradient(135deg, #ffc107, #fd7e14)';
            notification.style.color = '#000';
            break;
        case 'info':
        default:
            notification.style.background = 'linear-gradient(135deg, #0079d3, #00a8ff)';
    }
    
    // Style the content
    const content = notification.querySelector('.notification-content');
    content.style.cssText = `
        display: flex;
        align-items: center;
        gap: 12px;
        width: 100%;
    `;
    
    // Style the close button
    const closeBtn = notification.querySelector('.notification-close');
    closeBtn.style.cssText = `
        background: none;
        border: none;
        color: inherit;
        font-size: 20px;
        cursor: pointer;
        padding: 0;
        margin-left: auto;
        border-radius: 4px;
        width: 24px;
        height: 24px;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: background-color 0.2s;
    `;
    
    closeBtn.addEventListener('mouseenter', function() {
        this.style.backgroundColor = 'rgba(255,255,255,0.2)';
    });
    
    closeBtn.addEventListener('mouseleave', function() {
        this.style.backgroundColor = 'transparent';
    });
    
    // Add to page
    document.body.appendChild(notification);
    
    // Auto-remove after duration
    if (duration > 0) {
        setTimeout(function() {
            removeNotification(notification);
        }, duration);
    }
}

function getNotificationIcon(type) {
    switch (type) {
        case 'success':
            return 'fas fa-check-circle';
        case 'error':
            return 'fas fa-exclamation-circle';
        case 'warning':
            return 'fas fa-exclamation-triangle';
        case 'info':
        default:
            return 'fas fa-info-circle';
    }
}

function removeNotification(notification) {
    if (notification && notification.parentNode) {
        notification.style.animation = 'slideOut 0.3s ease-in';
        setTimeout(function() {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 300);
    }
}

// Utility function to format dates
function formatDate(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffTime < 60000) {
        return 'just now';
    } else if (diffTime < 3600000) {
        const minutes = Math.floor(diffTime / 60000);
        return `${minutes}m ago`;
    } else if (diffTime < 86400000) {
        const hours = Math.floor(diffTime / 3600000);
        return `${hours}h ago`;
    } else if (diffDays === 1) {
        return 'yesterday';
    } else if (diffDays < 7) {
        return `${diffDays} days ago`;
    } else {
        return date.toLocaleDateString();
    }
}

// Utility function to truncate text
function truncateText(text, maxLength) {
    if (text.length <= maxLength) {
        return text;
    }
    return text.substring(0, maxLength) + '...';
}

// Utility function to escape HTML
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Function to handle form validation
function validateForm(formElement) {
    const inputs = formElement.querySelectorAll('input[required], textarea[required]');
    let isValid = true;
    
    inputs.forEach(input => {
        const value = input.value.trim();
        
        if (!value) {
            addFieldError(input, 'This field is required');
            isValid = false;
        } else {
            removeFieldError(input);
            
            // Additional validation based on type
            if (input.type === 'email' && !isValidEmail(value)) {
                addFieldError(input, 'Please enter a valid email address');
                isValid = false;
            } else if (input.type === 'password' && !isValidPassword(value)) {
                addFieldError(input, 'Password must be at least 12 characters with 1 uppercase and 1 special character');
                isValid = false;
            }
        }
    });
    
    return isValid;
}

function addFieldError(input, message) {
    input.classList.add('error');
    
    // Remove existing error message
    const existingError = input.parentNode.querySelector('.field-error');
    if (existingError) {
        existingError.remove();
    }
    
    // Add new error message
    const errorDiv = document.createElement('div');
    errorDiv.className = 'field-error';
    errorDiv.textContent = message;
    input.parentNode.appendChild(errorDiv);
}

function removeFieldError(input) {
    input.classList.remove('error');
    const errorDiv = input.parentNode.querySelector('.field-error');
    if (errorDiv) {
        errorDiv.remove();
    }
}

function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

function isValidPassword(password) {
    return password.length >= 12 && 
           /[A-Z]/.test(password) && 
           /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password);
}

// Enhanced copy to clipboard function
async function copyToClipboard(text) {
    try {
        if (navigator.clipboard && window.isSecureContext) {
            await navigator.clipboard.writeText(text);
        } else {
            // Fallback for older browsers
            const textArea = document.createElement('textarea');
            textArea.value = text;
            textArea.style.position = 'fixed';
            textArea.style.left = '-999999px';
            textArea.style.top = '-999999px';
            document.body.appendChild(textArea);
            textArea.focus();
            textArea.select();
            document.execCommand('copy');
            textArea.remove();
        }
        showNotification('Copied to clipboard!', 'success');
        return true;
    } catch (err) {
        console.error('Failed to copy:', err);
        showNotification('Failed to copy to clipboard', 'error');
        return false;
    }
}

// Lazy loading for images (if needed)
function initializeLazyLoading() {
    if ('IntersectionObserver' in window) {
        const imageObserver = new IntersectionObserver((entries, observer) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    const img = entry.target;
                    img.src = img.dataset.src;
                    img.classList.remove('lazy');
                    imageObserver.unobserve(img);
                }
            });
        });

        document.querySelectorAll('img[data-src]').forEach(img => {
            imageObserver.observe(img);
        });
    }
}

// Add error styling
const errorStyle = document.createElement('style');
errorStyle.textContent = `
    .error {
        border-color: #dc3545 !important;
        box-shadow: 0 0 0 2px rgba(220, 53, 69, 0.2) !important;
    }
    
    .field-error {
        color: #dc3545;
        font-size: 12px;
        margin-top: 4px;
        display: block;
    }
    
    .notification-content {
        display: flex;
        align-items: center;
        gap: 12px;
        width: 100%;
    }
    
    .notification-icon {
        font-size: 18px;
        flex-shrink: 0;
    }
    
    .notification-message {
        flex: 1;
        word-wrap: break-word;
    }
`;
document.head.appendChild(errorStyle);

// Export functions for global use
window.Glorp = {
    showNotification,
    formatDate,
    truncateText,
    escapeHtml,
    validateForm,
    copyToClipboard,
    removeNotification
};

// Add CSS animations if not already present
if (!document.querySelector('#glorp-animations')) {
    const animationStyle = document.createElement('style');
    animationStyle.id = 'glorp-animations';
    animationStyle.textContent = `
        @keyframes slideIn {
            from {
                transform: translateX(100%);
                opacity: 0;
            }
            to {
                transform: translateX(0);
                opacity: 1;
            }
        }
        
        @keyframes slideOut {
            from {
                transform: translateX(0);
                opacity: 1;
            }
            to {
                transform: translateX(100%);
                opacity: 0;
            }
        }
        
        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }
        
        @keyframes fadeOut {
            from { opacity: 1; }
            to { opacity: 0; }
        }
        
        .loading-spinner {
            display: inline-block;
            width: 16px;
            height: 16px;
            border: 2px solid #f3f3f3;
            border-top: 2px solid currentColor;
            border-radius: 50%;
            animation: spin 1s linear infinite;
        }
        
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        
        .fade-in {
            animation: fadeIn 0.3s ease-in;
        }
        
        .fade-out {
            animation: fadeOut 0.3s ease-out;
        }
    `;
    document.head.appendChild(animationStyle);
}

// Function to delete a thread
async function deleteThread(threadId) {
    if (!confirm("Are you sure you want to delete this post? This action cannot be undone.")) {
        return;
    }

    try {
        const response = await fetch(`/api/threads/${threadId}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        const data = await response.json();

        if (response.ok) {
            showNotification(data.message || 'Post deleted successfully!', 'success');
            // Redirect to home page or user's profile after deletion
            setTimeout(() => {
                window.location.href = '/'; 
            }, 1000);
        } else {
            showNotification(data.message || 'Failed to delete post.', 'error');
        }
    } catch (error) {
        console.error('Error deleting thread:', error);
        showNotification('An error occurred while deleting the post.', 'error');
    }
}

// Function to delete a comment
async function deleteComment(commentId) {
    if (!confirm("Are you sure you want to delete this comment? This action cannot be undone.")) {
        return;
    }

    try {
        const response = await fetch(`/api/messages/${commentId}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        const data = await response.json();

        if (response.ok) {
            showNotification(data.message || 'Comment deleted successfully!', 'success');
            // Remove the comment from the DOM
            const commentElement = document.querySelector(`.comment-item[data-comment-id="${commentId}"]`);
            if (commentElement) {
                commentElement.remove();
            }
        } else {
            showNotification(data.message || 'Failed to delete comment.', 'error');
        }
    } catch (error) {
        console.error('Error deleting comment:', error);
        showNotification('An error occurred while deleting the comment.', 'error');
    }
}

// Function to vote on a thread
async function voteThread(threadId, voteType) {
    try {
        const response = await fetch(`/api/threads/${threadId}/vote`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ vote_type: voteType }),
        });

        const data = await response.json();

        if (response.ok) {
            document.getElementById('thread-score').textContent = data.score;
            const upvoteBtn = document.querySelector(`.thread-vote-section button[onclick*="voteThread(\'${threadId}\', 1)"]`);
            const downvoteBtn = document.querySelector(`.thread-vote-section button[onclick*="voteThread(\'${threadId}\', -1)"]`);

            if (upvoteBtn && downvoteBtn) {
                upvoteBtn.classList.remove('upvoted');
                downvoteBtn.classList.remove('downvoted');

                if (data.user_vote === 1) {
                    upvoteBtn.classList.add('upvoted');
                } else if (data.user_vote === -1) {
                    downvoteBtn.classList.add('downvoted');
                }
            }
            showNotification(data.message, 'success');
        } else {
            showNotification(data.message || 'Failed to record vote.', 'error');
        }
    } catch (error) {
        console.error('Error voting on thread:', error);
        showNotification('An error occurred while voting.', 'error');
    }
}

// Function to vote on a comment
async function voteComment(commentId, voteType) {
    try {
        const response = await fetch(`/api/messages/${commentId}/vote`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ vote_type: voteType }),
        });

        const data = await response.json();

        if (response.ok) {
            const commentItem = document.querySelector(`.comment-item[data-comment-id="${commentId}"]`);
            if (commentItem) {
                commentItem.querySelector('.comment-score').textContent = data.score;
                commentItem.querySelector('.comment-score-inline').textContent = `${data.score} points`;
                
                const upvoteBtn = commentItem.querySelector(`.comment-vote-btn[onclick*="voteComment(\'${commentId}\', 1)"]`);
                const downvoteBtn = commentItem.querySelector(`.comment-vote-btn[onclick*="voteComment(\'${commentId}\', -1)"]`);

                if (upvoteBtn && downvoteBtn) {
                    upvoteBtn.classList.remove('upvoted');
                    downvoteBtn.classList.remove('downvoted');

                    if (data.user_vote === 1) {
                        upvoteBtn.classList.add('upvoted');
                    } else if (data.user_vote === -1) {
                        downvoteBtn.classList.add('downvoted');
                    }
                }
            }
            showNotification(data.message, 'success');
        } else {
            showNotification(data.message || 'Failed to record vote.', 'error');
        }
    } catch (error) {
        console.error('Error voting on comment:', error);
        showNotification('An error occurred while voting.', 'error');
    }
}