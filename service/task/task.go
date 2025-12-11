package task

import (
	"errors"
	"fmt"
	"taskbot/domain"
	"taskbot/repository"
	"time"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) Create(user domain.User, task domain.Task) (domain.Task, error) {
	task.UserId = user.Id
	task.Status = domain.WAITING

	err := s.validateTask(task)
	if err != nil {
		return domain.Task{}, err
	}

	return s.repo.Create(task)
}

func (s *TaskService) Delete(task domain.Task) error {
	return s.repo.Delete(task)
}

func (s *TaskService) Update(task domain.Task) (domain.Task, error) {
	err := s.validateTask(task)
	if err != nil {
		return domain.Task{}, err
	}

	task.UpdatedAt = time.Now()
	return s.repo.Update(task)
}

func (s *TaskService) GetAll(user domain.User) ([]domain.Task, error) {
	if user.Id == 0 {
		return []domain.Task{}, errors.New("get tasks error: empty user id")
	}

	return s.repo.GetAll(user.Id)
}

func (s *TaskService) GetById(taskId int64) (domain.Task, error) {
	return s.repo.Get(taskId)
}

func (s *TaskService) GetAllByStatus(user domain.User, taskStatus domain.TaskStatus) ([]domain.Task, error) {
	tasks, err := s.GetAll(user)
	if err != nil {
		return []domain.Task{}, err
	}

	filteredTasks := make([]domain.Task, 0, len(tasks))
	for _, t := range tasks {
		if t.Status == taskStatus {
			filteredTasks = append(filteredTasks, t)
		}
	}

	return filteredTasks, nil
}

func (s *TaskService) validateTask(task domain.Task) error {
	if task.UserId == 0 {
		return fmt.Errorf("%w: task userId is empty")
	}

	if task.Title == "" {
		return fmt.Errorf("%w: task title is empty")
	}

	if task.Status == "" {
		return fmt.Errorf("%w: task status is empty")
	}

	return nil
}
