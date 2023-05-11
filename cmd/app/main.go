package main

import "ZenMobileService/internal/app"

const configPath = "configs/main"

// @title Zen Mobile Service
// @version v0.3
// @description Service uses Redis, PostgreSQL, HMAC-SHA-512
// @host localhost:8080
// @BasePath /
// @query.collection.format multi
func main() {
	app.Run(configPath)
}
