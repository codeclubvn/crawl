package controller

import (
	"crawl/initialization"
	"crawl/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CompanyController struct {
	Database       *initialization.Database
	CompanyUseCase usecase.ICompanyUseCase
}

// CreateOne create the company's information
// @Summary Create Company Information
// @Description Create the company's information
// @Tags Company
// @Accept json
// @Produce json
// @Router /api/v1/companies/create [post]
// @Security CookieAuth
func (b *CompanyController) CreateOne(ctx *gin.Context) {
	page := []string{"https://trangvangvietnam.com/categories/484645/logistics-dich-vu-logistics.html"}

	err := b.CompanyUseCase.CreateOne(ctx, page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
