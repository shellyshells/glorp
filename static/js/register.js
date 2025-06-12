document.getElementById('register-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    const formData = new FormData(this);
    const password = formData.get('password');
    const confirmPassword = formData.get('confirm-password');
    
    const errorDiv = document.getElementById('error-message');
    
    if (password !== confirmPassword) {
        errorDiv.textContent = 'Passwords do not match';
        errorDiv.style.display = 'block';
        return;
    }
    
    const data = {
        username: formData.get('username'),
        email: formData.get('email'),
        password: password
    };
    
    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });
        
        if (response.ok) {
            window.location.href = '/';
        } else {
            const result = await response.text();
            errorDiv.textContent = result;
            errorDiv.style.display = 'block';
        }
    } catch (error) {
        errorDiv.textContent = 'Registration failed. Please try again.';
        errorDiv.style.display = 'block';
    }
});