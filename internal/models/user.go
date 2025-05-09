package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// User represents a system user
type User struct {
	ID           int       `db:"id" json:"id"`
	Username     string    `db:"username" json:"username" validate:"required,min=3,max=50"`
	Email        string    `db:"email" json:"email" validate:"required,email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// LoginRequest for user authentication
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest for user registration
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// UserRepository handles database operations for users
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create adds a new user to the database
func (r *UserRepository) Create(user *User, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	user.PasswordHash = string(hashedPassword)
	
	// Insert the user
	query := `
		INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	
	return r.db.QueryRowx(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(username string) (*User, error) {
    user := &User{}
    err := r.db.QueryRow("SELECT id, username, email, password_hash, role FROM users WHERE username = $1", username).
        Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)
    
    if err == sql.ErrNoRows {
        // Không tìm thấy username, trả về nil, nil thay vì lỗi
        return nil, nil
    }
    
    if err != nil {
        // Lỗi database thực sự
        return nil, err
    }
    
    return user, nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id int) (*User, error) {
	user := &User{}
	err := r.db.Get(user, "SELECT * FROM users WHERE id = $1", id)
	return user, err
}

// List returns all users
func (r *UserRepository) List() ([]User, error) {
	var users []User
	err := r.db.Select(&users, "SELECT id, username, email, role, created_at, updated_at FROM users")
	return users, err
}

// CheckPassword verifies a user's password
func (r *UserRepository) CheckPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}