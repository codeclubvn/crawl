package main

import (
	"crawl/api/router"
	"crawl/infrastructor"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {

	app, client := infrastructor.App()

	env := app.Env

	db := app.MongoDB.Database(env.DBName)

	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second
	cacheTTL := time.Minute * 5

	_gin := gin.Default()

	router.SetUp(env, timeout, db, client, _gin, cacheTTL)
	fmt.Println("Location Server Web of us: http://localhost:8080")
	err := _gin.Run(env.ServerAddress)
	if err != nil {
		return
	}
}
