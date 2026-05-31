package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Jwt     JwtConfig      `mapstructure:"jwt"`
	Server  ServerConfig   `mapstructure:"server"`
	DB      PostgresConfig `mapstructure:"postgres"`
	Admin   AdminConfig    `mapstructure:"admin"`
	Mail    MailConfig     `mapstructure:"mail"`
	Captcha CaptchaConfig  `mapstructure:"captcha"`
	Redis   RedisConfig    `mapstructure:"redis"`
	Logger  LoggerConfig   `mapstructure:"logging"`
}

type JwtConfig struct {
	Secret       string   `mapstructure:"secret"`
	Expire       int64    `mapstructure:"expire"`
	SkippedPaths []string `mapstructure:"skippedpaths"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Ver  string `mapstructure:"ver"`
}

type PostgresConfig struct {
	Dsn string `mapstructure:"dsn"`
}

type AdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type MailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Nickname string `mapstructure:"nickname"`
}

type CaptchaConfig struct {
	Expire int    `mapstructure:"expire"`
	Resend int    `mapstructure:"resend"`
	Title  string `mapstructure:"title"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LoggerConfig struct {
	Debug bool   `mapstructure:"debug"`
	Path  string `mapstructure:"path"`
}

var Conf *Config

func InitConfig() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config/config.yaml"
	}

	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic("read config failed: " + err.Error())
	}

	Conf = &Config{}
	if err := viper.Unmarshal(Conf); err != nil {
		panic("unmarshal config failed: " + err.Error())
	}

	if err := validate(); err != nil {
		panic("config invalid: " + err.Error())
	}
}

// 新增的配置验证函数，确保关键配置项被正确设置
func validate() error {
	if Conf.Jwt.Secret == "" || Conf.Jwt.Secret == "change-me-to-a-random-string" {
		return errors.New("jwt.secret 未设置或仍为默认值")
	}
	if Conf.DB.Dsn == "" {
		return errors.New("postgres.dsn 未设置")
	}
	if Conf.Mail.Host == "" {
		return errors.New("mail.host 未设置")
	}
	if Conf.Redis.Host == "" {
		return errors.New("redis.host 未设置")
	}
	return nil
}
