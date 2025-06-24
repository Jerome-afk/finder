package services

import (
        "crypto/rand"
        "database/sql"
        "encoding/hex"
        "fmt"
        "finderr/internal/models"

        "golang.org/x/crypto/bcrypt"
)

type UserService struct {
        db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
        return &UserService{db: db}
}

func generateID() string {
        bytes := make([]byte, 16)
        rand.Read(bytes)
        return hex.EncodeToString(bytes)
}

func (s *UserService) CreateUser(username, email, password string) (*models.User, error) {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
                return nil, fmt.Errorf("failed to hash password: %w", err)
        }

        user := &models.User{
                ID:       generateID(),
                Username: username,
                Email:    email,
                Password: string(hashedPassword),
        }

        query := `
                INSERT INTO users (id, username, email, password_hash)
                VALUES (?, ?, ?, ?)
                RETURNING created_at, updated_at`

        err = s.db.QueryRow(query, user.ID, user.Username, user.Email, user.Password).
                Scan(&user.CreatedAt, &user.UpdatedAt)
        if err != nil {
                return nil, fmt.Errorf("failed to create user: %w", err)
        }

        return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
        user := &models.User{}
        query := `
                SELECT id, username, email, password_hash, created_at, updated_at
                FROM users WHERE email = ?`

        err := s.db.QueryRow(query, email).Scan(
                &user.ID, &user.Username, &user.Email, &user.Password,
                &user.CreatedAt, &user.UpdatedAt,
        )
        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, fmt.Errorf("user not found")
                }
                return nil, fmt.Errorf("failed to get user: %w", err)
        }

        return user, nil
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
        user := &models.User{}
        query := `
                SELECT id, username, email, password_hash, created_at, updated_at
                FROM users WHERE id = ?`

        err := s.db.QueryRow(query, id).Scan(
                &user.ID, &user.Username, &user.Email, &user.Password,
                &user.CreatedAt, &user.UpdatedAt,
        )
        if err != nil {
                if err == sql.ErrNoRows {
                        return nil, fmt.Errorf("user not found")
                }
                return nil, fmt.Errorf("failed to get user: %w", err)
        }

        return user, nil
}

func (s *UserService) ValidatePassword(hashedPassword, password string) error {
        return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

