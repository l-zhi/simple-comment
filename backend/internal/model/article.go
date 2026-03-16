package model

import (
	"time"

	"gorm.io/gorm"
)

// Article 帖子表
type Article struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Title   string `gorm:"size:200;not null;default:''" json:"title"`
	Content string `gorm:"type:text;not null" json:"content"`
}

func (Article) TableName() string {
	return "articles"
}
