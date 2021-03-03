package two_factor

import (
	"context"
	"github.com/atom-eight/tmt-backend/tools"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

const RedisAddress = "gitez.cc:7983"
const RedisPassword = "atom1234"

func TestTwoFactorValidator_GenerateNumberCode(t1 *testing.T) {
	type fields struct {
		Address               string
		Password              string
		DBId                  int
		CodeExpirationSeconds int
		rdb                   *redis.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "t1",
			fields: fields{
				Address:               RedisAddress,
				Password:              RedisPassword,
				DBId:                  0,
				CodeExpirationSeconds: 30,
				rdb:                   nil,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TwoFactorValidator{
				Address:               tt.fields.Address,
				Password:              tt.fields.Password,
				DBId:                  tt.fields.DBId,
				CodeExpirationSeconds: tt.fields.CodeExpirationSeconds,
				rdb:                   tt.fields.rdb,
			}
			t.InitDefault()
			if got := t.GenerateNumberCode(); len(got) != 6 {
				t1.Errorf("GenerateNumberCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

//
func TestTwoFactorValidator_SetCode(t1 *testing.T) {
	type fields struct {
		Address               string
		Password              string
		DBId                  int
		CodeExpirationSeconds int
		rdb                   *redis.Client
	}
	type args struct {
		ctx        context.Context
		userKey    string
		code       string
		expiration time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "t1",
			fields: fields{
				Address:               RedisAddress,
				Password:              RedisPassword,
				DBId:                  0,
				CodeExpirationSeconds: 30,
				rdb:                   nil,
			},
			args: args{},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TwoFactorValidator{
				Address:               tt.fields.Address,
				Password:              tt.fields.Password,
				DBId:                  tt.fields.DBId,
				CodeExpirationSeconds: tt.fields.CodeExpirationSeconds,
				rdb:                   tt.fields.rdb,
			}
			t.InitDefault()
			b, err := t.SetCode(tools.GetContextDefault(), "usermail", "123412", time.Second*10)
			if err != nil || !b {
				t1.Errorf("TwoFactorValidator_SetCode()")
			}
		})
	}
}

func TestTwoFactorValidator_ValidateCode(t1 *testing.T) {
	type fields struct {
		Address               string
		Password              string
		DBId                  int
		CodeExpirationSeconds int
		rdb                   *redis.Client
	}
	type args struct {
		ctx     context.Context
		userKey string
		code    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TwoFactorValidator{
				Address:               tt.fields.Address,
				Password:              tt.fields.Password,
				DBId:                  tt.fields.DBId,
				CodeExpirationSeconds: tt.fields.CodeExpirationSeconds,
				rdb:                   tt.fields.rdb,
			}
			got, err := t.ValidateCode(tt.args.ctx, tt.args.userKey, tt.args.code)
			if (err != nil) != tt.wantErr {
				t1.Errorf("ValidateCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t1.Errorf("ValidateCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
