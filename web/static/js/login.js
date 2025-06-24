// Login page functionality
document.addEventListener('DOMContentLoaded', function() {
    initializeLoginForm();
});

function initializeLoginForm() {
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', handleLogin);
    }
}

async function handleLogin(e) {
    e.preventDefault();
    
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value;
    const errorDiv = document.getElementById('login-error');
    
    // Clear previous errors
    errorDiv.style.display = 'none';
    
    // Validate inputs
    if (!email || !password) {
        showError('Please fill in all fields', 'login-error');
        return;
    }
    
    if (!isValidEmail(email)) {
        showError('Please enter a valid email address', 'login-error');
        return;
    }
    
    // Show loading state
    setLoading('login-form', true);
    
    try {
        const credentials = { email, password };
        await authManager.login(credentials);
        
        // Redirect to dashboard or intended page
        const urlParams = new URLSearchParams(window.location.search);
        const redirectTo = urlParams.get('redirect') || '/dashboard';
        window.location.href = redirectTo;
        
    } catch (error) {
        console.error('Login failed:', error);
        let errorMessage = 'Login failed. Please try again.';
        
        if (error.message.includes('credentials')) {
            errorMessage = 'Invalid email or password. Please check your credentials and try again.';
        } else if (error.message.includes('network') || error.message.includes('fetch')) {
            errorMessage = 'Network error. Please check your connection and try again.';
        }
        
        showError(errorMessage, 'login-error');
    } finally {
        setLoading('login-form', false);
    }
}

function isValidEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

function setLoading(formId, isLoading) {
    const form = document.getElementById(formId);
    const submitButton = form.querySelector('button[type="submit"]');
    
    if (isLoading) {
        submitButton.classList.add('loading');
        submitButton.disabled = true;
    } else {
        submitButton.classList.remove('loading');
        submitButton.disabled = false;
    }
}