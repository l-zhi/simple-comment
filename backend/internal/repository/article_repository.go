package repository

import (
	"context"

	"gorm.io/gorm"

	"simple-comment/internal/model"
)

type ArticleRepository interface {
	Create(ctx context.Context, a *model.Article) error
	DeleteByID(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*model.Article, error)
	List(ctx context.Context, offset, limit int) ([]*model.Article, int64, error)
}

type articleRepo struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepo{db: db}
}

func (r *articleRepo) Create(ctx context.Context, a *model.Article) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *articleRepo) DeleteByID(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Article{}).Error
}

func (r *articleRepo) GetByID(ctx context.Context, id uint) (*model.Article, error) {
	var a model.Article
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *articleRepo) List(ctx context.Context, offset, limit int) ([]*model.Article, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []*model.Article
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&list).Error
	return list, total, err
}
