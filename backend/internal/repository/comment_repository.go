package repository

import (
	"context"

	"gorm.io/gorm"

	"simple-comment/internal/model"
)

type CommentRepository interface {
	Create(ctx context.Context, c *model.Comment) error
	DeleteByID(ctx context.Context, id int64) error
	GetCommentByID(ctx context.Context, id int64) (*model.Comment, error)
	GetCommentsByIDs(ctx context.Context, ids []int64) map[int64]*model.Comment
	ListRootsByArticle(ctx context.Context, articleID uint, offset, limit int) ([]*model.Comment, int64, error)
	ListRepliesByParent(ctx context.Context, parentID int64, offset, limit int) ([]*model.Comment, int64, error)
	CountRepliesByParents(ctx context.Context, parentIDs []int64) map[int64]int64
	GetPreviewRepliesForRoots(ctx context.Context, rootIDs []int64, per int) map[int64][]*model.Comment
}

type commentRepo struct {
	db *gorm.DB
}

// NewCommentRepository 返回 GORM 实现的 CommentRepository
func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepo{db: db}
}

func (r *commentRepo) Create(ctx context.Context, c *model.Comment) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *commentRepo) DeleteByID(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Comment{}).Error
}

func (r *commentRepo) GetCommentByID(ctx context.Context, id int64) (*model.Comment, error) {
	var c model.Comment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *commentRepo) GetCommentsByIDs(ctx context.Context, ids []int64) map[int64]*model.Comment {
	if len(ids) == 0 {
		return nil
	}
	var list []*model.Comment
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&list).Error; err != nil {
		return nil
	}
	out := make(map[int64]*model.Comment)
	for _, c := range list {
		out[c.ID] = c
	}
	return out
}

func (r *commentRepo) ListRootsByArticle(ctx context.Context, articleID uint, offset, limit int) ([]*model.Comment, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("article_id = ? AND parent_id = 0", articleID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []*model.Comment
	err := r.db.WithContext(ctx).
		Where("article_id = ? AND parent_id = 0", articleID).
		Order("created_at ASC").
		Offset(offset).Limit(limit).
		Find(&list).Error
	return list, total, err
}

// ListRepliesByParent 返回某个根评论楼内的回复分页（按 reply_root_id）
func (r *commentRepo) ListRepliesByParent(ctx context.Context, rootID int64, offset, limit int) ([]*model.Comment, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("reply_root_id = ?", rootID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []*model.Comment
	err := r.db.WithContext(ctx).
		Where("reply_root_id = ?", rootID).
		Order("created_at ASC").
		Offset(offset).Limit(limit).
		Find(&list).Error
	return list, total, err
}

// CountRepliesByParents 统计每个根评论楼内的回复总数（按 reply_root_id）
func (r *commentRepo) CountRepliesByParents(ctx context.Context, rootIDs []int64) map[int64]int64 {
	if len(rootIDs) == 0 {
		return nil
	}
	var rows []struct {
		RootID int64 `gorm:"column:reply_root_id"`
		Cnt      int64 `gorm:"column:cnt"`
	}
	r.db.WithContext(ctx).Model(&model.Comment{}).
		Select("reply_root_id, count(*) as cnt").
		Where("reply_root_id IN ?", rootIDs).
		Group("reply_root_id").
		Find(&rows)
	out := make(map[int64]int64)
	for _, row := range rows {
		out[row.RootID] = row.Cnt
	}
	return out
}

// GetPreviewRepliesForRoots 为每个根评论取该楼内最多 per 条回复（按时间），用于首屏预览
func (r *commentRepo) GetPreviewRepliesForRoots(ctx context.Context, rootIDs []int64, per int) map[int64][]*model.Comment {
	if len(rootIDs) == 0 || per <= 0 {
		return nil
	}
	var all []*model.Comment
	if err := r.db.WithContext(ctx).
		Where("reply_root_id IN ?", rootIDs).
		Order("reply_root_id ASC, created_at ASC").
		Find(&all).Error; err != nil {
		return nil
	}
	out := make(map[int64][]*model.Comment)
	for _, c := range all {
		rootID := c.ReplyRootID
		if len(out[rootID]) < per {
			out[rootID] = append(out[rootID], c)
		}
	}
	return out
}
