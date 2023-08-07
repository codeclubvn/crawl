package main

import (
	"crawl/conf"
	"crawl/route"
)

func main() {
	conf.SetEnv()
	app := route.NewService()
	if err := app.Start(); err != nil {
		panic(err)
	}
}
