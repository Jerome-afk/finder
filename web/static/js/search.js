// Search page functionality
let currentQuery = '';
let currentPage = 1;
let currentType = 'all';
let searchTimeout;

document.addEventListener('DOMContentLoaded', function() {
    initializeSearchPage();
});

function initializeSearchPage() {
    const searchInput = document.getElementById('search-input');
    const filterTabs = document.querySelectorAll('.filter-tab');
    
    // Get initial query from URL
    const urlParams = new URLSearchParams(window.location.search);
    currentQuery = urlParams.get('q') || '';
    
    if (searchInput) {
        searchInput.value = currentQuery;
        searchInput.addEventListener('input', debounceSearch);
        searchInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                performSearch();
            }
        });
    }
    
    // Setup filter tabs
    filterTabs.forEach(tab => {
        tab.addEventListener('click', function() {
            switchFilter(this.dataset.type);
        });
    });
    
    // Perform initial search if there's a query
    if (currentQuery) {
        performSearch();
    }
}

const debounceSearch = debounce(function() {
    const searchInput = document.getElementById('search-input');
    if (searchInput) {
        currentQuery = searchInput.value.trim();
        currentPage = 1;
        if (currentQuery) {
            performSearch();
            updateURL();
        } else {
            clearResults();
        }
    }
}, 500);

function switchFilter(type) {
    currentType = type;
    currentPage = 1;
    
    // Update active tab
    document.querySelectorAll('.filter-tab').forEach(tab => {
        tab.classList.remove('active');
    });
    document.querySelector(`[data-type="${type}"]`).classList.add('active');
    
    if (currentQuery) {
        performSearch();
    }
}

async function performSearch() {
    if (!currentQuery) return;
    
    const resultsContainer = document.getElementById('search-results');
    const paginationContainer = document.getElementById('pagination');
    
    // Show loading
    resultsContainer.innerHTML = createLoadingSpinner('Searching...');
    paginationContainer.style.display = 'none';
    
    try {
        let results;
        
        if (currentType === 'movie') {
            results = await api.searchMovies(currentQuery, currentPage);
            displayMovieResults(results);
        } else if (currentType === 'tv') {
            results = await api.searchTVShows(currentQuery, currentPage);
            displayTVResults(results);
        } else {
            // Search both
            results = await api.search(currentQuery, currentPage);
            displayMixedResults(results);
        }
        
    } catch (error) {
        console.error('Search failed:', error);
        resultsContainer.innerHTML = createEmptyState('Search failed. Please try again.', 'fas fa-exclamation-triangle');
    }
}

function displayMovieResults(results) {
    const container = document.getElementById('search-results');
    
    if (results.results.length === 0) {
        container.innerHTML = createEmptyState('No movies found for your search.', 'fas fa-film');
        return;
    }
    
    let html = '<div class="media-grid">';
    results.results.forEach(movie => {
        html += createMediaCard(movie, 'movie');
    });
    html += '</div>';
    
    container.innerHTML = html;
    displayPagination(results.page, results.total_pages);
}

function displayTVResults(results) {
    const container = document.getElementById('search-results');
    
    if (results.results.length === 0) {
        container.innerHTML = createEmptyState('No TV shows found for your search.', 'fas fa-tv');
        return;
    }
    
    let html = '<div class="media-grid">';
    results.results.forEach(tv => {
        html += createMediaCard(tv, 'tv');
    });
    html += '</div>';
    
    container.innerHTML = html;
    displayPagination(results.page, results.total_pages);
}

function displayMixedResults(results) {
    const container = document.getElementById('search-results');
    
    const hasMovies = results.movies && results.movies.results && results.movies.results.length > 0;
    const hasTVShows = results.tv_shows && results.tv_shows.results && results.tv_shows.results.length > 0;
    
    if (!hasMovies && !hasTVShows) {
        container.innerHTML = createEmptyState('No results found for your search.', 'fas fa-search');
        return;
    }
    
    let html = '';
    
    if (hasMovies) {
        html += '<section class="search-section">';
        html += '<h2>Movies</h2>';
        html += '<div class="media-grid">';
        results.movies.results.slice(0, 10).forEach(movie => {
            html += createMediaCard(movie, 'movie');
        });
        html += '</div>';
        html += '</section>';
    }
    
    if (hasTVShows) {
        html += '<section class="search-section">';
        html += '<h2>TV Shows</h2>';
        html += '<div class="media-grid">';
        results.tv_shows.results.slice(0, 10).forEach(tv => {
            html += createMediaCard(tv, 'tv');
        });
        html += '</div>';
        html += '</section>';
    }
    
    container.innerHTML = html;
    
    // For mixed results, use the first available pagination
    const paginationData = hasMovies ? results.movies : results.tv_shows;
    displayPagination(paginationData.page, paginationData.total_pages);
}

function displayPagination(currentPage, totalPages) {
    const container = document.getElementById('pagination');
    
    if (totalPages <= 1) {
        container.style.display = 'none';
        return;
    }
    
    let html = '';
    
    // Previous button
    html += `<button ${currentPage === 1 ? 'disabled' : ''} onclick="goToPage(${currentPage - 1})">
        <i class="fas fa-chevron-left"></i>
    </button>`;
    
    // Page numbers
    const startPage = Math.max(1, currentPage - 2);
    const endPage = Math.min(totalPages, currentPage + 2);
    
    if (startPage > 1) {
        html += `<button onclick="goToPage(1)">1</button>`;
        if (startPage > 2) {
            html += `<span>...</span>`;
        }
    }
    
    for (let i = startPage; i <= endPage; i++) {
        html += `<button ${i === currentPage ? 'class="active"' : ''} onclick="goToPage(${i})">${i}</button>`;
    }
    
    if (endPage < totalPages) {
        if (endPage < totalPages - 1) {
            html += `<span>...</span>`;
        }
        html += `<button onclick="goToPage(${totalPages})">${totalPages}</button>`;
    }
    
    // Next button
    html += `<button ${currentPage === totalPages ? 'disabled' : ''} onclick="goToPage(${currentPage + 1})">
        <i class="fas fa-chevron-right"></i>
    </button>`;
    
    container.innerHTML = html;
    container.style.display = 'flex';
}

function goToPage(page) {
    currentPage = page;
    performSearch();
    window.scrollTo({ top: 0, behavior: 'smooth' });
}

function clearResults() {
    const resultsContainer = document.getElementById('search-results');
    const paginationContainer = document.getElementById('pagination');
    
    resultsContainer.innerHTML = createEmptyState('Enter a search term to find movies and TV shows', 'fas fa-search');
    paginationContainer.style.display = 'none';
}

function updateURL() {
    const url = new URL(window.location);
    if (currentQuery) {
        url.searchParams.set('q', currentQuery);
    } else {
        url.searchParams.delete('q');
    }
    window.history.replaceState({}, '', url);
}