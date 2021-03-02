package tools

import (
	"context"
	"time"
)

func GetContext(second int) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(second))
	return ctx
}

func GetContextDefault() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	return ctx
}
