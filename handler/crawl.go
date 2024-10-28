package handler

import (
	"crawl/model"
	"crawl/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Crawl struct {
	service service.IUser
}

func NewCrawl(service service.IUser) *Crawl {
	return &Crawl{
		service: service,
	}
}

func (h *Crawl) CrawlYellowPage(ctx *gin.Context) {
	var req model.CrawlInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CrawlYellowPage(req.Url); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Ví dụ: trả về JSON
	ctx.JSON(200, gin.H{
		"message": "success",
	})
}

func (h *Crawl) GetOne(ctx *gin.Context) {
	var req model.GetOneInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.service.GetOne(req.Filter)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, result)
}

func (h *Crawl) GetList(ctx *gin.Context) {
	var req model.GetListInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.service.GetList(req.Filter, req.Limit, req.Skip)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, result)
}
