package model

import (
	"time"

	"gorm.io/gorm"

	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"BBingyan/internal/config"
)

type User struct {
	ID         int            `json:"id"         gorm:"primaryKey;autoIncrement"`
	Username   string         `json:"username"   gorm:"uniqueIndex;not null"`
	Password   string         `json:"-"          gorm:"not null"` // json:"-" = 绝不序列化
	Email      string         `json:"email"      gorm:"not null"`
	Nickname   string         `json:"nickname"`
	Permission int            `json:"permission" gorm:"default:0"` // 0=user, 1=admin
	Followed   int            `json:"followed"   gorm:"default:0"`
	Followers  int            `json:"followers"  gorm:"default:0"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-"          gorm:"index"` // 软删除
}

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exists")
)

func InitDefaultAdmin() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := GetUserByUsername(ctx, "admin")
	if err == nil {
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(config.Conf.Admin.Password), bcrypt.DefaultCost)
	admin := &User{
		Username:   config.Conf.Admin.Username,
		Password:   string(hash),
		Permission: 1,
	}
	if err := DB.WithContext(ctx).Create(admin).Error; err != nil {
		panic("create default admin: " + err.Error())
	}
}

func CreateUser(ctx context.Context, user *User) error {
	_, err := GetUserByUsername(ctx, user.Username)
	if err == nil {
		return ErrUserAlreadyExist
	}
	if !errors.Is(err, ErrUserNotFound) {
		return err
	}
	return DB.WithContext(ctx).Create(user).Error
}

func GetUserByID(ctx context.Context, id int) (*User, error) {
	var user User
	result := DB.WithContext(ctx).Where("id = ?", id).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, result.Error
}

func GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	result := DB.WithContext(ctx).Where("username = ?", username).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, result.Error
}

func DeleteUser(ctx context.Context, id int) error {
	result := DB.WithContext(ctx).Delete(&User{}, id)
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return result.Error
}

func IncrFollowed(ctx context.Context, id int) error {
	return DB.WithContext(ctx).Model(&User{}).Where("id = ?", id).
		UpdateColumn("followed", gorm.Expr("followed + 1")).Error
}

func DecrFollowed(ctx context.Context, id int) error {
	return DB.WithContext(ctx).Model(&User{}).Where("id = ?", id).
		UpdateColumn("followed", gorm.Expr("CASE WHEN followed > 0 THEN followed - 1 ELSE 0 END")).Error
}

func IncrFollowers(ctx context.Context, id int) error {
	return DB.WithContext(ctx).Model(&User{}).Where("id = ?", id).
		UpdateColumn("followers", gorm.Expr("followers + 1")).Error
}

func DecrFollowers(ctx context.Context, id int) error {
	return DB.WithContext(ctx).Model(&User{}).Where("id = ?", id).
		UpdateColumn("followers", gorm.Expr("CASE WHEN followers > 0 THEN followers - 1 ELSE 0 END")).Error
}
