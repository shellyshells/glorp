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