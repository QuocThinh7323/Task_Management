package models

import (
	"time"
	
	"github.com/jmoiron/sqlx"
)

// Task represents a task in the system
type Task struct {
	ID          int        `db:"id" json:"id"`
	Title       string     `db:"title" json:"title" validate:"required,min=3,max=100"`
	Description string     `db:"description" json:"description"`
	UserID      int        `db:"user_id" json:"user_id"`
	CategoryID  *int       `db:"category_id" json:"category_id"`
	Status      string     `db:"status" json:"status"`
	DueDate     *time.Time `db:"due_date" json:"due_date"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// Category represents a task category
type Category struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name" validate:"required,min=3,max=50"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// TaskRepository handles database operations for tasks
type TaskRepository struct {
	db *sqlx.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create adds a new task to the database
func (r *TaskRepository) Create(task *Task) error {
	query := `
		INSERT INTO tasks (title, description, user_id, category_id, status, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	
	return r.db.QueryRowx(
		query,
		task.Title,
		task.Description,
		task.UserID,
		task.CategoryID,
		task.Status,
		task.DueDate,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

// Update modifies an existing task
func (r *TaskRepository) Update(task *Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, category_id = $3, status = $4, due_date = $5, updated_at = NOW()
		WHERE id = $6 AND user_id = $7
		RETURNING updated_at
	`
	
	return r.db.QueryRowx(
		query,
		task.Title,
		task.Description,
		task.CategoryID,
		task.Status,
		task.DueDate,
		task.ID,
		task.UserID,
	).Scan(&task.UpdatedAt)
}

// Delete removes a task by ID
func (r *TaskRepository) Delete(id, userID int) error {
	_, err := r.db.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", id, userID)
	return err
}

// FindByID finds a task by ID
func (r *TaskRepository) FindByID(id int) (*Task, error) {
	task := &Task{}
	err := r.db.Get(task, "SELECT * FROM tasks WHERE id = $1", id)
	return task, err
}

// ListByUser returns all tasks for a specific user
func (r *TaskRepository) ListByUser(userID int) ([]Task, error) {
	var tasks []Task
	err := r.db.Select(&tasks, "SELECT * FROM tasks WHERE user_id = $1 ORDER BY due_date ASC", userID)
	return tasks, err
}

// ListAllTasks returns all tasks (admin only)
func (r *TaskRepository) ListAll() ([]Task, error) {
	var tasks []Task
	err := r.db.Select(&tasks, "SELECT * FROM tasks ORDER BY due_date ASC")
	return tasks, err
}

// CategoryRepository handles database operations for categories
type CategoryRepository struct {
	db *sqlx.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create adds a new category
func (r *CategoryRepository) Create(category *Category) error {
	query := `
		INSERT INTO categories (name, created_at)
		VALUES ($1, NOW())
		RETURNING id, created_at
	`
	
	return r.db.QueryRowx(query, category.Name).Scan(&category.ID, &category.CreatedAt)
}

// List returns all categories
func (r *CategoryRepository) List() ([]Category, error) {
	var categories []Category
	err := r.db.Select(&categories, "SELECT * FROM categories ORDER BY name")
	return categories, err
}

// Delete removes a category
func (r *CategoryRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	return err
}