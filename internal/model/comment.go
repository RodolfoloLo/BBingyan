package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	CID       int       `json:"cid"       gorm:"primaryKey;autoIncrement"`
	UID       int       `json:"uid"       gorm:"index;not null"`
	PID       int       `json:"pid"       gorm:"index;not null"`
	Content   string    `json:"content"   gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateComment(ctx context.Context, comment *Comment) error {
	tx := DB.WithContext(ctx).Begin()

	if err := tx.Create(comment).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&Post{}).Where("pid = ?", comment.PID).
		Update("comments", gorm.Expr("comments + 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func ListCommentsByPID(ctx context.Context, pid int) ([]Comment, error) {
	var comments []Comment
	result := DB.WithContext(ctx).Where("pid = ?", pid).Order("created_at asc").Find(&comments)
	return comments, result.Error
}

func DeleteComment(ctx context.Context, cid int) error {
	var comment Comment
	if err := DB.WithContext(ctx).Where("cid = ?", cid).First(&comment).Error; err != nil {
		return err
	}

	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Comment{}, cid).Error; err != nil {
			return err
		}
		return tx.Model(&Post{}).Where("pid = ?", comment.PID).
			Update("comments", gorm.Expr("CASE WHEN comments > 0 THEN comments - 1 ELSE 0 END")).Error
	})
}
