// Dashboard page functionality
document.addEventListener('DOMContentLoaded', async function() {
    // Check authentication
    if (!authManager.requireAuth()) {
        return;
    }
    
    await loadDashboardData();
});

async function loadDashboardData() {
    try {
        // Load user stats and recent activity concurrently
        const [stats, recentWatchlist, recommendations] = await Promise.all([
            api.getUserStats(),
            api.getWatchlist('', null), // Get all watchlist items
            api.getRecommendations(6)
        ]);
        
        displayUserStats(stats);
        displayRecentActivity(recentWatchlist.items);
        displayRecommendations(recommendations.recommendations);
        
    } catch (error) {
        console.error('Failed to load dashboard data:', error);
        showDashboardError();
    }
}

function displayUserStats(stats) {
    const container = document.getElementById('stats-grid');
    if (!container) return;
    
    const html = `
        <div class="stat-card">
            <div class="stat-number">${stats.total_movies}</div>
            <div class="stat-label">Movies in Watchlist</div>
        </div>
        <div class="stat-card">
            <div class="stat-number">${stats.total_tv_shows}</div>
            <div class="stat-label">TV Shows in Watchlist</div>
        </div>
        <div class="stat-card">
            <div class="stat-number">${stats.watched_movies}</div>
            <div class="stat-label">Movies Watched</div>
        </div>
        <div class="stat-card">
            <div class="stat-number">${stats.watched_tv_shows}</div>
            <div class="stat-label">TV Shows Watched</div>
        </div>
        <div class="stat-card">
            <div class="stat-number">${stats.average_rating ? stats.average_rating.toFixed(1) : 'N/A'}</div>
            <div class="stat-label">Average Rating</div>
        </div>
    `;
    
    container.innerHTML = html;
}

function displayRecentActivity(watchlistItems) {
    const container = document.getElementById('recent-watchlist');
    if (!container) return;
    
    if (!watchlistItems || watchlistItems.length === 0) {
        container.innerHTML = createEmptyState('No watchlist items yet. Start by searching for movies and TV shows!', 'fas fa-list');
        return;
    }
    
    // Sort by created_at and take the 6 most recent
    const recentItems = watchlistItems
        .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
        .slice(0, 6);
    
    let html = '';
    recentItems.forEach(item => {
        html += createWatchlistCard(item);
    });
    
    container.innerHTML = html;
}

function displayRecommendations(recommendations) {
    const container = document.getElementById('recommendations');
    if (!container) return;
    
    if (!recommendations || recommendations.length === 0) {
        container.innerHTML = createEmptyState('No recommendations yet. Rate some movies and TV shows to get personalized suggestions!', 'fas fa-star');
        return;
    }
    
    let html = '';
    recommendations.forEach(item => {
        html += createWatchlistCard(item);
    });
    
    container.innerHTML = html;
}

function createWatchlistCard(item) {
    const posterURL = getPosterURL(item.poster_path);
    const watchedClass = item.is_watched ? 'watched' : '';
    const ratingDisplay = item.rating ? `★ ${item.rating}/10` : '';
    
    return `
        <div class="media-card ${watchedClass}" data-id="${item.media_id}" data-type="${item.media_type}">
            <div class="media-poster">
                <img src="${posterURL}" alt="${item.title}" loading="lazy" 
                     onerror="this.src='/static/img/no-poster.jpg'">
                ${item.is_watched ? '<div class="watched-badge"><i class="fas fa-check"></i></div>' : ''}
            </div>
            <div class="media-info">
                <div class="media-title">${item.title}</div>
                ${ratingDisplay ? `<div class="media-rating">${ratingDisplay}</div>` : ''}
                <div class="media-actions">
                    <button class="btn btn-small ${item.is_watched ? 'btn-outline' : 'btn-primary'}" 
                            onclick="toggleWatchStatus(${item.media_id}, '${item.media_type}', ${!item.is_watched})">
                        <i class="fas fa-${item.is_watched ? 'undo' : 'check'}"></i>
                        ${item.is_watched ? 'Mark Unwatched' : 'Mark Watched'}
                    </button>
                    ${!item.rating ? `
                        <button class="btn btn-small btn-outline" onclick="showRatingModal(${item.media_id}, '${item.media_type}', '${item.title}')">
                            <i class="fas fa-star"></i> Rate
                        </button>
                    ` : ''}
                </div>
            </div>
        </div>
    `;
}

async function toggleWatchStatus(mediaId, mediaType, isWatched) {
    try {
        await api.updateWatchStatus(mediaId, mediaType, isWatched);
        showSuccess(isWatched ? 'Marked as watched' : 'Marked as unwatched');
        
        // Reload dashboard to update stats
        await loadDashboardData();
        
    } catch (error) {
        console.error('Failed to update watch status:', error);
        showError('Failed to update watch status');
    }
}

function showRatingModal(mediaId, mediaType, title) {
    // Simple prompt for rating - in a real app, this would be a proper modal
    const rating = prompt(`Rate "${title}" (1-10):`);
    
    if (rating) {
        const numRating = parseInt(rating);
        if (numRating >= 1 && numRating <= 10) {
            rateMedia(mediaId, mediaType, numRating);
        } else {
            showError('Please enter a rating between 1 and 10');
        }
    }
}

async function rateMedia(mediaId, mediaType, rating) {
    try {
        await api.rateMedia(mediaId, mediaType, rating);
        showSuccess(`Rating saved: ${rating}/10`);
        
        // Reload dashboard to update display
        await loadDashboardData();
        
    } catch (error) {
        console.error('Failed to rate media:', error);
        showError('Failed to save rating');
    }
}

function showDashboardError() {
    const containers = ['stats-grid', 'recent-watchlist', 'recommendations'];
    
    containers.forEach(containerId => {
        const container = document.getElementById(containerId);
        if (container) {
            container.innerHTML = createEmptyState('Failed to load data. Please refresh the page.', 'fas fa-exclamation-triangle');
        }
    });
}

// Add CSS for watched items
const watchedCSS = `
.media-card.watched {
    opacity: 0.7;
}

.media-card.watched .media-poster::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(70, 211, 105, 0.1);
    z-index: 2;
}

.watched-badge {
    position: absolute;
    top: 10px;
    right: 10px;
    background: var(--success-color);
    color: white;
    border-radius: 50%;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.8rem;
    z-index: 3;
}
`;

// Inject CSS
const style = document.createElement('style');
style.textContent = watchedCSS;
document.head.appendChild(style);