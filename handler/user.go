package handler

import (
	"crawl/model"
	"crawl/service"
	"github.com/gin-gonic/gin"
)

type User struct {
	service service.IUser
}

func NewUser(service service.IUser) *User {
	return &User{
		service: service,
	}
}

func (h *User) Login(c *gin.Context) {
	// Lấy thông tin từ request
	userRequest := model.UserRequest{}
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	// call service
	user, err := h.service.Login(c, userRequest)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Trả về thông tin của bản ghi
	c.JSON(200, user)
}
