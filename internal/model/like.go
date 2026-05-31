package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Like struct {
	LID       int       `json:"lid"       gorm:"primaryKey;autoIncrement"`
	UID       int       `json:"uid"       gorm:"index;not null"`
	PID       int       `json:"pid"       gorm:"index;not null"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateLike(ctx context.Context, uid, pid int) error {
	// 防重复
	var count int64
	DB.WithContext(ctx).Model(&Like{}).Where("uid = ? AND pid = ?", uid, pid).Count(&count)
	if count > 0 {
		return nil
	}

	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&Like{UID: uid, PID: pid}).Error; err != nil {
			return err
		}
		return IncrLikes(ctx, pid)
	})
}

func DeleteLike(ctx context.Context, uid, pid int) error {
	result := DB.WithContext(ctx).Where("uid = ? AND pid = ?", uid, pid).Delete(&Like{})
	if result.RowsAffected == 0 {
		return nil
	}
	return DecrLikes(ctx, pid)
}
