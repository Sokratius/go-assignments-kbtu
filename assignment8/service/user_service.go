package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MBFG9000/golang-practice-8/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(id int) (*repository.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) CreateUser(user *repository.User) error {
	return s.repo.CreateUser(user)
}

func (s *UserService) RegisterUser(user *repository.User, email string) error {
	existingUser, err := s.repo.GetByEmail(email)
	if err != nil {
		return fmt.Errorf("error getting user with this email: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("user with this email already exists")
	}
	if err := s.repo.CreateUser(user); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (s *UserService) UpdateUserName(id int, newName string) error {
	if strings.TrimSpace(newName) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user with id %d not found", id)
	}
	user.Name = newName
	if err := s.repo.UpdateUser(user); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (s *UserService) DeleteUser(id int) error {
	if id == 1 {
		return errors.New("it is not allowed to delete admin user")
	}
	if err := s.repo.DeleteUser(id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
