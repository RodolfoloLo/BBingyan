package utils

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"fmt"

	"BBingyan/internal/config"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.Conf.Redis.Host,
		Password: config.Conf.Redis.Password,
		DB:       config.Conf.Redis.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, err := Redis.Ping(ctx).Result(); err != nil {
		panic("redis connect: " + err.Error())
	}
}

func GenerateCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

func SetCode(ctx context.Context, email, code string) error {
	ttl := time.Duration(config.Conf.Captcha.Expire) * time.Second
	return Redis.Set(ctx, email, code, ttl).Err()
}

func ValidateCode(ctx context.Context, email, code string) (bool, error) {
	stored, err := Redis.Get(ctx, email).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return stored == code, nil
}

func CanResend(ctx context.Context, email string) (bool, time.Duration, error) {
	ttl, err := Redis.TTL(ctx, email).Result()
	if err != nil {
		return true, 0, err
	}
	if ttl <= 0 {
		return true, 0, nil
	}

	expireSec := time.Duration(config.Conf.Captcha.Expire) * time.Second
	resendSec := time.Duration(config.Conf.Captcha.Resend) * time.Second
	elapsed := expireSec - ttl
	if elapsed < resendSec {
		return false, resendSec - elapsed, nil
	}
	return true, 0, nil
}
