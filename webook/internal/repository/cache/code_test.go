package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestLuaScript(t *testing.T) {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := c.Ping(context.Background()).Err(); err != nil {
		t.Fatal(err)
	}
	res, err := c.Eval(context.Background(), luaSendCode, []string{"key2"}, "123").Int()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
