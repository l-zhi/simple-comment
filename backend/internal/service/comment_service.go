package service

import (
	"context"
	"errors"
	"strings"

	"simple-comment/internal/model"
	"simple-comment/internal/repository"
)

type CommentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

type CreateCommentInput struct {
	ArticleID uint
	UserID    uint
	UserName  string
	Avatar    string
	ParentID  uint // 0=根评论；回复时=被回复的那条评论 id
	Content   string
}

var (
	ErrInvalidUserName = errors.New("invalid user name")
	ErrInvalidContent  = errors.New("invalid content")
)

func (s *CommentService) CreateComment(ctx context.Context, in *CreateCommentInput) (*model.Comment, error) {
	userName := strings.TrimSpace(in.UserName)
	if len(userName) == 0 || len(userName) > 50 {
		return nil, ErrInvalidUserName
	}
	content := strings.TrimSpace(in.Content)
	if len(content) == 0 || len(content) > 2000 {
		return nil, ErrInvalidContent
	}
	if len(in.Avatar) > 255 {
		in.Avatar = in.Avatar[:255]
	}

	var replyRootID int64 = 0
	if in.ParentID != 0 {
		parent, err := s.repo.GetCommentByID(ctx, int64(in.ParentID))
		if err == nil && parent != nil {
			if parent.ParentID == 0 {
				replyRootID = parent.ID
			} else if parent.ReplyRootID != 0 {
				replyRootID = parent.ReplyRootID
			} else {
				replyRootID = int64(parent.ParentID)
			}
		}
	}

	c := &model.Comment{
		ArticleID: in.ArticleID,
		UserID:    in.UserID,
		UserName:  userName,
		Avatar:    in.Avatar,
		ParentID:  in.ParentID,
		ReplyRootID: replyRootID,
		Content:   content,
		Status:    model.StatusNormal,
		Likes:     0,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

const (
	RootsPageSize   = 10
	RepliesPageSize = 20
)

const replyToContentMaxLen = 80

// enrichRepliesWithTarget 根据 parent_id（被直接回复的评论 id）查被回复记录，填 replyToUserName + replyToContent；
// 被回复评论已删除则不展示该条回复
func (s *CommentService) enrichRepliesWithTarget(ctx context.Context, root *model.Comment, replies []*model.Comment) []*model.Comment {
	if len(replies) == 0 {
		return nil
	}
	ids := make([]int64, 0, len(replies)+1)
	if root != nil {
		ids = append(ids, root.ID)
	}
	for _, r := range replies {
		tid := int64(r.ParentID)
		if tid > 0 {
			ids = append(ids, tid)
		}
	}
	byID := s.repo.GetCommentsByIDs(ctx, ids)
	if root != nil && byID != nil {
		byID[root.ID] = root
	}

	var out []*model.Comment
	for _, r := range replies {
		tid := int64(r.ParentID)
		if tid == 0 {
			r.ReplyToUserName = ""
			r.ReplyToContent = ""
			out = append(out, r)
			continue
		}
		target := byID[tid]
		if target == nil {
			continue
		}
		r.ReplyToUserName = target.UserName
		r.ReplyToContent = target.Content
		if len(r.ReplyToContent) > replyToContentMaxLen {
			r.ReplyToContent = r.ReplyToContent[:replyToContentMaxLen] + "…"
		}
		out = append(out, r)
	}
	return out
}

// ListRootsWithPreview 一级评论分页，每条约带回复总数（不再返回预览回复）
func (s *CommentService) ListRootsWithPreview(ctx context.Context, articleID uint, page, pageSize int) ([]*model.RootWithPreview, int64, error) {
	if pageSize <= 0 || pageSize > 50 {
		pageSize = RootsPageSize
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize

	roots, total, err := s.repo.ListRootsByArticle(ctx, articleID, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	if len(roots) == 0 {
		return []*model.RootWithPreview{}, total, nil
	}

	rootIDs := make([]int64, 0, len(roots))
	for _, r := range roots {
		rootIDs = append(rootIDs, r.ID)
	}
	replyCounts := s.repo.CountRepliesByParents(ctx, rootIDs)

	out := make([]*model.RootWithPreview, 0, len(roots))
	for _, r := range roots {
		out = append(out, &model.RootWithPreview{
			Comment:    *r,
			ReplyCount: replyCounts[r.ID],
			Replies:    []model.Comment{},
		})
	}
	return out, total, nil
}

// GetReplies 某条一级评论下的二级回复分页；被回复评论已删除的回复不展示
func (s *CommentService) GetReplies(ctx context.Context, parentID int64, offset, limit int) ([]*model.Comment, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = RepliesPageSize
	}
	list, total, err := s.repo.ListRepliesByParent(ctx, parentID, offset, limit)
	if err != nil || len(list) == 0 {
		return list, total, err
	}
	root, _ := s.repo.GetCommentByID(ctx, parentID)
	enriched := s.enrichRepliesWithTarget(ctx, root, list)
	return enriched, total, nil
}

// DeleteComment 软删除评论
func (s *CommentService) DeleteComment(ctx context.Context, id int64) error {
	return s.repo.DeleteByID(ctx, id)
}
