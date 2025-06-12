document.getElementById('thread-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const formData = new FormData(this);
    const selectedTags = Array.from(formData.getAll('tags')).map(id => parseInt(id));
    
    const data = {
        title: formData.get('title'),
        description: formData.get('description'),
        tags: selectedTags
    };
    
    const errorDiv = document.getElementById('error-message');
    
    try {
        const response = await fetch('/api/threads/{{.Thread.ID}}', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });
        
        if (response.ok) {
            window.location.href = '/threads/{{.Thread.ID}}';
        } else {
            const result = await response.text();
            errorDiv.textContent = result;
            errorDiv.style.display = 'block';
        }
    } catch (error) {
        errorDiv.textContent = 'Failed to update thread. Please try again.';
        errorDiv.style.display = 'block';
    }
});