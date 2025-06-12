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

// Join community functionality
async function joinCommunity(communityId) {
    let message = '';
    
    // If approval required, ask for message
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

// Leave community functionality
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

// Process join request
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

// Voting and other functions from main thread list
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

function shareThread(threadId) {
    if (window.Glorp) {
        Glorp.showNotification('Share feature coming soon!', 'info');
    }
}

function showThreadOptions(threadId) {
    if (window.Glorp) {
        Glorp.showNotification('Thread options coming soon!', 'info');
    }
}