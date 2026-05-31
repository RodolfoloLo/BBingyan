package model

import (
	"context"
	"time"

	"BBingyan/internal/controller/param"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Post struct {
	PID       int       `json:"pid"       gorm:"primaryKey;autoIncrement"`
	UID       int       `json:"uid"       gorm:"index;not null"`
	Title     string    `json:"title"     gorm:"not null"`
	NID       int       `json:"nid"       gorm:"index"`
	Likes     int       `json:"likes"     gorm:"default:0"`
	Comments  int       `json:"comments"  gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   *string   `json:"content,omitempty" gorm:"-"` // 不存 DB，查询时从 Body 表填充
}

type Body struct {
	PID     int    `json:"pid"     gorm:"primaryKey"`
	Content string `json:"content" gorm:"not null"`
}

func CreatePost(ctx context.Context, post *Post) (*Post, error) {
	tx := DB.WithContext(ctx).Begin()
	content := post.Content
	post.Content = nil

	result := tx.Clauses(clause.Returning{}).Create(post)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	if content != nil {
		body := &Body{PID: post.PID, Content: *content}
		if err := tx.Create(body).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if post.NID != 0 {
		if err := incrArticleCount(tx, post.NID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return post, tx.Commit().Error
}

func GetPosts(ctx context.Context, p param.Paging) ([]Post, int64, error) {
	var posts []Post
	var total int64

	query := DB.WithContext(ctx).Model(&Post{})
	if p.ID != 0 {
		query = query.Where("nid = ?", p.ID)
	}
	query.Count(&total)

	result := query.
		Limit(p.PageSize).Offset((p.Page - 1) * p.PageSize).
		Order(p.SortClause()).
		Find(&posts)
	return posts, total, result.Error
}

func GetPostByPID(ctx context.Context, pid int) (*Post, error) {
	var post Post
	if err := DB.WithContext(ctx).Where("pid = ?", pid).First(&post).Error; err != nil {
		return nil, err
	}
	var body Body
	if err := DB.WithContext(ctx).Where("pid = ?", pid).First(&body).Error; err == nil {
		post.Content = &body.Content
	}
	return &post, nil
}

func DeletePost(ctx context.Context, pid int) error {
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Where("pid = ?", pid).Delete(&Like{})
		tx.Where("pid = ?", pid).Delete(&Comment{})
		tx.Where("pid = ?", pid).Delete(&Body{})
		return tx.Delete(&Post{}, pid).Error
	})
}

func IncrLikes(ctx context.Context, pid int) error {
	return DB.WithContext(ctx).Model(&Post{}).Where("pid = ?", pid).
		Update("likes", gorm.Expr("likes + 1")).Error
}

func DecrLikes(ctx context.Context, pid int) error {
	return DB.WithContext(ctx).Model(&Post{}).Where("pid = ?", pid).
		Update("likes", gorm.Expr("CASE WHEN likes > 0 THEN likes - 1 ELSE 0 END")).Error
}

func IncrComments(ctx context.Context, pid int) error {
	return DB.WithContext(ctx).Model(&Post{}).Where("pid = ?", pid).
		Update("comments", gorm.Expr("comments + 1")).Error
}

func DecrComments(ctx context.Context, pid int) error {
	return DB.WithContext(ctx).Model(&Post{}).Where("pid = ?", pid).
		Update("comments", gorm.Expr("CASE WHEN comments > 0 THEN comments - 1 ELSE 0 END")).Error
}

func incrArticleCount(tx *gorm.DB, nid int) error {
	return tx.Model(&Node{}).Where("nid = ?", nid).
		Update("article", gorm.Expr("article + 1")).Error
}
