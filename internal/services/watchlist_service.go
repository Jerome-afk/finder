package services

import (
        "crypto/rand"
        "database/sql"
        "encoding/hex"
        "fmt"
        "finderr/internal/models"
)

type WatchlistService struct {
        db *sql.DB
}

func NewWatchlistService(db *sql.DB) *WatchlistService {
        return &WatchlistService{db: db}
}

func generateWatchlistID() string {
        bytes := make([]byte, 16)
        rand.Read(bytes)
        return hex.EncodeToString(bytes)
}

func (s *WatchlistService) AddToWatchlist(userID string, mediaID int, mediaType, title, posterPath string) (*models.WatchlistItem, error) {
        item := &models.WatchlistItem{
                ID:         generateWatchlistID(),
                UserID:     userID,
                MediaID:    mediaID,
                MediaType:  mediaType,
                Title:      title,
                PosterPath: posterPath,
                IsWatched:  false,
        }

        query := `
                INSERT OR REPLACE INTO watchlist_items (id, user_id, media_id, media_type, title, poster_path, is_watched)
                VALUES (?, ?, ?, ?, ?, ?, ?)`

        _, err := s.db.Exec(query, item.ID, item.UserID, item.MediaID, item.MediaType, 
                item.Title, item.PosterPath, item.IsWatched)
        if err != nil {
                return nil, fmt.Errorf("failed to add to watchlist: %w", err)
        }

        // Get the created item
        getQuery := `SELECT created_at, updated_at FROM watchlist_items WHERE id = ?`
        err = s.db.QueryRow(getQuery, item.ID).Scan(&item.CreatedAt, &item.UpdatedAt)
        if err != nil {
                return nil, fmt.Errorf("failed to get created item: %w", err)
        }

        return item, nil
}

func (s *WatchlistService) RemoveFromWatchlist(userID string, mediaID int, mediaType string) error {
        query := `DELETE FROM watchlist_items WHERE user_id = ? AND media_id = ? AND media_type = ?`
        
        result, err := s.db.Exec(query, userID, mediaID, mediaType)
        if err != nil {
                return fmt.Errorf("failed to remove from watchlist: %w", err)
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
                return fmt.Errorf("failed to get rows affected: %w", err)
        }

        if rowsAffected == 0 {
                return fmt.Errorf("item not found in watchlist")
        }

        return nil
}

func (s *WatchlistService) GetUserWatchlist(userID string, mediaType string, isWatched *bool) ([]models.WatchlistItem, error) {
        var items []models.WatchlistItem
        var args []interface{}
        var conditions []string

        query := `
                SELECT id, user_id, media_id, media_type, title, poster_path, is_watched, rating, created_at, updated_at
                FROM watchlist_items WHERE user_id = ?`
        args = append(args, userID)

        if mediaType != "" {
                conditions = append(conditions, "media_type = ?")
                args = append(args, mediaType)
        }

        if isWatched != nil {
                conditions = append(conditions, "is_watched = ?")
                args = append(args, *isWatched)
        }

        if len(conditions) > 0 {
                query += " AND " + conditions[0]
                for i := 1; i < len(conditions); i++ {
                        query += " AND " + conditions[i]
                }
        }

        query += " ORDER BY created_at DESC"

        rows, err := s.db.Query(query, args...)
        if err != nil {
                return nil, fmt.Errorf("failed to get watchlist: %w", err)
        }
        defer rows.Close()

        for rows.Next() {
                var item models.WatchlistItem
                err := rows.Scan(
                        &item.ID, &item.UserID, &item.MediaID, &item.MediaType,
                        &item.Title, &item.PosterPath, &item.IsWatched, &item.Rating,
                        &item.CreatedAt, &item.UpdatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan watchlist item: %w", err)
                }
                items = append(items, item)
        }

        if err = rows.Err(); err != nil {
                return nil, fmt.Errorf("failed to iterate watchlist rows: %w", err)
        }

        return items, nil
}

func (s *WatchlistService) UpdateWatchStatus(userID string, mediaID int, mediaType string, isWatched bool) error {
        query := `
                UPDATE watchlist_items 
                SET is_watched = ?, updated_at = CURRENT_TIMESTAMP
                WHERE user_id = ? AND media_id = ? AND media_type = ?`

        result, err := s.db.Exec(query, isWatched, userID, mediaID, mediaType)
        if err != nil {
                return fmt.Errorf("failed to update watch status: %w", err)
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
                return fmt.Errorf("failed to get rows affected: %w", err)
        }

        if rowsAffected == 0 {
                return fmt.Errorf("item not found in watchlist")
        }

        return nil
}

func (s *WatchlistService) RateMedia(userID string, mediaID int, mediaType string, rating int) error {
        if rating < 1 || rating > 10 {
                return fmt.Errorf("rating must be between 1 and 10")
        }

        query := `
                UPDATE watchlist_items 
                SET rating = ?, updated_at = CURRENT_TIMESTAMP
                WHERE user_id = ? AND media_id = ? AND media_type = ?`

        result, err := s.db.Exec(query, rating, userID, mediaID, mediaType)
        if err != nil {
                return fmt.Errorf("failed to rate media: %w", err)
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
                return fmt.Errorf("failed to get rows affected: %w", err)
        }

        if rowsAffected == 0 {
                return fmt.Errorf("item not found in watchlist")
        }

        return nil
}

func (s *WatchlistService) IsInWatchlist(userID string, mediaID int, mediaType string) (bool, error) {
        var exists bool
        query := `SELECT EXISTS(SELECT 1 FROM watchlist_items WHERE user_id = ? AND media_id = ? AND media_type = ?)`
        
        err := s.db.QueryRow(query, userID, mediaID, mediaType).Scan(&exists)
        if err != nil {
                return false, fmt.Errorf("failed to check watchlist: %w", err)
        }

        return exists, nil
}

func (s *WatchlistService) GetRecommendations(userID string, limit int) ([]models.WatchlistItem, error) {
        // Simple recommendation: get highest rated items from user's watchlist
        // In a real app, this would be more sophisticated
        var items []models.WatchlistItem

        query := `
                SELECT id, user_id, media_id, media_type, title, poster_path, is_watched, rating, created_at, updated_at
                FROM watchlist_items 
                WHERE user_id = ? AND rating IS NOT NULL AND rating >= 8
                ORDER BY rating DESC, created_at DESC
                LIMIT ?`

        rows, err := s.db.Query(query, userID, limit)
        if err != nil {
                return nil, fmt.Errorf("failed to get recommendations: %w", err)
        }
        defer rows.Close()

        for rows.Next() {
                var item models.WatchlistItem
                err := rows.Scan(
                        &item.ID, &item.UserID, &item.MediaID, &item.MediaType,
                        &item.Title, &item.PosterPath, &item.IsWatched, &item.Rating,
                        &item.CreatedAt, &item.UpdatedAt,
                )
                if err != nil {
                        return nil, fmt.Errorf("failed to scan recommendation item: %w", err)
                }
                items = append(items, item)
        }

        if err = rows.Err(); err != nil {
                return nil, fmt.Errorf("failed to iterate recommendation rows: %w", err)
        }

        return items, nil
}

func (s *WatchlistService) GetUserStats(userID string) (*models.UserStats, error) {
        stats := &models.UserStats{}
        
        query := `
                SELECT 
                        COUNT(CASE WHEN media_type = 'movie' THEN 1 END) as total_movies,
                        COUNT(CASE WHEN media_type = 'tv' THEN 1 END) as total_tv_shows,
                        COUNT(CASE WHEN media_type = 'movie' AND is_watched = 1 THEN 1 END) as watched_movies,
                        COUNT(CASE WHEN media_type = 'tv' AND is_watched = 1 THEN 1 END) as watched_tv_shows,
                        COALESCE(AVG(CASE WHEN rating IS NOT NULL THEN rating END), 0) as avg_rating
                FROM watchlist_items 
                WHERE user_id = ?`

        err := s.db.QueryRow(query, userID).Scan(
                &stats.TotalMovies, &stats.TotalTVShows,
                &stats.WatchedMovies, &stats.WatchedTVShows,
                &stats.AverageRating,
        )
        if err != nil {
                return nil, fmt.Errorf("failed to get user stats: %w", err)
        }

        return stats, nil
}