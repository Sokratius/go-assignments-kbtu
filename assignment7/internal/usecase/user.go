package usecase

import (
	"fmt"
	"practice-7/internal/entity"
	"practice-7/internal/usecase/repo"
	"practice-7/utils"

	"github.com/google/uuid"
)

type UserUseCase struct {
	repo      *repo.UserRepo
	jwtSecret string
}

func NewUserUseCase(r *repo.UserRepo, jwtSecret string) *UserUseCase {
	return &UserUseCase{repo: r, jwtSecret: jwtSecret}
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	if user.Role == "" {
		user.Role = "user"
	}

	storedUser, err := u.repo.RegisterUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("register user: %w", err)
	}

	sessionID := uuid.NewString()
	return storedUser, sessionID, nil
}

func (u *UserUseCase) LoginUser(user *entity.LoginUserDTO) (string, error) {
	userFromRepo, err := u.repo.GetByUsername(user.Username)
	if err != nil {
		return "", fmt.Errorf("get user by username: %w", err)
	}

	if !utils.CheckPassword(userFromRepo.Password, user.Password) {
		return "", fmt.Errorf("invalid username or password")
	}

	token, err := utils.GenerateJWT(userFromRepo.ID, userFromRepo.Role, []byte(u.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("generate jwt: %w", err)
	}

	return token, nil
}

func (u *UserUseCase) GetMe(userID uuid.UUID) (*entity.User, error) {
	user, err := u.repo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("get me: %w", err)
	}
	return user, nil
}

func (u *UserUseCase) PromoteUser(targetID uuid.UUID, role string) (*entity.User, error) {
	if role == "" {
		return nil, fmt.Errorf("role is required")
	}

	updatedUser, err := u.repo.UpdateRole(targetID, role)
	if err != nil {
		return nil, fmt.Errorf("promote user: %w", err)
	}

	return updatedUser, nil
}
