// Home page functionality
document.addEventListener('DOMContentLoaded', async function() {
    initializeSearchBar();
    await loadTrendingContent();
});

function initializeSearchBar() {
    const searchInput = document.getElementById('main-search');
    if (searchInput) {
        searchInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                const query = this.value.trim();
                if (query) {
                    window.location.href = `/search?q=${encodeURIComponent(query)}`;
                }
            }
        });
    }
}

async function loadTrendingContent() {
    try {
        const data = await api.getTrending();
        
        displayTrendingContent(data);
        displayPopularMovies(data.popular_movies);
        displayPopularTVShows(data.popular_tv);
        
    } catch (error) {
        console.error('Failed to load trending content:', error);
        showTrendingError();
    }
}

function displayTrendingContent(data) {
    const container = document.getElementById('trending-content');
    if (!container) return;

    const trendingMovies = data.trending_movies || [];
    const trendingTV = data.trending_tv || [];
    
    if (trendingMovies.length === 0 && trendingTV.length === 0) {
        container.innerHTML = createEmptyState('No trending content available');
        return;
    }

    let html = '';
    
    if (trendingMovies.length > 0) {
        html += '<h3>Trending Movies</h3>';
        html += '<div class="media-grid">';
        trendingMovies.slice(0, 6).forEach(movie => {
            html += createMediaCard(movie, 'movie');
        });
        html += '</div>';
    }
    
    if (trendingTV.length > 0) {
        html += '<h3>Trending TV Shows</h3>';
        html += '<div class="media-grid">';
        trendingTV.slice(0, 6).forEach(tv => {
            html += createMediaCard(tv, 'tv');
        });
        html += '</div>';
    }

    container.innerHTML = html;
}

function displayPopularMovies(movies) {
    const container = document.getElementById('popular-movies');
    if (!container || !movies || movies.length === 0) {
        if (container) {
            container.innerHTML = createEmptyState('No popular movies available');
        }
        return;
    }

    let html = '';
    movies.slice(0, 8).forEach(movie => {
        html += createMediaCard(movie, 'movie');
    });

    container.innerHTML = html;
}

function displayPopularTVShows(tvShows) {
    const container = document.getElementById('popular-tv');
    if (!container || !tvShows || tvShows.length === 0) {
        if (container) {
            container.innerHTML = createEmptyState('No popular TV shows available');
        }
        return;
    }

    let html = '';
    tvShows.slice(0, 8).forEach(tv => {
        html += createMediaCard(tv, 'tv');
    });

    container.innerHTML = html;
}

function showTrendingError() {
    const containers = ['trending-content', 'popular-movies', 'popular-tv'];
    
    containers.forEach(containerId => {
        const container = document.getElementById(containerId);
        if (container) {
            container.innerHTML = createEmptyState('Failed to load content. Please try again later.', 'fas fa-exclamation-triangle');
        }
    });
}