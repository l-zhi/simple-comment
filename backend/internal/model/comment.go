package model

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论表结构（GORM 模型，支持软删除与索引）
type Comment struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `gorm:"index:idx_article_created,priority:2" json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ArticleID uint   `gorm:"index:idx_article_created,priority:1;index:idx_article_parent,priority:1;not null" json:"articleId"`
	UserID    uint   `gorm:"index;not null" json:"userId"`
	UserName  string `gorm:"size:50;not null;default:''" json:"userName"`
	Avatar    string `gorm:"size:255;not null;default:''" json:"avatar"`

	ParentID    uint  `gorm:"index:idx_article_parent,priority:2;index:idx_parent_id;default:0" json:"parentId"` // 0=根评论；回复时=被回复评论id（可为根或回复）
	ReplyRootID int64 `gorm:"index:idx_reply_root_id;default:0" json:"replyRootId"`                               // 0=根评论；回复时=所属根评论id

	ReplyToUserName string `gorm:"-" json:"replyToUserName,omitempty"` // 展示用：由 parent_id 查被回复记录
	ReplyToContent  string `gorm:"-" json:"replyToContent,omitempty"`  // 展示用：被回复评论内容摘要

	Content string `gorm:"type:text;not null" json:"content"`
	Status  int8   `gorm:"type:tinyint;default:1" json:"status"`
	Likes   int    `gorm:"default:0" json:"likes"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}

// CommentTree 树形展示用，带子回复
type CommentTree struct {
	Comment
	Replies []*CommentTree `json:"replies,omitempty" gorm:"-"`
}

// RootWithPreview 一级评论 + 回复总数 + 预览的 2 条二级回复（API 返回用）
type RootWithPreview struct {
	Comment    Comment   `json:"comment"`
	ReplyCount int64     `json:"replyCount"`
	Replies    []Comment `json:"replies"`
}

const (
	StatusPending = 0 // 审核中
	StatusNormal  = 1 // 正常可见
	StatusHidden  = 2 // 已折叠/隐藏
)
