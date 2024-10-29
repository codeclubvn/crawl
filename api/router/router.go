package router

import (
	"crawl/api/controller"
	"crawl/api/middleware"
	"crawl/docs"
	"crawl/domain"
	"crawl/initialization"
	"crawl/pkg/crawl_data"
	"crawl/repository"
	"crawl/usecase"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func SetUp(env *initialization.Database, timeout time.Duration, db *mongo.Database, client *mongo.Client, gin *gin.Engine, cacheTTL time.Duration) {
	publicRouterV1 := gin.Group("/api/v1")

	// Middleware
	publicRouterV1.Use(
		middleware.Recover(),
		gzip.Gzip(gzip.DefaultCompression,
			gzip.WithExcludedPaths([]string{",*"})),
	)

	CompanyRouter(env, timeout, db, publicRouterV1, cacheTTL)
	SwaggerRouter(env, timeout, db, gin.Group(""))

	// Đếm các route
	routeCount := countRoutes(gin)
	fmt.Printf("The number of API endpoints: %d\n", routeCount)
}

func countRoutes(r *gin.Engine) int {
	count := 0
	routes := r.Routes()
	for range routes {
		count++
	}
	return count
}

func init() {
	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/"),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.DeepLinking(true),
		ginSwagger.PersistAuthorization(true),
	)

	// Save pprof handlers first.
	pprofMux := http.DefaultServeMux
	http.DefaultServeMux = http.NewServeMux()

	// Pprof server.
	go func() {
		fmt.Println(http.ListenAndServe("localhost:8000", pprofMux))
	}()
}

func SwaggerRouter(env *initialization.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup) {
	router := group.Group("")

	docs.SwaggerInfo.BasePath = ""
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Thực hiện tự động chuyển hướng khi chạy chương trình
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/swagger/index.html")
	})
}

func CompanyRouter(env *initialization.Database, timeout time.Duration, db *mongo.Database, group *gin.RouterGroup, cacheTTL time.Duration) {
	var companies domain.Companies

	co := repository.NewCompanyRepository("company", db)
	sc := crawl_data.NewCompaniesCrawl(companies, co)

	company := &controller.CompanyController{
		CompanyUseCase: usecase.NewCompanyUseCase(co, timeout, sc),
		Database:       env,
	}

	router := group.Group("/companies")
	router.POST("/create", company.CreateOne)
	router.GET("/get/all", company.GetAll)
}
