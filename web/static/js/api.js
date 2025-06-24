// API utility functions
const API_BASE = '';

class ApiClient {
    constructor() {
        this.baseURL = API_BASE;
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const config = {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            ...options,
        };

        try {
            const response = await fetch(url, config);
            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || `HTTP error! status: ${response.status}`);
            }

            return data;
        } catch (error) {
            console.error('API request failed:', error);
            throw error;
        }
    }

    // Auth methods
    async register(userData) {
        return this.request('/auth/register', {
            method: 'POST',
            body: JSON.stringify(userData),
        });
    }

    async login(credentials) {
        return this.request('/auth/login', {
            method: 'POST',
            body: JSON.stringify(credentials),
        });
    }

    async logout() {
        return this.request('/auth/logout', {
            method: 'POST',
        });
    }

    async getCurrentUser() {
        return this.request('/auth/user');
    }

    // Movie/TV methods
    async searchMovies(query, page = 1) {
        return this.request(`/api/search/movies?q=${encodeURIComponent(query)}&page=${page}`);
    }

    async searchTVShows(query, page = 1) {
        return this.request(`/api/search/tv?q=${encodeURIComponent(query)}&page=${page}`);
    }

    async search(query, page = 1) {
        return this.request(`/api/search?q=${encodeURIComponent(query)}&page=${page}`);
    }

    async getMovieDetails(movieId) {
        return this.request(`/api/movie/${movieId}`);
    }

    async getTVDetails(tvId) {
        return this.request(`/api/tv/${tvId}`);
    }

    async getTrending() {
        return this.request('/api/trending');
    }

    async getMoviesByGenre(genreId, page = 1) {
        return this.request(`/api/movies/genre/${genreId}?page=${page}`);
    }

    async getTVByGenre(genreId, page = 1) {
        return this.request(`/api/tv/genre/${genreId}?page=${page}`);
    }

    // Watchlist methods
    async getWatchlist(type = '', watched = null) {
        let url = '/api/watchlist';
        const params = new URLSearchParams();
        if (type) params.append('type', type);
        if (watched !== null) params.append('watched', watched);
        if (params.toString()) url += '?' + params.toString();

        return this.request(url);
    }

    async addToWatchlist(mediaData) {
        return this.request('/api/watchlist', {
            method: 'POST',
            body: JSON.stringify(mediaData),
        });
    }

    async removeFromWatchlist(mediaId, mediaType) {
        return this.request(`/api/watchlist/${mediaId}/${mediaType}`, {
            method: 'DELETE',
        });
    }

    async updateWatchStatus(mediaId, mediaType, isWatched) {
        return this.request(`/api/watchlist/${mediaId}/${mediaType}/watched`, {
            method: 'PUT',
            body: JSON.stringify({ is_watched: isWatched }),
        });
    }

    async rateMedia(mediaId, mediaType, rating) {
        return this.request(`/api/watchlist/${mediaId}/${mediaType}/rating`, {
            method: 'PUT',
            body: JSON.stringify({ rating: rating }),
        });
    }

    async checkWatchlist(mediaId, mediaType) {
        return this.request(`/api/watchlist/check/${mediaId}/${mediaType}`);
    }

    async getRecommendations(limit = 10) {
        return this.request(`/api/recommendations?limit=${limit}`);
    }

    async getUserStats() {
        return this.request('/api/stats');
    }
}

// Global API instance
const api = new ApiClient();

// Utility functions
function showError(message, containerId = 'error-container') {
    const container = document.getElementById(containerId);
    if (container) {
        container.innerHTML = `<div class="error-message">${message}</div>`;
        container.style.display = 'block';
    } else {
        console.error('Error:', message);
    }
}

function showSuccess(message, containerId = 'success-container') {
    const container = document.getElementById(containerId);
    if (container) {
        container.innerHTML = `<div class="success-message">${message}</div>`;
        container.style.display = 'block';
    } else {
        console.log('Success:', message);
    }
}

function hideMessage(containerId) {
    const container = document.getElementById(containerId);
    if (container) {
        container.style.display = 'none';
    }
}

function setLoading(buttonId, isLoading) {
    const button = document.getElementById(buttonId);
    if (button) {
        if (isLoading) {
            button.classList.add('loading');
            button.disabled = true;
        } else {
            button.classList.remove('loading');
            button.disabled = false;
        }
    }
}

function formatDate(dateString) {
    if (!dateString) return 'Unknown';
    const date = new Date(dateString);
    return date.getFullYear().toString();
}

function formatRating(rating) {
    if (!rating) return 'N/A';
    return parseFloat(rating).toFixed(1);
}

function getPosterURL(posterPath, size = 'w500') {
    if (!posterPath) return '/static/img/no-poster.jpg';
    return `https://image.tmdb.org/t/p/${size}${posterPath}`;
}

function getBackdropURL(backdropPath, size = 'w1280') {
    if (!backdropPath) return '/static/img/no-backdrop.jpg';
    return `https://image.tmdb.org/t/p/${size}${backdropPath}`;
}

function createMediaCard(media, mediaType) {
    const title = media.title || media.name;
    const releaseDate = media.release_date || media.first_air_date;
    const year = formatDate(releaseDate);
    const rating = formatRating(media.vote_average);
    const posterURL = getPosterURL(media.poster_path);

    return `
        <div class="media-card" data-id="${media.id}" data-type="${mediaType}">
            <div class="media-poster">
                <img src="${posterURL}" alt="${title}" loading="lazy" 
                     onerror="this.src='/static/img/no-poster.jpg'">
            </div>
            <div class="media-info">
                <div class="media-title">${title}</div>
                <div class="media-year">${year}</div>
                <div class="media-rating">
                    <span class="rating-stars">★</span>
                    <span>${rating}</span>
                </div>
                <div class="media-actions">
                    <button class="btn btn-small btn-watchlist" 
                            onclick="toggleWatchlist(${media.id}, '${mediaType}', '${title}', '${media.poster_path || ''}')">
                        <i class="fas fa-plus"></i> Watchlist
                    </button>
                </div>
            </div>
        </div>
    `;
}

function createLoadingSpinner(message = 'Loading...') {
    return `
        <div class="loading-spinner">
            <i class="fas fa-spinner fa-spin"></i>
            <p>${message}</p>
        </div>
    `;
}

function createEmptyState(message, icon = 'fas fa-film') {
    return `
        <div class="empty-state">
            <i class="${icon}"></i>
            <p>${message}</p>
        </div>
    `;
}

// Debounce function for search
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Watchlist management
async function toggleWatchlist(mediaId, mediaType, title, posterPath) {
    try {
        const isLoggedIn = await checkAuthStatus();
        if (!isLoggedIn) {
            window.location.href = '/login';
            return;
        }

        const checkResult = await api.checkWatchlist(mediaId, mediaType);
        
        if (checkResult.in_watchlist) {
            await api.removeFromWatchlist(mediaId, mediaType);
            showSuccess('Removed from watchlist');
        } else {
            await api.addToWatchlist({
                media_id: mediaId,
                media_type: mediaType,
                title: title,
                poster_path: posterPath
            });
            showSuccess('Added to watchlist');
        }
        
        // Update button state
        updateWatchlistButton(mediaId, mediaType, !checkResult.in_watchlist);
        
    } catch (error) {
        showError('Failed to update watchlist');
        console.error('Watchlist error:', error);
    }
}

function updateWatchlistButton(mediaId, mediaType, inWatchlist) {
    const card = document.querySelector(`[data-id="${mediaId}"][data-type="${mediaType}"]`);
    if (card) {
        const button = card.querySelector('.btn-watchlist');
        if (button) {
            if (inWatchlist) {
                button.innerHTML = '<i class="fas fa-check"></i> In Watchlist';
                button.classList.add('active');
            } else {
                button.innerHTML = '<i class="fas fa-plus"></i> Watchlist';
                button.classList.remove('active');
            }
        }
    }
}

// Check authentication status
async function checkAuthStatus() {
    try {
        await api.getCurrentUser();
        return true;
    } catch (error) {
        return false;
    }
}