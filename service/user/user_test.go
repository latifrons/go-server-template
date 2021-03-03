package user

import (
	"context"
	"github.com/atom-eight/tmt-backend/dbgorm"
	"github.com/atom-eight/tmt-backend/tools"
	"testing"
)

func GetDbOperatorTestInstance() *dbgorm.DbOperator {
	dbOperator := &dbgorm.DbOperator{
		//Source: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		//	viper.GetString("mysql.username"), viper.GetString("mysql.password"),
		//	viper.GetString("mysql.url"), viper.GetString("mysql.schema")),
		Source: "test.db",
	}
	dbOperator.InitDefault()
	return dbOperator
}

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "t1",
			args: args{
				"this is a password",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HashPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	type fields struct {
		DbOperator *dbgorm.DbOperator
	}
	type args struct {
		ctx   context.Context
		email string
		phone string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				DbOperator: GetDbOperatorTestInstance(),
			},
			args: args{
				ctx:   tools.GetContextDefault(),
				email: "latifrons88@gmail.com",
				phone: "13918280266",
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserService{
				DbOperator: tt.fields.DbOperator,
			}
			got, err := u.CreateUser(tt.args.ctx, tt.args.email, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
