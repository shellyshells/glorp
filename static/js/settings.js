    // Settings navigation
    document.querySelectorAll('.settings-nav-item').forEach(item => {
        item.addEventListener('click', function() {
            // Remove active class from all nav items and sections
            document.querySelectorAll('.settings-nav-item').forEach(nav => nav.classList.remove('active'));
            document.querySelectorAll('.settings-section').forEach(section => {
                section.classList.remove('active');
                section.style.display = 'none';
            });
            
            // Add active class to clicked nav item
            this.classList.add('active');
            
            // Show corresponding section
            const sectionId = this.dataset.section + '-section';
            const section = document.getElementById(sectionId);
            section.classList.add('active');
            section.style.display = 'block';
        });
    });
    
    // Avatar style selection
    function selectAvatarStyle(style) {
        // Remove selected class from all options
        document.querySelectorAll('.avatar-style-option').forEach(option => {
            option.classList.remove('selected');
        });
        
        // Add selected class to clicked option
        event.currentTarget.classList.add('selected');
        
        // Send update to server
        updateAvatarStyle(style);
    }
    
    async function updateAvatarStyle(style) {
        try {
            const response = await fetch('/api/profile/avatar-style', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ style: style })
            });
            
            if (response.ok) {
                // Reload the page to ensure the new avatar style is applied everywhere
                location.reload();
            } else {
                const result = await response.text();
                showFormMessage('avatar-message', result, 'error');
            }
        } catch (error) {
            console.error('Error:', error);
            showFormMessage('avatar-message', 'Failed to update avatar style. Please try again.', 'error');
        }
    }
    
    function getAvatarClass(style) {
        switch (style) {
            case 'red': return 'avatar-red';
            case 'blue': return 'avatar-blue';
            case 'green': return 'avatar-green';
            case 'purple': return 'avatar-purple';
            case 'orange': return 'avatar-orange';
            case 'pink': return 'avatar-pink';
            case 'teal': return 'avatar-teal';
            case 'admin': return 'avatar-admin';
            default: return 'avatar-default';
        }
    }
    
    // Profile form handler
    document.getElementById('profile-form')?.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const formData = new FormData(this);
        const data = {
            display_name: formData.get('display_name'),
            bio: formData.get('bio'),
            location: formData.get('location'),
            website: formData.get('website'),
            show_email: document.getElementById('show-email')?.checked || false,
            show_online: document.getElementById('show-online')?.checked || false,
            allow_messages: document.getElementById('allow-messages')?.checked || false,
            public_profile: document.getElementById('public-profile')?.checked || false
        };
        
        const submitBtn = this.querySelector('button[type="submit"]');
        const originalText = submitBtn.innerHTML;
        submitBtn.disabled = true;
        submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Saving...';
        
        try {
            const response = await fetch('/api/profile/update', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });
            
            if (response.ok) {
                showFormMessage('profile-message', 'Profile updated successfully!', 'success');
            } else {
                const result = await response.text();
                showFormMessage('profile-message', result, 'error');
            }
        } catch (error) {
            console.error('Error:', error);
            showFormMessage('profile-message', 'Failed to update profile. Please try again.', 'error');
        } finally {
            submitBtn.disabled = false;
            submitBtn.innerHTML = originalText;
        }
    });
    
    // Privacy form handler
    document.getElementById('privacy-form')?.addEventListener('submit', async function(e) {
        e.preventDefault();
        
        const formData = new FormData(this);
        const data = {
            display_name: document.getElementById('display-name')?.value || '{{.User.Username}}',
            bio: document.getElementById('bio')?.value || '',
            location: document.getElementById('location')?.value || '',
            website: document.getElementById('website')?.value || '',
            show_email: formData.has('show_email'),
            show_online: formData.has('show_online'),
            allow_messages: formData.has('allow_messages'),
            public_profile: formData.has('public_profile')
        };
        
        const submitBtn = this.querySelector('button[type="submit"]');
        const originalText = submitBtn.innerHTML;
        submitBtn.disabled = true;
        submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Saving...';
        
        try {
            const response = await fetch('/api/profile/update', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });
            
            if (response.ok) {
                showFormMessage('privacy-message', 'Privacy settings updated successfully!', 'success');
            } else {
                const result = await response.text();
                showFormMessage('privacy-message', result, 'error');
            }
        } catch (error) {
            console.error('Error:', error);
            showFormMessage('privacy-message', 'Failed to update privacy settings. Please try again.', 'error');
        } finally {
            submitBtn.disabled = false;
            submitBtn.innerHTML = originalText;
        }
    });
    
    // Other form handlers (placeholder functionality)
    document.getElementById('account-form')?.addEventListener('submit', function(e) {
        e.preventDefault();
        showFormMessage('account-message', 'Account settings saved!', 'success');
    });
    
    document.getElementById('notifications-form')?.addEventListener('submit', function(e) {
        e.preventDefault();
        showFormMessage('notifications-message', 'Notification settings saved!', 'success');
    });
    
    document.getElementById('security-form')?.addEventListener('submit', function(e) {
        e.preventDefault();
        showFormMessage('security-message', 'Password change feature coming soon!', 'success');
    });
    
    function showFormMessage(messageId, text, type) {
        const messageDiv = document.getElementById(messageId);
        messageDiv.textContent = text;
        messageDiv.className = `form-message ${type}`;
        messageDiv.style.display = 'block';
        
        setTimeout(() => {
            messageDiv.style.display = 'none';
        }, 5000);
    }
    
    function confirmDeleteAccount() {
        const confirmed = confirm('Are you sure you want to delete your account? This action cannot be undone.');
        if (confirmed) {
            const doubleConfirm = prompt('Type "DELETE" to confirm account deletion:');
            if (doubleConfirm === 'DELETE') {
                alert('Account deletion feature will be implemented soon.');
            }
        }
    }