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
	rp := repo.NewRepo(db)

	crawlService := service.NewCrawl(rp)
	crawl := handler.NewCrawl(crawlService)
	health := handler.NewHealth()

	router := s.Router
	v1 := router.Group("/v1")

	// crawl
	v1.POST("/crawl/trang-vang", crawl.CrawlYellowPage)
	v1.GET("company/get-one", crawl.GetOne)
	v1.GET("company/get-list", crawl.GetList)
	router.POST("/check-health", health.CheckHealth)

	return &s
}
