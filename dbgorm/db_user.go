package dbgorm

import (
	"context"
	"errors"
)

func (d *DbOperator) QueryUserByMail(ctx context.Context, mail string) (result DbUser, err error) {
	query := &DbUser{
		Email: mail,
	}
	queryResult := d.db.WithContext(ctx).Where(query).Limit(1).Find(&result)
	err = queryResult.Error
	return
}

func (d *DbOperator) CreateUserT(ctx context.Context, user *DbUser) (err error) {
	d.Begin()
	// check user existence
	var existingUser DbUser
	existingUser, err = d.QueryUserByMail(ctx, user.Email)
	if err != nil {
		return
	}

	if existingUser.ID != 0 {
		return errors.New("record existing")
	}

	result := d.db.WithContext(ctx).Create(user)
	d.Commit()

	return result.Error
}
