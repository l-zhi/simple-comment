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

type ArticleHandler struct {
	svc *service.ArticleService
}

func NewArticleHandler(svc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

type createArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *ArticleHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/posts", h.listPosts)
	rg.POST("/posts", h.createPost)
	rg.GET("/posts/:id", h.getPost)
	rg.DELETE("/posts/:id", h.deletePost)
}

func (h *ArticleHandler) listPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	items, total, err := h.svc.ListArticles(c.Request.Context(), page, pageSize)
	if err != nil {
		log.Printf("[posts] list error: %v", err)
		response.Error(c, http.StatusInternalServerError, 5101, "failed to list posts")
		return
	}
	if items == nil {
		items = []*model.Article{}
	}
	response.Success(c, gin.H{"items": items, "total": total})
}

func (h *ArticleHandler) getPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, 4101, "invalid post id")
		return
	}
	a, err := h.svc.GetArticle(c.Request.Context(), uint(id))
	if err != nil {
		log.Printf("[posts] get error: %v", err)
		response.Error(c, http.StatusNotFound, 4102, "post not found")
		return
	}
	response.Success(c, a)
}

func (h *ArticleHandler) createPost(c *gin.Context) {
	var req createArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request body")
		return
	}
	a, err := h.svc.CreateArticle(c.Request.Context(), req.Title, req.Content)
	if err != nil {
		if err == service.ErrInvalidTitle {
			response.Error(c, http.StatusBadRequest, 4002, err.Error())
			return
		}
		log.Printf("[posts] create error: %v", err)
		response.Error(c, http.StatusInternalServerError, 5001, "failed to create post")
		return
	}
	response.Success(c, a)
}

func (h *ArticleHandler) deletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, 4103, "invalid post id")
		return
	}
	if err := h.svc.DeleteArticle(c.Request.Context(), uint(id)); err != nil {
		log.Printf("[posts] delete error: %v", err)
		response.Error(c, http.StatusInternalServerError, 5102, "failed to delete post")
		return
	}
	response.Success(c, nil)
}
