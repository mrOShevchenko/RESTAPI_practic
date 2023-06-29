package repository

import (
	"Nix_trainee_practic/internal/models"
	"fmt"
	"github.com/upper/db/v4"
	"strings"
	"time"
)

const UsersTable = "users"

type user struct {
	ID          int64      `db:"id,omitempty"`
	Email       string     `db:"email"`
	Name        string     `db:"name"`
	Password    string     `db:"password,omitempty"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

//go:generate mockery --dir . --name UserRepo --output ./mock
type UserRepo interface {
	Save(user models.User) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindByID(id int64) (models.User, error)
	Delete(id int64) error
}

type userRepo struct {
	coll db.Collection
}

func NewUserRepo(dbSession db.Session) UserRepo {
	return userRepo{
		coll: dbSession.Collection(UsersTable),
	}
}

func (u userRepo) Save(user models.User) (models.User, error) {
	domainUser := u.mapDomainToModel(user)
	domainUser.CreatedDate = time.Now()
	domainUser.UpdatedDate = time.Now()
	err := u.coll.InsertReturning(&domainUser)
	if err != nil {
		return models.User{}, err
	}
	return u.mapModelToDomain(domainUser), nil
}

func (u userRepo) FindByEmail(email string) (models.User, error) {
	var domainUser user
	email = strings.ToLower(email)
	err := u.coll.Find(db.Cond{
		"email":        email,
		"deleted_date": nil,
	}).One(&domainUser)
	if err != nil {
		return models.User{}, fmt.Errorf("user repository save user: %w", err)
	}
	return u.mapModelToDomain(domainUser), nil
}

func (u userRepo) FindByID(id int64) (models.User, error) {
	var domainUser user

	err := u.coll.Find(db.Cond{
		"id":           id,
		"deleted_date": nil,
	}).One(&domainUser)
	if err != nil {
		return models.User{}, fmt.Errorf("user repository finde by id user: %w", err)
	}
	return u.mapModelToDomain(domainUser), nil
}

func (u userRepo) Delete(id int64) error {
	err := u.coll.Find(db.Cond{
		"id":           id,
		"deleted_date": nil,
	}).Update(map[string]interface{}{"deleted_date": time.Now()})
	if err != nil {
		return fmt.Errorf("user repository delete user: %w", err)
	}
	return nil
}

func (u userRepo) mapDomainToModel(d models.User) user {
	return user{
		ID:       d.ID,
		Email:    strings.ToLower(d.Email),
		Password: d.Password,
		Name:     d.Name,
	}
}

func (u userRepo) mapModelToDomain(d user) models.User {
	return models.User{
		ID:          d.ID,
		Email:       d.Email,
		Password:    d.Password,
		Name:        d.Name,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}
