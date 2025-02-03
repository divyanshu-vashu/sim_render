document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const errorDiv = document.getElementById('errorMessage');
    errorDiv.style.display = 'none';
    
    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            credentials: 'same-origin',
            body: JSON.stringify({
                username: document.getElementById('username').value,
                password: document.getElementById('password').value,
            }),
        });

        const data = await response.json();

        if (response.ok && data.success) {
            console.log('Login successful, redirecting...');
            window.location.replace('/');
        } else {
            errorDiv.textContent = data.error || 'Login failed';
            errorDiv.style.display = 'block';
        }
    } catch (error) {
        console.error('Login error:', error);
        errorDiv.textContent = 'An error occurred. Please try again.';
        errorDiv.style.display = 'block';
    }
});