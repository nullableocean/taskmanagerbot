package repository

import "taskbot/domain"

type TaskRepository interface {
	Create(task domain.Task) (domain.Task, error)
	Update(task domain.Task) (domain.Task, error)
	GetAll(userId int64) ([]domain.Task, error)
	Get(id int64) (domain.Task, error)
	Delete(task domain.Task) error
}
