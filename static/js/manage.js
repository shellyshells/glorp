// Update community settings
document.getElementById('communitySettingsForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const formData = {
        display_name: document.getElementById('displayName').value,
        description: document.getElementById('description').value,
        visibility: document.getElementById('visibility').value,
        join_approval: document.getElementById('joinApproval').value
    };
    
    try {
        const response = await fetch(`/api/communities/${communityData.id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });
        
        if (response.ok) {
            if (window.Glorp) {
                Glorp.showNotification('Community settings updated successfully', 'success');
            } else {
                alert('Community settings updated successfully');
            }
            location.reload();
        } else {
            const result = await response.text();
            if (window.Glorp) {
                Glorp.showNotification(result, 'error');
            } else {
                alert('Failed to update settings: ' + result);
            }
        }
    } catch (error) {
        if (window.Glorp) {
            Glorp.showNotification('Failed to update settings. Please try again.', 'error');
        } else {
            alert('Failed to update settings. Please try again.');
        }
    }
});

// Update moderator role
async function updateModeratorRole(userId, role) {
    try {
        const response = await fetch(`/api/communities/${communityData.id}/moderators/${userId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                action: 'add',
                role: role
            })
        });
        
        if (response.ok) {
            if (window.Glorp) {
                Glorp.showNotification('Moderator role updated successfully', 'success');
            } else {
                alert('Moderator role updated successfully');
            }
            location.reload();
        } else {
            const result = await response.text();
            if (window.Glorp) {
                Glorp.showNotification(result, 'error');
            } else {
                alert('Failed to update role: ' + result);
            }
        }
    } catch (error) {
        if (window.Glorp) {
            Glorp.showNotification('Failed to update role. Please try again.', 'error');
        } else {
            alert('Failed to update role. Please try again.');
        }
    }
}

// Remove moderator
async function removeModerator(userId) {
    if (!confirm('Are you sure you want to remove this moderator?')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/communities/${communityData.id}/moderators/${userId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                action: 'remove'
            })
        });
        
        if (response.ok) {
            if (window.Glorp) {
                Glorp.showNotification('Moderator removed successfully', 'success');
            } else {
                alert('Moderator removed successfully');
            }
            location.reload();
        } else {
            const result = await response.text();
            if (window.Glorp) {
                Glorp.showNotification(result, 'error');
            } else {
                alert('Failed to remove moderator: ' + result);
            }
        }
    } catch (error) {
        if (window.Glorp) {
            Glorp.showNotification('Failed to remove moderator. Please try again.', 'error');
        } else {
            alert('Failed to remove moderator. Please try again.');
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
            if (window.Glorp) {
                Glorp.showNotification(`Join request ${approved ? 'approved' : 'rejected'} successfully`, 'success');
            } else {
                alert(`Join request ${approved ? 'approved' : 'rejected'} successfully`);
            }
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