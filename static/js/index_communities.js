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
    try {
        const response = await fetch(`/api/communities/${communityId}/join`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({})
        });
        
        if (response.ok) {
            location.reload();
        } else {
            const result = await response.text();
            alert('Failed to join community: ' + result);
        }
    } catch (error) {
        alert('Failed to join community. Please try again.');
    }
}