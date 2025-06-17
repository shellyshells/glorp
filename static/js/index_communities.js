function updateFilter(type, value) {
    const url = new URL(window.location);
    if (value) {
        url.searchParams.set(type, value);
    } else {
        url.searchParams.delete(type);
    }
    url.searchParams.delete('page');
    window.location.href = url.toString();
}

async function joinCommunity(communityId) {
    const community = communityData.communities.find(c => c.id === parseInt(communityId));
    let message = '';
    
    if (community && community.joinApproval === 'approval_required') {
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