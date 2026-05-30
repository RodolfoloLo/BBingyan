package config

import (
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
	viper.SetConfigFile("config/config.yaml") // 告诉 Viper 去哪找文件
	if err := viper.ReadInConfig(); err != nil {
		panic("读取配置文件失败: " + err.Error())
	}

	Conf = &Config{}
	if err := viper.Unmarshal(Conf); err != nil {
		panic("解析配置到结构体失败: " + err.Error())
	}
}
