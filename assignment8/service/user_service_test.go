package service

import (
	"fmt"
	"testing"

	"github.com/MBFG9000/golang-practice-8/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Legendov"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := userService.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)

}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)
	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)
	err := userService.CreateUser(user)
	assert.NoError(t, err)
}

func TestRegisterUser(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRepo := repository.NewMockUserRepository(mockCtrl)
		userService := NewUserService(mockRepo)

		user := &repository.User{ID: 1, Name: "Testobekov Testovich", Email: "testobekov@test.ru"}
		mockRepo.EXPECT().GetByEmail(user.Email).Return(nil, nil)
		mockRepo.EXPECT().CreateUser(user).Return(nil)
		err := userService.RegisterUser(user, user.Email)

		assert.NoError(t, err)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRepo := repository.NewMockUserRepository(mockCtrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 1, Name: "Testobekov Testovich", Email: "testobekov@test.ru"}

		mockRepo.EXPECT().GetByEmail(user.Email).Return(user, nil)
		err := userService.RegisterUser(user, user.Email)

		assert.EqualError(t, err, "user with this email already exists")
	})

	t.Run("repository error on GetByEmail", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRepo := repository.NewMockUserRepository(mockCtrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 1, Name: "Testobekov Testovich", Email: "testobekov@test.ru"}

		mockRepo.EXPECT().GetByEmail(user.Email).Return(nil, fmt.Errorf("db is unavailable"))
		err := userService.RegisterUser(user, user.Email)

		assert.ErrorContains(t, err, "error getting user with this email")
	})

	t.Run("repository error on CreateUser", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRepo := repository.NewMockUserRepository(mockCtrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 1, Name: "Testobekov Testovich", Email: "testobekov@test.ru"}

		mockRepo.EXPECT().GetByEmail(user.Email).Return(nil, nil)
		mockRepo.EXPECT().CreateUser(user).Return(fmt.Errorf("insert failed"))

		err := userService.RegisterUser(user, user.Email)

		assert.ErrorContains(t, err, "error creating user")
	})
}

func TestUpdateUserName(t *testing.T) {
	t.Run("empty name", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRepo := repository.NewMockUserRepository(mockCtrl)
		userService := NewUserService(mockRepo)

		err := userService.UpdateUserName(1, "")
		assert.EqualError(t, err, "name cannot be empty")
	})

	t.Run("user not found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRepo := repository.NewMockUserRepository(mockCtrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 1, Name: "Testobekov Testovich", Email: "testobekov@test.ru"}

		mockRepo.EXPECT().GetUserByID(user.ID).Return(nil, fmt.Errorf("user with id %d not found", user.ID))
		err := userService.UpdateUserName(user.ID, "NewName")
		assert.EqualError(t, err, fmt.Sprintf("user with id %d not found", user.ID))
	})

	t.Run("repository error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRepo := repository.NewMockUserRepository(mockCtrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 1, Name: "Testobekov Testovich", Email: "testobekov@test.ru"}

		mockRepo.EXPECT().GetUserByID(user.ID).Return(nil, fmt.Errorf("repository error"))
		err := userService.UpdateUserName(user.ID, "NewName")
		assert.Error(t, err)
	})

	t.Run("successful update", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRepo := repository.NewMockUserRepository(mockCtrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 1, Name: "Testobekov Testovich", Email: "testobekov@test.ru"}

		initial := user
		updated := initial
		updated.Name = "NewName"
		mockRepo.EXPECT().GetUserByID(initial.ID).Return(initial, nil)
		mockRepo.EXPECT().UpdateUser(initial).Return(nil)
		mockRepo.EXPECT().GetUserByID(initial.ID).Return(updated, nil)

		err := userService.UpdateUserName(initial.ID, updated.Name)

		assert.NoError(t, err)

		updatedUserService, err := userService.GetUserByID(initial.ID)

		assert.Equal(t, updated.Name, updatedUserService.Name)
		assert.NoError(t, err)
	})

	t.Run("update user fail", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRepo := repository.NewMockUserRepository(mockCtrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 1, Name: "Testobekov Testovich", Email: "testobekov@test.ru"}

		mockRepo.EXPECT().GetUserByID(user.ID).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(user).Return(fmt.Errorf("repository error"))

		err := userService.UpdateUserName(user.ID, "NewName")
		assert.ErrorContains(t, err, "error updating user")
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("admin delete", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRepo := repository.NewMockUserRepository(mockCtrl)
		userService := NewUserService(mockRepo)

		err := userService.DeleteUser(1)
		assert.EqualError(t, err, "it is not allowed to delete admin user")
	})

	t.Run("successful deletion", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRepo := repository.NewMockUserRepository(mockCtrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 2, Name: "Testobekov Testovich", Email: "testobekov@test.ru"}

		mockRepo.EXPECT().DeleteUser(user.ID).Return(nil)
		mockRepo.EXPECT().GetUserByID(user.ID).Return(nil, fmt.Errorf("user with id %d not found", user.ID))

		err := userService.DeleteUser(user.ID)

		assert.Nil(t, err)

		userGetByService, err := userService.GetUserByID(user.ID)

		assert.Nil(t, userGetByService)
		assert.EqualError(t, err, fmt.Sprintf("user with id %d not found", user.ID))
	})

	t.Run("repository error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockRepo := repository.NewMockUserRepository(mockCtrl)
		userService := NewUserService(mockRepo)
		user := &repository.User{ID: 2, Name: "Testobekov Testovich", Email: "testobekov@test.ru"}

		mockRepo.EXPECT().DeleteUser(user.ID).Return(fmt.Errorf("repository error"))
		err := userService.DeleteUser(user.ID)
		assert.ErrorContains(t, err, "error deleting user")
	})
}
