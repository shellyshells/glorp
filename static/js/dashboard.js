async function updateThreadStatus(threadId, status) {
    if (!status) return;
    
    const confirmMessage = `Are you sure you want to ${status} this thread?`;
    if (!confirm(confirmMessage)) {
        return;
    }
    
    try {
        const response = await fetch(`/api/admin/threads/${threadId}/status`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ status })
        });
        
        if (response.ok) {
            location.reload(); // Refresh to show updated status
        } else {
            const result = await response.text();
            alert('Failed to update thread status: ' + result);
        }
    } catch (error) {
        alert('Failed to update thread status. Please try again.');
    }
}

async function deleteThread(threadId) {
    if (!confirm('Are you sure you want to delete this thread? This action cannot be undone.')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/threads/${threadId}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            location.reload(); // Refresh to remove deleted thread
        } else {
            const result = await response.text();
            alert('Failed to delete thread: ' + result);
        }
    } catch (error) {
        alert('Failed to delete thread. Please try again.');
    }
}

async function banUser(userId, username) {
    const action = confirm(`Are you sure you want to ban user "${username}"?`) ? 'ban' : null;
    if (!action) return;
    
    try {
        const response = await fetch(`/api/admin/ban/${userId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ action })
        });
        
        if (response.ok) {
            alert(`User "${username}" has been banned successfully.`);
            location.reload();
        } else {
            const result = await response.text();
            alert('Failed to ban user: ' + result);
        }
    } catch (error) {
        alert('Failed to ban user. Please try again.');
    }
}

async function deleteUser(userId) {
    if (!confirm('Are you sure you want to delete this user? This action cannot be undone.')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/admin/users/${userId}`, {
            method: 'DELETE'
        });
        
        if (response.ok) {
            location.reload();
        } else {
            const result = await response.text();
            alert('Failed to delete user: ' + result);
        }
    } catch (error) {
        alert('Failed to delete user. Please try again.');
    }
}

async function deleteMessage(messageId) {
    if (!confirm('Are you sure you want to delete this message? This action cannot be undone.')) {
        return;
    }

    try {
        const response = await fetch(`/api/admin/messages/${messageId}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            location.reload();
        } else {
            const result = await response.text();
            alert('Failed to delete message: ' + result);
        }
    } catch (error) {
        alert('Failed to delete message. Please try again.');
    }
}

async function editMessage(messageId) {
    const newContent = prompt('Enter new content for the message:');
    if (newContent === null || newContent.trim() === '') {
        return;
    }

    try {
        const response = await fetch(`/api/admin/messages/${messageId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ content: newContent })
        });

        if (response.ok) {
            location.reload();
        } else {
            const result = await response.text();
            alert('Failed to edit message: ' + result);
        }
    } catch (error) {
        alert('Failed to edit message. Please try again.');
    }
}

async function deleteCommunity(communityId) {
    if (!confirm('Are you sure you want to delete this community? This action cannot be undone.')) {
        return;
    }

    try {
        const response = await fetch(`/api/admin/communities/${communityId}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            location.reload();
        } else {
            const result = await response.text();
            alert('Failed to delete community: ' + result);
        }
    } catch (error) {
        alert('Failed to delete community. Please try again.');
    }
}

async function editCommunity(communityId, currentName, currentDescription) {
    const newName = prompt('Enter new name for the community:', currentName);
    if (newName === null || newName.trim() === '') {
        return;
    }

    const newDescription = prompt('Enter new description for the community:', currentDescription);
    if (newDescription === null) {
        return; // User cancelled description edit
    }

    try {
        const response = await fetch(`/api/admin/communities/${communityId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name: newName, description: newDescription })
        });

        if (response.ok) {
            location.reload();
        } else {
            const result = await response.text();
            alert('Failed to edit community: ' + result);
        }
    } catch (error) {
        alert('Failed to edit community. Please try again.');
    }
}