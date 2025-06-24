// Register page functionality
document.addEventListener('DOMContentLoaded', function() {
    initializeRegisterForm();
});

function initializeRegisterForm() {
    const registerForm = document.getElementById('register-form');
    if (registerForm) {
        registerForm.addEventListener('submit', handleRegister);
    }
}

async function handleRegister(e) {
    e.preventDefault();
    
    const username = document.getElementById('username').value.trim();
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value;
    const errorDiv = document.getElementById('register-error');
    
    // Clear previous errors
    errorDiv.style.display = 'none';
    
    // Validate inputs
    const validation = validateRegistrationForm(username, email, password);
    if (!validation.isValid) {
        showError(validation.message, 'register-error');
        return;
    }
    
    // Show loading state
    setLoading('register-form', true);
    
    try {
        const userData = { username, email, password };
        await authManager.register(userData);
        
        // Redirect to dashboard
        window.location.href = '/dashboard';
        
    } catch (error) {
        console.error('Registration failed:', error);
        let errorMessage = 'Registration failed. Please try again.';
        
        if (error.message.includes('already exists') || error.message.includes('duplicate')) {
            errorMessage = 'An account with this email or username already exists. Please try different credentials.';
        } else if (error.message.includes('validation')) {
            errorMessage = 'Please check your input and try again.';
        } else if (error.message.includes('network') || error.message.includes('fetch')) {
            errorMessage = 'Network error. Please check your connection and try again.';
        }
        
        showError(errorMessage, 'register-error');
    } finally {
        setLoading('register-form', false);
    }
}

function validateRegistrationForm(username, email, password) {
    // Username validation
    if (!username) {
        return { isValid: false, message: 'Username is required' };
    }
    if (username.length < 3) {
        return { isValid: false, message: 'Username must be at least 3 characters long' };
    }
    if (username.length > 50) {
        return { isValid: false, message: 'Username must be less than 50 characters' };
    }
    if (!/^[a-zA-Z0-9_]+$/.test(username)) {
        return { isValid: false, message: 'Username can only contain letters, numbers, and underscores' };
    }
    
    // Email validation
    if (!email) {
        return { isValid: false, message: 'Email is required' };
    }
    if (!isValidEmail(email)) {
        return { isValid: false, message: 'Please enter a valid email address' };
    }
    
    // Password validation
    if (!password) {
        return { isValid: false, message: 'Password is required' };
    }
    if (password.length < 6) {
        return { isValid: false, message: 'Password must be at least 6 characters long' };
    }
    if (password.length > 128) {
        return { isValid: false, message: 'Password must be less than 128 characters' };
    }
    
    // Check password strength
    const hasUpperCase = /[A-Z]/.test(password);
    const hasLowerCase = /[a-z]/.test(password);
    const hasNumbers = /\d/.test(password);
    const hasSpecial = /[!@#$%^&*(),.?":{}|<>]/.test(password);
    
    let strengthScore = 0;
    if (hasUpperCase) strengthScore++;
    if (hasLowerCase) strengthScore++;
    if (hasNumbers) strengthScore++;
    if (hasSpecial) strengthScore++;
    
    if (strengthScore < 2) {
        return { 
            isValid: false, 
            message: 'Password should contain at least 2 of: uppercase letters, lowercase letters, numbers, or special characters' 
        };
    }
    
    return { isValid: true };
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