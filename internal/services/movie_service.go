package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"finderr/internal/models"
)

type MovieService struct {
	tmdbAPIKey string
	omdbAPIKey string
	httpClient *http.Client
}

func NewMovieService(tmdbAPIKey, omdbAPIKey string) *MovieService {
	return &MovieService{
		tmdbAPIKey: tmdbAPIKey,
		omdbAPIKey: omdbAPIKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

const (
	tmdbBaseURL = "https://api.themoviedb.org/3"
	omdbBaseURL = "http://www.omdbapi.com"
)

func (s *MovieService) SearchMovies(query string, page int) (*models.TMDBMovieResponse, error) {
	url := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s&page=%d",
		tmdbBaseURL, s.tmdbAPIKey, url.QueryEscape(query), page)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result models.TMDBMovieResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

func (s *MovieService) SearchTVShows(query string, page int) (*models.TMDBTVResponse, error) {
	url := fmt.Sprintf("%s/search/tv?api_key=%s&query=%s&page=%d",
		tmdbBaseURL, s.tmdbAPIKey, url.QueryEscape(query), page)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to search TV shows: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result models.TMDBTVResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

func (s *MovieService) GetMovieDetails(movieID int) (*models.Movie, error) {
	url := fmt.Sprintf("%s/movie/%d?api_key=%s&append_to_response=credits",
		tmdbBaseURL, movieID, s.tmdbAPIKey)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var movie models.Movie
	if err := json.Unmarshal(body, &movie); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Enhance with OMDB data if available
	if s.omdbAPIKey != "" {
		if omdbData, err := s.getOMDBData(movie.Title, "movie"); err == nil {
			s.enhanceMovieWithOMDB(&movie, omdbData)
		}
	}

	return &movie, nil
}

func (s *MovieService) GetTVShowDetails(tvID int) (*models.TVShow, error) {
	url := fmt.Sprintf("%s/tv/%d?api_key=%s&append_to_response=credits",
		tmdbBaseURL, tvID, s.tmdbAPIKey)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get TV show details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tvShow models.TVShow
	if err := json.Unmarshal(body, &tvShow); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Enhance with OMDB data if available
	if s.omdbAPIKey != "" {
		if omdbData, err := s.getOMDBData(tvShow.Name, "series"); err == nil {
			s.enhanceTVShowWithOMDB(&tvShow, omdbData)
		}
	}

	return &tvShow, nil
}

func (s *MovieService) GetTrendingContent() (*models.TrendingContent, error) {
	content := &models.TrendingContent{}

	// Get trending movies
	trendingMoviesURL := fmt.Sprintf("%s/trending/movie/week?api_key=%s", tmdbBaseURL, s.tmdbAPIKey)
	if movieResp, err := s.fetchTMDBMovies(trendingMoviesURL); err == nil {
		content.TrendingMovies = movieResp.Results
	}

	// Get trending TV shows
	trendingTVURL := fmt.Sprintf("%s/trending/tv/week?api_key=%s", tmdbBaseURL, s.tmdbAPIKey)
	if tvResp, err := s.fetchTMDBTVShows(trendingTVURL); err == nil {
		content.TrendingTV = tvResp.Results
	}

	// Get popular movies
	popularMoviesURL := fmt.Sprintf("%s/movie/popular?api_key=%s", tmdbBaseURL, s.tmdbAPIKey)
	if movieResp, err := s.fetchTMDBMovies(popularMoviesURL); err == nil {
		content.PopularMovies = movieResp.Results
	}

	// Get popular TV shows
	popularTVURL := fmt.Sprintf("%s/tv/popular?api_key=%s", tmdbBaseURL, s.tmdbAPIKey)
	if tvResp, err := s.fetchTMDBTVShows(popularTVURL); err == nil {
		content.PopularTV = tvResp.Results
	}

	return content, nil
}

func (s *MovieService) GetMoviesByGenre(genreID int, page int) (*models.TMDBMovieResponse, error) {
	url := fmt.Sprintf("%s/discover/movie?api_key=%s&with_genres=%d&page=%d",
		tmdbBaseURL, s.tmdbAPIKey, genreID, page)
	return s.fetchTMDBMovies(url)
}

func (s *MovieService) GetTVShowsByGenre(genreID int, page int) (*models.TMDBTVResponse, error) {
	url := fmt.Sprintf("%s/discover/tv?api_key=%s&with_genres=%d&page=%d",
		tmdbBaseURL, s.tmdbAPIKey, genreID, page)
	return s.fetchTMDBTVShows(url)
}

func (s *MovieService) fetchTMDBMovies(url string) (*models.TMDBMovieResponse, error) {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result models.TMDBMovieResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

func (s *MovieService) fetchTMDBTVShows(url string) (*models.TMDBTVResponse, error) {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TV shows: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result models.TMDBTVResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

func (s *MovieService) getOMDBData(title, mediaType string) (*models.OMDBResponse, error) {
	url := fmt.Sprintf("%s/?apikey=%s&t=%s&type=%s",
		omdbBaseURL, s.omdbAPIKey, url.QueryEscape(title), mediaType)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get OMDB data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OMDB response: %w", err)
	}

	var omdbResp models.OMDBResponse
	if err := json.Unmarshal(body, &omdbResp); err != nil {
		return nil, fmt.Errorf("failed to parse OMDB response: %w", err)
	}

	if omdbResp.Response == "False" {
		return nil, fmt.Errorf("OMDB error: %s", omdbResp.Error)
	}

	return &omdbResp, nil
}

func (s *MovieService) enhanceMovieWithOMDB(movie *models.Movie, omdbData *models.OMDBResponse) {
	movie.Plot = omdbData.Plot
	movie.Director = omdbData.Director
	movie.IMDBRating = omdbData.ImdbRating

	if omdbData.Actors != "" {
		movie.Cast = strings.Split(omdbData.Actors, ", ")
	}

	if runtime, err := strconv.Atoi(strings.Fields(omdbData.Runtime)[0]); err == nil {
		movie.Runtime = runtime
	}

	// Extract Rotten Tomatoes rating
	for _, rating := range omdbData.Ratings {
		if rating.Source == "Rotten Tomatoes" {
			movie.RottenTomatoes = rating.Value
			break
		}
	}
}

func (s *MovieService) enhanceTVShowWithOMDB(tvShow *models.TVShow, omdbData *models.OMDBResponse) {
	tvShow.Plot = omdbData.Plot
	tvShow.Creator = omdbData.Director // OMDB uses "Director" field for TV show creators
	tvShow.IMDBRating = omdbData.ImdbRating

	if omdbData.Actors != "" {
		tvShow.Cast = strings.Split(omdbData.Actors, ", ")
	}
}
