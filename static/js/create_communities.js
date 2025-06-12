let nameCheckTimeout;

// Auto-fill display name from community name
document.getElementById('community-name').addEventListener('input', function() {
    const name = this.value;
    const displayName = document.getElementById('display-name');
    
    if (!displayName.value || displayName.value === displayName.dataset.autoFilled) {
        displayName.value = name;
        displayName.dataset.autoFilled = name;
    }
    
    // Check name availability
    clearTimeout(nameCheckTimeout);
    if (name.length >= 3) {
        nameCheckTimeout = setTimeout(() => checkNameAvailability(name), 500);
    } else {
        document.getElementById('name-availability').innerHTML = '';
    }
});

async function checkNameAvailability(name) {
    if (!/^[a-zA-Z0-9_]+$/.test(name)) {
        document.getElementById('name-availability').innerHTML = 
            '<span class="availability-error"><i class="fas fa-times"></i> Invalid characters</span>';
        return;
    }
    
    try {
        const response = await fetch(`/api/communities/check-name?name=${encodeURIComponent(name)}`);
        const result = await response.json();
        
        const availability = document.getElementById('name-availability');
        if (result.available) {
            availability.innerHTML = '<span class="availability-success"><i class="fas fa-check"></i> Available</span>';
        } else {
            availability.innerHTML = '<span class="availability-error"><i class="fas fa-times"></i> Name already taken</span>';
        }
    } catch (error) {
        console.error('Error checking name availability:', error);
    }
}

// Form submission
document.getElementById('create-community-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const formData = new FormData(this);
    const data = {
        name: formData.get('name'),
        display_name: formData.get('display_name'),
        description: formData.get('description'),
        visibility: formData.get('visibility'),
        join_approval: formData.get('join_approval')
    };
    
    const submitBtn = this.querySelector('button[type="submit"]');
    const originalText = submitBtn.innerHTML;
    
    submitBtn.disabled = true;
    submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Creating...';
    
    try {
        const response = await fetch('/api/communities', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });
        
        if (response.ok) {
            const result = await response.json();
            showFormMessage('Community created successfully! Redirecting...', 'success');
            setTimeout(() => {
                window.location.href = `/r/${result.community.name}`;
            }, 1500);
        } else {
            const result = await response.text();
            showFormMessage(result, 'error');
        }
    } catch (error) {
        showFormMessage('Failed to create community. Please try again.', 'error');
        console.error('Error creating community:', error);
    } finally {
        submitBtn.disabled = false;
        submitBtn.innerHTML = originalText;
    }
});

function showFormMessage(message, type) {
    const messageDiv = document.getElementById('form-message');
    messageDiv.textContent = message;
    messageDiv.className = `form-message ${type}`;
    messageDiv.style.display = 'block';
    
    if (type === 'error') {
        setTimeout(() => {
            messageDiv.style.display = 'none';
        }, 5000);
    }
}