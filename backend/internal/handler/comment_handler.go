package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"simple-comment/internal/model"
	"simple-comment/internal/service"
	"simple-comment/pkg/response"
)

type CommentHandler struct {
	svc *service.CommentService
}

func NewCommentHandler(svc *service.CommentService) *CommentHandler {
	return &CommentHandler{svc: svc}
}

type createCommentRequest struct {
	ArticleID uint   `json:"articleId"`
	UserID    uint   `json:"userId"`
	UserName  string `json:"userName"`
	Avatar    string `json:"avatar"`
	ParentID  uint   `json:"parentId"` // 0=根评论；回复时=被回复的那条评论 id
	Content   string `json:"content"`
}

func (h *CommentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/comments/replies", h.listReplies)
	rg.POST("/comments", h.createComment)
	rg.GET("/comments", h.listComments)
	rg.DELETE("/comments/:id", h.deleteComment)
}

func (h *CommentHandler) createComment(c *gin.Context) {
	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request body")
		return
	}
	if req.ArticleID == 0 {
		req.ArticleID = 1
	}

	comment, err := h.svc.CreateComment(c.Request.Context(), &service.CreateCommentInput{
		ArticleID: req.ArticleID,
		UserID:    req.UserID,
		UserName:  req.UserName,
		Avatar:    req.Avatar,
		ParentID:  req.ParentID,
		Content:   req.Content,
	})
	if err != nil {
		switch err {
		case service.ErrInvalidUserName, service.ErrInvalidContent:
			response.Error(c, http.StatusBadRequest, 4002, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5001, "failed to create comment")
		}
		return
	}

	response.Success(c, comment)
}

// listComments 一级评论分页，每条约带 replyCount 与最多 2 条预览回复
func (h *CommentHandler) listComments(c *gin.Context) {
	articleID, _ := strconv.ParseUint(c.DefaultQuery("articleId", "1"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	items, total, err := h.svc.ListRootsWithPreview(c.Request.Context(), uint(articleID), page, pageSize)
	if err != nil {
		log.Printf("[comments] list error: %v", err)
		response.Error(c, http.StatusInternalServerError, 5002, "failed to list comments")
		return
	}
	if items == nil {
		items = []*model.RootWithPreview{}
	}

	response.Success(c, gin.H{
		"items": items,
		"total": total,
	})
}

// listReplies 某条一级评论下的二级回复分页（查看更多 / 加载更多）
func (h *CommentHandler) listReplies(c *gin.Context) {
	parentID, _ := strconv.ParseInt(c.Query("parentId"), 10, 64)
	if parentID <= 0 {
		response.Error(c, http.StatusBadRequest, 4003, "parentId required")
		return
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	items, total, err := h.svc.GetReplies(c.Request.Context(), parentID, offset, limit)
	if err != nil {
		log.Printf("[comments] replies error: %v", err)
		response.Error(c, http.StatusInternalServerError, 5003, "failed to list replies")
		return
	}
	if items == nil {
		items = []*model.Comment{}
	}

	response.Success(c, gin.H{
		"items": items,
		"total": total,
	})
}

func (h *CommentHandler) deleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, 4004, "invalid comment id")
		return
	}
	if err := h.svc.DeleteComment(c.Request.Context(), id); err != nil {
		log.Printf("[comments] delete error: %v", err)
		response.Error(c, http.StatusInternalServerError, 5004, "failed to delete comment")
		return
	}
	response.Success(c, nil)
}
