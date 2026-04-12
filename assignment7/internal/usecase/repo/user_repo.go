package repo

import (
	"practice-7/internal/entity"
	"practice-7/pkg/postgres"

	"github.com/google/uuid"
)

type UserRepo struct {
	PG *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) (*UserRepo, error) {
	r := &UserRepo{PG: pg}
	if err := r.PG.Conn.AutoMigrate(&entity.User{}); err != nil {
		return nil, err
	}
	return r, nil
}

func (u *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	if err := u.PG.Conn.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) GetByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := u.PG.Conn.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) GetByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	if err := u.PG.Conn.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) UpdateRole(id uuid.UUID, role string) (*entity.User, error) {
	if err := u.PG.Conn.Model(&entity.User{}).Where("id = ?", id).Update("role", role).Error; err != nil {
		return nil, err
	}
	return u.GetByID(id)
}
