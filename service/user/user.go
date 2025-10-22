package user

import (
	"fmt"
	"taskbot/domain"
	"taskbot/pkg/password"
	"taskbot/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (s *UserService) Save(u domain.User) (domain.User, error) {
	err := s.validateNewUser(u)
	if err != nil {
		return u, err
	}

	hashPass, _ := password.HashPassword(u.Password)
	u.Password = hashPass

	return s.repo.Save(u)
}

func (s *UserService) validateNewUser(data domain.User) error {
	if data.Username == "" {
		return fmt.Errorf("%w. empty username", ErrValidateData)
	}
	if data.Password == "" {
		return fmt.Errorf("%w. empty password", ErrValidateData)
	}

	return nil
}
