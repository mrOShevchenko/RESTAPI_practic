package service

import (
	"Nix_trainee_practic/internal/models"
	"Nix_trainee_practic/internal/repository"
	"fmt"
	"log"
)

//go:generate mockery --dir . --name UserService --output ./mocks
type UserService interface {
	Save(user models.User) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindByID(id int64) (models.User, error)
	Delete(id int64) error
}

type userService struct {
	userRepo    repository.UserRepo
	passwordGen Generator
}

func NewUser(ur repository.UserRepo, gs Generator) UserService {
	return userService{
		userRepo:    ur,
		passwordGen: gs,
	}
}

func (u userService) Save(user models.User) (models.User, error) {
	var err error

	user.Password, err = u.passwordGen.GeneratePasswordHash(user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("user service save user, could not generate hash: %w", err)
	}

	saveUser, err := u.userRepo.Save(user)
	if err != nil {
		return models.User{}, fmt.Errorf("user service save user: %w", err)
	}
	return saveUser, nil
}

func (u userService) FindByEmail(email string) (models.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		log.Println(err)
		return models.User{}, fmt.Errorf("user service find by email user: %w", err)
	}
	return user, nil
}

func (u userService) FindByID(id int64) (models.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		log.Println(err)
		return models.User{}, fmt.Errorf("user service find by id user: %w", err)
	}
	return user, nil
}

func (u userService) Delete(id int64) error {
	err := u.userRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("user service delete user: %w", err)
	}
	return nil
}
