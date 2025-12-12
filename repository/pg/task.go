package pg

import (
	"database/sql"
	"fmt"
	"taskbot/domain"
	"taskbot/repository"
	"time"
)

type TaskRepositoryi interface {
	Create(task domain.Task) (domain.Task, error)
	Update(task domain.Task) (domain.Task, error)
	GetAll(userId int64) ([]domain.Task, error)
	Get(id int64) (domain.Task, error)
	Delete(task domain.Task) error
}

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) Create(task domain.Task) (domain.Task, error) {
	query := `
		INSERT INTO tasks (user_id, title, body, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return task, err
	}
	defer stmt.Close()

	var id int64
	currentTime := time.Now()
	err = stmt.QueryRow(
		task.UserId,
		task.Title,
		task.Body,
		task.Status,
		currentTime,
		currentTime,
	).Scan(&id)

	if err != nil {
		return task, err
	}

	task.Id = id
	task.CreatedAt = currentTime
	task.UpdatedAt = currentTime

	return task, nil
}

func (r *TaskRepository) Update(task domain.Task) (domain.Task, error) {
	query := `
		UPDATE tasks 
		SET title = $1, body = $2, status = $3, updated_at = $4
		WHERE id = $5
	`

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return task, err
	}
	defer stmt.Close()

	currentTime := time.Now()
	res, err := stmt.Exec(
		task.Title,
		task.Body,
		task.Status,
		currentTime,
		task.Id,
	)

	if err != nil {
		return task, err
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return task, err
	}
	if rowAffected == 0 {
		return task, repository.ErrNotFound
	}

	task.UpdatedAt = currentTime
	return task, nil
}

func (r *TaskRepository) GetAll(userId int64) ([]domain.Task, error) {
	query := `
        SELECT id, user_id, title, body, status, created_at, updated_at
        FROM tasks
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task

	for rows.Next() {
		task, err := r.scanRowsToTask(rows)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil
}

func (r *TaskRepository) Get(id int64) (domain.Task, error) {
	query := `
        SELECT id, user_id, title, body, status, created_at, updated_at
        FROM tasks
        WHERE id = $1
    `

	var task domain.Task
	err := r.db.QueryRow(query, id).Scan(
		&task.Id,
		&task.UserId,
		&task.Title,
		&task.Body,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Task{}, repository.ErrNotFound
		}
		return domain.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) Delete(task domain.Task) error {
	query := `
        DELETE FROM tasks
        WHERE id = $1
    `

	res, err := r.db.Exec(query, task.Id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *TaskRepository) scanRowsToTask(rows *sql.Rows) (domain.Task, error) {
	task := domain.Task{}
	err := rows.Scan(
		&task.Id,
		&task.UserId,
		&task.Title,
		&task.Body,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	return task, err
}
