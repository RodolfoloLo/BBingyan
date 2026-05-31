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

// GenerateCode 生成 6 位数字验证码
func GenerateCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

// SetCode 将验证码存到 Redis，过期时间由配置决定
func SetCode(ctx context.Context, email, code string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	ttl := time.Duration(config.Conf.Captcha.Expire) * time.Second
	return Redis.Set(ctx, email, code, ttl).Err()
}

// ValidateCode 验证验证码是否正确
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

// CanResend 检查是否可以重新发送验证码
func CanResend(ctx context.Context, email string) (bool, time.Duration, error) {
	ttl, err := Redis.TTL(ctx, email).Result()
	if err != nil {
		return true, 0, err // Redis 出错就放行，别卡死用户
	}
	if ttl <= 0 {
		return true, 0, nil // key 已过期，随便发
	}

	expireSec := time.Duration(config.Conf.Captcha.Expire) * time.Second
	resendSec := time.Duration(config.Conf.Captcha.Resend) * time.Second

	// TTL 还剩很多 → 说明刚发过 → 不让发
	wait := ttl - (expireSec - resendSec)
	if wait > 0 {
		return false, wait, nil
	}
	return true, 0, nil
}
