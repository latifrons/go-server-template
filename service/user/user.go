package user

import (
	"context"
	"github.com/atom-eight/tmt-backend/dbgorm"
	"golang.org/x/crypto/bcrypt"
)

const DefaultPassword = "@t0m8"

type UserService struct {
	DbOperator *dbgorm.DbOperator
}

func (u *UserService) CreateUser(ctx context.Context, email string, phone string) (uint, error) {
	passwordHash, err := HashPassword(DefaultPassword)
	if err != nil {
		return 0, err
	}
	user := &dbgorm.DbUser{
		Email:              email,
		Phone:              phone,
		Password:           passwordHash,
		NeedChangePassword: false,
	}
	err = u.DbOperator.CreateUserT(ctx, user)
	return user.ID, err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
