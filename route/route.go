package route

import (
	"crawl/conf"
	"crawl/handler"
	"crawl/repo"
	"crawl/service"
)

type Service struct {
	*conf.App
}

func NewService() *Service {
	s := Service{
		conf.NewApp(),
	}

	db := s.GetDB()
	repo := repo.NewRepo(db)

	userService := service.NewUser(repo)
	user := handler.NewUser(userService)
	health := handler.NewHealth()

	router := s.Router
	v1 := router.Group("/v1")

	// user
	v1.POST("/login", user.Login)
	router.POST("/check-health", health.CheckHealth)

	return &s
}
