package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type IntArray []int

func (a *IntArray) Scan(value any) error {
	if value == nil {
		*a = IntArray{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("IntArray: type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}

func (a IntArray) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return json.Marshal(a)
}

type Node struct {
	NID         int       `json:"nid"         gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name"        gorm:"not null"`
	Description string    `json:"description"`
	Article     int       `json:"article"     gorm:"default:0"`
	Moderators  IntArray  `json:"moderators"  gorm:"type:json"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func CreateNode(ctx context.Context, node *Node) error {
	return DB.WithContext(ctx).Create(node).Error
}

func ListNodes(ctx context.Context) ([]Node, error) {
	var nodes []Node
	result := DB.WithContext(ctx).Order("nid asc").Find(&nodes)
	return nodes, result.Error
}

func DeleteNode(ctx context.Context, nid int) error {
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Where("nid = ?", nid).Delete(&Post{})
		return tx.Delete(&Node{}, nid).Error
	})
}

func IncrArticleCount(ctx context.Context, nid int) error {
	return DB.WithContext(ctx).Model(&Node{}).Where("nid = ?", nid).
		Update("article", gorm.Expr("article + 1")).Error
}

func DecrArticleCount(ctx context.Context, nid int) error {
	return DB.WithContext(ctx).Model(&Node{}).Where("nid = ?", nid).
		Update("article", gorm.Expr("CASE WHEN article > 0 THEN article - 1 ELSE 0 END")).Error
}
