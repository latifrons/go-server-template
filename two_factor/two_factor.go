package two_factor

import (
	"context"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"strconv"
	"time"
)

const CodePrefix = "code_"

type TwoFactorValidator struct {
	Address               string
	Password              string
	DBId                  int
	CodeExpirationSeconds int

	rdb *redis.Client
}

func (t *TwoFactorValidator) InitDefault() {
	t.rdb = redis.NewClient(&redis.Options{
		Addr:     t.Address,
		Password: t.Password,
		DB:       t.DBId,
	})
}

func (t *TwoFactorValidator) SetCode(ctx context.Context, userKey string, code string, expiration time.Duration) (bool, error) {
	result := t.rdb.SetNX(ctx, CodePrefix+userKey, code, expiration)
	return result.Result()
}

func (t *TwoFactorValidator) ValidateCode(ctx context.Context, userKey string, code string) (bool, error) {
	cmd := t.rdb.Get(ctx, userKey)
	v, err := cmd.Result()
	if err != nil {
		return false, err
	}
	return v == code, nil
}

func (t *TwoFactorValidator) GenerateNumberCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(899999) + 100000
	res := strconv.Itoa(code)
	return res
}
