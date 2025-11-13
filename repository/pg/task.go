package pg

import (
	"database/sql"
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
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE user_id = $1", userId).Scan(&count)
	if err != nil {
		return []domain.Task{}, err
	}

	query := `
		SELECT id, user_id, title, body, status, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
	`
	rows, err := r.db.Query(query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return []domain.Task{}, repository.ErrNotFound
		}

		return []domain.Task{}, err
	}

	tasks := make([]domain.Task, 0, count)

	defer rows.Close()
	for rows.Next() {
		t, err := r.scanRowsToTask(rows)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TaskRepository) Get(id int64) (domain.Task, error) {
	query := `
		SELECT id, user_id, title, body, status, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`
	rows, err := r.db.Query(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Task{}, repository.ErrNotFound
		}

		return domain.Task{}, err
	}

	defer rows.Close()
	rows.Next()

	return r.scanRowsToTask(rows)
}

func (r *TaskRepository) Delete(task domain.Task) error {
	query := `
		DELETE
		FROM tasks
		WHERE id = $1
	`
	res, err := r.db.Exec(query, task.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return repository.ErrNotFound
		}

		return err
	}

	rowAffected, err := res.RowsAffected()
	if rowAffected == 0 || err != nil {
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
