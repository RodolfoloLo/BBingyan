package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Follow struct {
	UID       int       `json:"uid"       gorm:"primaryKey;index"`
	Followee  int       `json:"followee"  gorm:"primaryKey;index"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateFollow(ctx context.Context, uid, followee int) error {
	var count int64
	DB.WithContext(ctx).Model(&Follow{}).Where("uid = ? AND followee = ?", uid, followee).Count(&count)
	if count > 0 {
		return nil
	}

	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&Follow{UID: uid, Followee: followee}).Error; err != nil {
			return err
		}
		if err := IncrFollowed(ctx, uid); err != nil {
			return err
		}
		return IncrFollowers(ctx, followee)
	})
}

func DeleteFollow(ctx context.Context, uid, followee int) error {
	result := DB.WithContext(ctx).Where("uid = ? AND followee = ?", uid, followee).Delete(&Follow{})
	if result.RowsAffected == 0 {
		return nil
	}
	if err := DecrFollowed(ctx, uid); err != nil {
		return err
	}
	return DecrFollowers(ctx, followee)
}
