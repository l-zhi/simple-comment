package service

import (
	"context"
	"errors"
	"strings"

	"simple-comment/internal/model"
	"simple-comment/internal/repository"
)

type ArticleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

var ErrInvalidTitle = errors.New("invalid title")

func (s *ArticleService) CreateArticle(ctx context.Context, title, content string) (*model.Article, error) {
	title = strings.TrimSpace(title)
	if len(title) == 0 || len(title) > 200 {
		return nil, ErrInvalidTitle
	}
	content = strings.TrimSpace(content)
	if len(content) > 50000 {
		content = content[:50000]
	}
	a := &model.Article{Title: title, Content: content}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *ArticleService) DeleteArticle(ctx context.Context, id uint) error {
	return s.repo.DeleteByID(ctx, id)
}

func (s *ArticleService) GetArticle(ctx context.Context, id uint) (*model.Article, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ArticleService) ListArticles(ctx context.Context, page, pageSize int) ([]*model.Article, int64, error) {
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	return s.repo.List(ctx, offset, pageSize)
}
