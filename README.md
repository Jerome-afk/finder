# FindeRR - Movie/TV Show Discovery Web App

A comprehensive entertainment discovery platform built with Go, where users can search for movies and TV shows, view detailed information, manage personal watchlists, and discover trending content.

## Features

### Core Functionality
- **Search** - Real-time search for movies and TV shows with pagination
- **Detailed Views** - Comprehensive information including cast, ratings, and plot
- **Watchlist Management** - Add/remove titles, mark as watched, rate content
- **Trending Content** - Popular movies and shows from TMDB
- **User Dashboard** - Personal stats and recommendations
- **Genre Filtering** - Browse content by category

### User Experience
- **Authentication** - Secure user registration and login
- **Responsive Design** - Works on desktop and mobile
- **Real-time Updates** - Instant watchlist updates
- **Error Handling** - Graceful error management with user feedback
- **Performance** - Debounced search and efficient API calls

## Tech Stack

### Backend
- **Go** - Main server language
- **Gin** - Web framework
- **PostgreSQL** - Database with UUID support
- **TMDB API** - Primary movie/TV data source
- **OMDB API** - Additional ratings and details

### Frontend
- **HTML5** - Semantic markup
- **CSS3** - Modern styling with Netflix-inspired design
- **Vanilla JavaScript** - Client-side functionality
- **Font Awesome** - Icons

## API Integration

### TMDB API
- Movie and TV show search
- Trending content
- Detailed information
- Genre-based filtering
- High-quality poster images

### OMDB API
- IMDB ratings
- Rotten Tomatoes scores
- Additional plot details
- Cast information

## Installation

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 12 or higher
- TMDB API key
- OMDB API key (optional)

### Setup Steps

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd finderr
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Configure database**
   - Create a PostgreSQL database
   - Update `DATABASE_URL` in .env

5. **Get API keys**
   - TMDB: Register at https://www.themoviedb.org/settings/api
   - OMDB: Register at http://www.omdbapi.com/apikey.aspx
   - Add keys to .env file

6. **Run migrations**
   ```bash
   go run main.go
   ```
   Migrations run automatically on startup.

7. **Start the server**
   ```bash
   go run main.go
   ```

The application will be available at `http://localhost:8080`

## Project Structure

```
finderr/
├── main.go                 # Application entry point
├── internal/
│   ├── api/               # HTTP handlers
│   │   ├── auth_handler.go
│   │   ├── movie_handler.go
│   │   ├── watchlist_handler.go
│   │   ├── web_handler.go
│   │   └── routes.go
│   ├── config/            # Configuration management
│   │   └── config.go
│   ├── database/          # Database setup and migrations
│   │   └── database.go
│   ├── middleware/        # HTTP middleware
│   │   └── middleware.go
│   ├── models/           # Data models
│   │   └── models.go
│   └── services/         # Business logic
│       ├── user_service.go
│       ├── movie_service.go
│       └── watchlist_service.go
├── migrations/           # Database migrations
│   ├── 001_initial_schema.up.sql
│   └── 001_initial_schema.down.sql
├── web/                 # Frontend assets
│   ├── static/
│   │   ├── css/
│   │   │   └── style.css
│   │   └── js/
│   │       ├── api.js
│   │       ├── auth.js
│   │       ├── home.js
│   │       ├── search.js
│   │       ├── login.js
│   │       ├── register.js
│   │       └── dashboard.js
│   └── templates/       # HTML templates
│       ├── base.html
│       ├── index.html
│       ├── search.html
│       ├── login.html
│       ├── register.html
│       └── dashboard.html
└── go.mod
```

## API Endpoints

### Authentication
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout
- `GET /auth/user` - Get current user

### Movies & TV Shows
- `GET /api/search` - Search movies and TV shows
- `GET /api/search/movies` - Search movies only
- `GET /api/search/tv` - Search TV shows only
- `GET /api/movie/:id` - Get movie details
- `GET /api/tv/:id` - Get TV show details
- `GET /api/trending` - Get trending content
- `GET /api/movies/genre/:id` - Get movies by genre
- `GET /api/tv/genre/:id` - Get TV shows by genre

### Watchlist (Protected)
- `GET /api/watchlist` - Get user's watchlist
- `POST /api/watchlist` - Add to watchlist
- `DELETE /api/watchlist/:mediaId/:mediaType` - Remove from watchlist
- `PUT /api/watchlist/:mediaId/:mediaType/watched` - Update watch status
- `PUT /api/watchlist/:mediaId/:mediaType/rating` - Rate media
- `GET /api/watchlist/check/:mediaId/:mediaType` - Check if in watchlist
- `GET /api/recommendations` - Get recommendations
- `GET /api/stats` - Get user statistics

## Database Schema

### Users Table
- `id` - UUID primary key
- `username` - Unique username
- `email` - Unique email address
- `password_hash` - Bcrypt hashed password
- `created_at`, `updated_at` - Timestamps

### Watchlist Items Table
- `id` - UUID primary key
- `user_id` - Foreign key to users
- `media_id` - TMDB ID
- `media_type` - 'movie' or 'tv'
- `title` - Display title
- `poster_path` - Image path
- `is_watched` - Watched status
- `rating` - User rating (1-10)
- `created_at`, `updated_at` - Timestamps

### Sessions Table
- `id` - Session ID
- `data` - Session data
- `expires_at` - Expiration timestamp

## Security Features

- **Password Hashing** - Bcrypt for secure password storage
- **Session Management** - Secure session handling
- **Input Validation** - Server-side validation for all inputs
- **SQL Injection Protection** - Parameterized queries
- **XSS Protection** - Content-Type headers and input sanitization
- **CSRF Protection** - Session-based CSRF tokens

## Performance Optimizations

- **Debounced Search** - Reduces API calls
- **Concurrent Requests** - Parallel API calls where possible
- **Database Indexing** - Optimized queries
- **Image Lazy Loading** - Improved page load times
- **Caching Headers** - Static asset caching

## Error Handling

- **API Rate Limiting** - Graceful handling of rate limits
- **Network Errors** - Retry logic and user feedback
- **Database Errors** - Transaction rollbacks
- **Validation Errors** - Clear user messaging
- **404 Handling** - Proper error pages

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:
1. Check the existing issues on GitHub
2. Create a new issue with detailed information
3. Include steps to reproduce any bugs

## Roadmap

- [ ] Advanced filtering (year, rating, runtime)
- [ ] Dark/light theme toggle
- [ ] Social features and sharing
- [ ] Watch provider integration
- [ ] Trailer integration
- [ ] Export functionality
- [ ] Mobile app
- [ ] API rate limiting improvements
- [ ] Caching layer (Redis)
- [ ] Full-text search
- [ ] Recommendation improvements