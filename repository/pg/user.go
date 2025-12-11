package pg

import (
	"database/sql"
	"taskbot/domain"
	"taskbot/repository"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(user domain.User) (domain.User, error) {
	if user.Id == 0 {
		return r.Create(user)
	}
	return r.Update(user)
}

func (r *UserRepository) Create(user domain.User) (domain.User, error) {
	taken, err := r.isUsernameTaken(user.Username, 0)
	if err != nil {
		return domain.User{}, err
	}
	if taken {
		return domain.User{}, repository.ErrUsernameTaken
	}

	query := `
		INSERT INTO users (username, first_name, pass, telegram_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	currentTime := time.Now()
	err = r.db.QueryRow(
		query,
		user.Username,
		user.Name,
		user.Password,
		user.TelegramId,
		currentTime,
		currentTime,
	).Scan(&user.Id)

	if err != nil {
		return domain.User{}, err
	}

	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	return user, nil
}

func (r *UserRepository) Update(user domain.User) (domain.User, error) {
	_, err := r.Get(user.Id)
	if err != nil {
		return domain.User{}, err
	}

	taken, err := r.isUsernameTaken(user.Username, user.Id)
	if err != nil {
		return domain.User{}, err
	}
	if taken {
		return domain.User{}, repository.ErrUsernameTaken
	}

	query := `
		UPDATE users 
		SET username = $1, first_name = $2, pass = $3, telegram_id = $4, updated_at = $5
		WHERE id = $6
	`

	currentTime := time.Now()
	result, err := r.db.Exec(
		query,
		user.Username,
		user.Name,
		user.Password,
		user.TelegramId,
		currentTime,
		user.Id,
	)

	if err != nil {
		return domain.User{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.User{}, err
	}
	if rowsAffected == 0 {
		return domain.User{}, repository.ErrNotFound
	}

	user.UpdatedAt = currentTime

	return user, nil
}

func (r *UserRepository) Get(id int64) (domain.User, error) {
	query := `
		SELECT id, username, first_name, pass, telegram_id, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRow(query, id).Scan(
		&user.Id,
		&user.Username,
		&user.Name,
		&user.Password,
		&user.TelegramId,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, repository.ErrNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByTelegramId(tid int64) (domain.User, error) {
	query := `
		SELECT id, username, first_name, pass, telegram_id, created_at, updated_at
		FROM users 
		WHERE telegram_id = $1
	`

	var user domain.User
	err := r.db.QueryRow(query, tid).Scan(
		&user.Id,
		&user.Username,
		&user.Name,
		&user.Password,
		&user.TelegramId,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, repository.ErrNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetByUsername(username string) (domain.User, error) {
	query := `
		SELECT id, username, first_name, pass, telegram_id, created_at, updated_at
		FROM users 
		WHERE username = $1
	`

	var user domain.User
	err := r.db.QueryRow(query, username).Scan(
		&user.Id,
		&user.Username,
		&user.Name,
		&user.Password,
		&user.TelegramId,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, repository.ErrNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetAll() ([]domain.User, error) {
	query := `
		SELECT id, username, first_name, pass, telegram_id, created_at, updated_at
		FROM users 
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Name,
			&user.Password,
			&user.TelegramId,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *UserRepository) isUsernameTaken(username string, excludeUserId int64) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM users 
		WHERE username = $1 AND id != $2
	`

	var count int
	err := r.db.QueryRow(query, username, excludeUserId).Scan(&count)
	if err != nil {
		return false, repository.ErrUsernameTaken
	}

	return count > 0, nil
}
