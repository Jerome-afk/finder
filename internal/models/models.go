package models

import (
        "time"
)

type User struct {
        ID        string    `json:"id" db:"id"`
        Username  string    `json:"username" db:"username"`
        Email     string    `json:"email" db:"email"`
        Password  string    `json:"-" db:"password_hash"`
        CreatedAt time.Time `json:"created_at" db:"created_at"`
        UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Movie struct {
        ID              int      `json:"id"`
        Title           string   `json:"title"`
        Overview        string   `json:"overview"`
        ReleaseDate     string   `json:"release_date"`
        PosterPath      string   `json:"poster_path"`
        BackdropPath    string   `json:"backdrop_path"`
        VoteAverage     float64  `json:"vote_average"`
        VoteCount       int      `json:"vote_count"`
        Genres          []Genre  `json:"genres"`
        Runtime         int      `json:"runtime,omitempty"`
        IMDBRating      string   `json:"imdb_rating,omitempty"`
        RottenTomatoes  string   `json:"rotten_tomatoes,omitempty"`
        Plot            string   `json:"plot,omitempty"`
        Director        string   `json:"director,omitempty"`
        Cast            []string `json:"cast,omitempty"`
        Popularity      float64  `json:"popularity"`
        Adult           bool     `json:"adult"`
        OriginalTitle   string   `json:"original_title"`
        OriginalLanguage string  `json:"original_language"`
}

type TVShow struct {
        ID              int      `json:"id"`
        Name            string   `json:"name"`
        Overview        string   `json:"overview"`
        FirstAirDate    string   `json:"first_air_date"`
        LastAirDate     string   `json:"last_air_date,omitempty"`
        PosterPath      string   `json:"poster_path"`
        BackdropPath    string   `json:"backdrop_path"`
        VoteAverage     float64  `json:"vote_average"`
        VoteCount       int      `json:"vote_count"`
        Genres          []Genre  `json:"genres"`
        NumberOfSeasons int      `json:"number_of_seasons,omitempty"`
        NumberOfEpisodes int     `json:"number_of_episodes,omitempty"`
        IMDBRating      string   `json:"imdb_rating,omitempty"`
        Plot            string   `json:"plot,omitempty"`
        Creator         string   `json:"creator,omitempty"`
        Cast            []string `json:"cast,omitempty"`
        Popularity      float64  `json:"popularity"`
        OriginalName    string   `json:"original_name"`
        OriginalLanguage string  `json:"original_language"`
        Status          string   `json:"status,omitempty"`
}

type Genre struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
}

type WatchlistItem struct {
        ID         string    `json:"id" db:"id"`
        UserID     string    `json:"user_id" db:"user_id"`
        MediaID    int       `json:"media_id" db:"media_id"`
        MediaType  string    `json:"media_type" db:"media_type"` // "movie" or "tv"
        Title      string    `json:"title" db:"title"`
        PosterPath string    `json:"poster_path" db:"poster_path"`
        IsWatched  bool      `json:"is_watched" db:"is_watched"`
        Rating     *int      `json:"rating" db:"rating"` // User's personal rating 1-10
        CreatedAt  time.Time `json:"created_at" db:"created_at"`
        UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type SearchResult struct {
        Movies     []Movie  `json:"movies"`
        TVShows    []TVShow `json:"tv_shows"`
        TotalPages int      `json:"total_pages"`
        Page       int      `json:"page"`
        TotalResults int    `json:"total_results"`
}

type TrendingContent struct {
        TrendingMovies []Movie  `json:"trending_movies"`
        TrendingTV     []TVShow `json:"trending_tv"`
        PopularMovies  []Movie  `json:"popular_movies"`
        PopularTV      []TVShow `json:"popular_tv"`
}

type UserStats struct {
        TotalMovies    int `json:"total_movies"`
        TotalTVShows   int `json:"total_tv_shows"`
        WatchedMovies  int `json:"watched_movies"`
        WatchedTVShows int `json:"watched_tv_shows"`
        AverageRating  float64 `json:"average_rating"`
}

// API Response structures for TMDB
type TMDBMovieResponse struct {
        Page         int     `json:"page"`
        Results      []Movie `json:"results"`
        TotalPages   int     `json:"total_pages"`
        TotalResults int     `json:"total_results"`
}

type TMDBTVResponse struct {
        Page         int      `json:"page"`
        Results      []TVShow `json:"results"`
        TotalPages   int      `json:"total_pages"`
        TotalResults int      `json:"total_results"`
}

// OMDB API Response
type OMDBResponse struct {
        Title          string `json:"Title"`
        Year           string `json:"Year"`
        Rated          string `json:"Rated"`
        Released       string `json:"Released"`
        Runtime        string `json:"Runtime"`
        Genre          string `json:"Genre"`
        Director       string `json:"Director"`
        Writer         string `json:"Writer"`
        Actors         string `json:"Actors"`
        Plot           string `json:"Plot"`
        Language       string `json:"Language"`
        Country        string `json:"Country"`
        Awards         string `json:"Awards"`
        Poster         string `json:"Poster"`
        Ratings        []struct {
                Source string `json:"Source"`
                Value  string `json:"Value"`
        } `json:"Ratings"`
        Metascore    string `json:"Metascore"`
        ImdbRating   string `json:"imdbRating"`
        ImdbVotes    string `json:"imdbVotes"`
        ImdbID       string `json:"imdbID"`
        Type         string `json:"Type"`
        TotalSeasons string `json:"totalSeasons"`
        Response     string `json:"Response"`
        Error        string `json:"Error"`
}