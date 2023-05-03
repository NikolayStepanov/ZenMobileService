package app

import (
	"ZenMobileService/internal/config"
	httpDelivery "ZenMobileService/internal/delivery/http"
	"ZenMobileService/internal/server"
	"ZenMobileService/internal/service/redis"
	"context"
	log "github.com/sirupsen/logrus"
	"os/signal"
	"sync"
	"syscall"
)

func Run(configPath string) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	cfg := new(config.Config)
	cfg.Init(configPath)

	redisCache := redis.NewRedisCache(cfg)
	//services := service.NewServices(service.ServicesDependencies{Cache: redisCache})
	log.Println(cfg.Redis.Port)
	log.Println(cfg.Redis.Host)
	log.Println(cfg.Redis.Password)
	redisCache.Set(ctx, "age", 27)
	log.Println(redisCache.Get(ctx, "age"))
	wg := sync.WaitGroup{}
	handlers := httpDelivery.NewHandler()
	//HTTP Server
	srv := server.NewServer(handlers.Init())
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Run(); err != nil {
			log.Printf("error occurred while running http server: %s\n", err.Error())
		}
	}()
	log.Print("Server started")
	<-ctx.Done()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Stop(context.Background()); err != nil {
			log.Printf("error occured on server shutting down: %s", err.Error())
		}
	}()
	wg.Wait()
	log.Print("Server stopped")
}
