package service

import (
	"Nix_trainee_practic/internal/models"
	"Nix_trainee_practic/internal/repository"
	"Nix_trainee_practic/mocks"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserService_FindByEmail(t *testing.T) {
	testTable := []struct {
		name            string
		email           string
		repoConstructor func() repository.UserRepo
		expect          models.User
		expectErr       bool
	}{
		{

			"OK FindMyEmail",
			"test@mail.com",
			func() repository.UserRepo {
				mock := mocks.NewUserRepo(t)
				mock.
					On("FindByEmail", "test@mail.com").
					Return(models.User{Email: "test@mail.com"}, nil).Times(1)
				return mock
			},
			models.User{
				Email: "test@mail.com",
			},
			false,
		},
		{
			"ERROR FindMyEmail",
			"test@mail.com",
			func() repository.UserRepo {
				mock := mocks.NewUserRepo(t)
				mock.
					On("FindByEmail", "test@mail.com").
					Return(models.User{}, errors.New("testError")).Times(1)
				return mock
			},
			models.User{},
			true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			us := userService{
				userRepo: tt.repoConstructor(),
			}
			user, err := us.FindByEmail(tt.email)
			fmt.Println(user)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, user, tt.expect)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, user, tt.expect)
			}
		})
	}
}

func TestUserService_FindByID(t *testing.T) {
	testTable := []struct {
		name            string
		id              int64
		repoConstructor func() repository.UserRepo
		expect          models.User
		expectErr       bool
	}{
		{
			"OK FindByID",
			2,
			func() repository.UserRepo {
				mock := mocks.NewUserRepo(t)
				mock.
					On("FindByID", int64(2)).
					Return(models.User{ID: 2}, nil).Times(1)
				return mock
			},
			models.User{
				ID: 2,
			},
			false,
		},
		{
			"ERROR FindByID",
			2,
			func() repository.UserRepo {
				mock := mocks.NewUserRepo(t)
				mock.
					On("FindByID", int64(2)).
					Return(models.User{}, errors.New("testError")).Times(1)
				return mock
			},
			models.User{},
			true,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			us := userService{
				userRepo: tt.repoConstructor(),
			}
			user, err := us.FindByID(tt.id)
			fmt.Println(user)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, user, tt.expect)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, user, tt.expect)
			}
		})
	}
}

func TestUserService_Save(t *testing.T) {
	tests := []struct {
		name                 string
		user                 models.User
		repoConstructor      func(user models.User) repository.UserRepo
		generatorConstructor func(password string) Generator
		expect               models.User
		expectErr            bool
	}{
		{
			"OK userSave",
			models.User{
				Email:    "test@mail.com",
				Name:     "user",
				Password: "11111",
			},
			func(user models.User) repository.UserRepo {
				user.Password = "password"
				mock := mocks.NewUserRepo(t)
				mock.
					On("Save", user).
					Return(models.User{
						Email:    user.Email,
						Name:     user.Name,
						Password: user.Password,
					}, nil).Times(1)
				return mock
			},
			func(password string) Generator {
				mock := NewMockGenerator(t)
				mock.
					On("GeneratePasswordHash", password).
					Return("password", nil).Times(1)
				return mock
			},
			models.User{
				Email:    "test@mail.com",
				Name:     "user",
				Password: "password",
			},
			false,
		},
		{
			"ERROR userSave",
			models.User{
				Email:    "test@mail.com",
				Name:     "user",
				Password: "1111",
			},
			func(user models.User) repository.UserRepo {
				mock := mocks.NewUserRepo(t)
				user.Password = "password"
				mock.
					On("Save", user).
					Return(models.User{}, errors.New("testError")).Times(1)
				return mock
			},
			func(password string) Generator {
				mock := NewMockGenerator(t)
				mock.
					On("GeneratePasswordHash", password).
					Return("password", nil).Times(1)
				return mock
			},
			models.User{},
			true,
		},
		{
			"ERROR generator in userSave",
			models.User{
				Email:    "test@mail.com",
				Name:     "user",
				Password: "password",
			},
			func(user models.User) repository.UserRepo {
				mock := mocks.NewUserRepo(t)
				return mock
			},
			func(password string) Generator {
				mock := NewMockGenerator(t)
				mock.
					On("GeneratePasswordHash", password).
					Return(password, errors.New("Error with generator")).Times(1)
				return mock
			},
			models.User{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := userService{
				passwordGen: tt.generatorConstructor(tt.user.Password),
				userRepo:    tt.repoConstructor(tt.user),
			}
			user, err := NewUser(us.userRepo, us.passwordGen).Save(tt.user)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, user, tt.expect)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, user, tt.expect)
			}
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	testTable := []struct {
		name            string
		id              int64
		repoConstructor func(id int64) repository.UserRepo
		expect          error
		expectErr       bool
	}{
		{
			"OK userDelete",
			2,
			func(id int64) repository.UserRepo {
				mock := mocks.NewUserRepo(t)
				mock.
					On("Delete", id).
					Return(nil).Times(1)
				return mock
			},
			nil,
			false,
		},
		{
			"ERROR userDelete",
			2,
			func(id int64) repository.UserRepo {
				mock := mocks.NewUserRepo(t)
				mock.
					On("Delete", id).
					Return(errors.New("error: can't delete")).Times(1)
				return mock
			},
			errors.New("error: can't delete"),
			true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			u := userService{
				userRepo: tt.repoConstructor(tt.id),
			}
			err := u.Delete(tt.id)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
