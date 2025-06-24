// Authentication management
class AuthManager {
    constructor() {
        this.currentUser = null;
        this.init();
    }

    async init() {
        await this.checkAuthStatus();
        this.updateNavigation();
    }

    async checkAuthStatus() {
        try {
            const response = await api.getCurrentUser();
            this.currentUser = response.user;
            return true;
        } catch (error) {
            this.currentUser = null;
            return false;
        }
    }

    updateNavigation() {
        const authNav = document.getElementById('auth-nav');
        if (!authNav) return;

        if (this.currentUser) {
            authNav.innerHTML = `
                <a href="/dashboard" class="nav-link">Dashboard</a>
                <a href="/dashboard/watchlist" class="nav-link">Watchlist</a>
                <div class="nav-user">
                    <span class="nav-username">${this.currentUser.username}</span>
                    <button class="btn btn-outline btn-small" onclick="authManager.logout()">
                        <i class="fas fa-sign-out-alt"></i> Logout
                    </button>
                </div>
            `;
        } else {
            authNav.innerHTML = `
                <a href="/login" class="btn btn-outline">Login</a>
                <a href="/register" class="btn btn-primary">Sign Up</a>
            `;
        }
    }

    async login(credentials) {
        try {
            const response = await api.login(credentials);
            this.currentUser = response.user;
            this.updateNavigation();
            return response;
        } catch (error) {
            throw error;
        }
    }

    async register(userData) {
        try {
            const response = await api.register(userData);
            this.currentUser = response.user;
            this.updateNavigation();
            return response;
        } catch (error) {
            throw error;
        }
    }

    async logout() {
        try {
            await api.logout();
            this.currentUser = null;
            this.updateNavigation();
            window.location.href = '/';
        } catch (error) {
            console.error('Logout failed:', error);
            // Force logout on client side even if server request fails
            this.currentUser = null;
            this.updateNavigation();
            window.location.href = '/';
        }
    }

    isAuthenticated() {
        return this.currentUser !== null;
    }

    requireAuth() {
        if (!this.isAuthenticated()) {
            window.location.href = '/login';
            return false;
        }
        return true;
    }
}

// Global auth manager instance
const authManager = new AuthManager();

// Utility function for checking auth status
async function checkAuthStatus() {
    return await authManager.checkAuthStatus();
}

// CSS for user navigation
const userNavCSS = `
.nav-user {
    display: flex;
    align-items: center;
    gap: 1rem;
}

.nav-username {
    color: var(--text-light);
    font-weight: 500;
}

@media (max-width: 768px) {
    .nav-user {
        flex-direction: column;
        gap: 0.5rem;
    }
    
    .nav-username {
        font-size: 0.9rem;
    }
}
`;

// Inject CSS
const style = document.createElement('style');
style.textContent = userNavCSS;
document.head.appendChild(style);