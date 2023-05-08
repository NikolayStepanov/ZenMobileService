package http

import (
	_ "ZenMobileService/docs"
	"ZenMobileService/internal/service"
	"github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	redisRoute   = "/redis"
	signRoute    = "/sign"
	swaggerRoute = "/swagger/*"
	swaggerURL   = "http://localhost:8080/swagger/doc.json"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() *chi.Mux {
	router := chi.NewRouter()
	log := logrus.New()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(logger.Logger("router", log))
	router.Use(render.SetContentType(render.ContentTypeJSON))
	h.MountRoutes(router)
	return router
}

func (h *Handler) MountRoutes(router *chi.Mux) {
	router.Mount(redisRoute, h.initRedisRoutes())
	router.Mount(signRoute, h.initSignRoutes())
	router.Get(swaggerRoute, httpSwagger.Handler(
		httpSwagger.URL(swaggerURL),
	))
}
