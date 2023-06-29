package service

import "golang.org/x/crypto/bcrypt"

//go:generate mockery --dir . --name Generator --output ./mocks
type Generator interface {
	GeneratePasswordHash(password string) (string, error)
}

type generatePasswordHash struct {
	cost int
}

func NewGeneratePasswordHash(cost int) Generator {
	return generatePasswordHash{
		cost: cost,
	}
}

func (g generatePasswordHash) GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), g.cost)
	return string(bytes), err
}
