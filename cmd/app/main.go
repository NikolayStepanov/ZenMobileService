package main

import "ZenMobileService/internal/app"

const configPath = "configs/main"

// @title Zen Mobile Service
// @version 0.1
// @description Service use Redis
// @host localhost:8080
// @BasePath /
// @query.collection.format multi
func main() {
	app.Run(configPath)
}
