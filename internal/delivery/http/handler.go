package http

import (
	_ "ZenMobileService/docs"
	"ZenMobileService/internal/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	cache service.CacheServicer
}

func NewHandler(cacheService service.CacheServicer) *Handler {
	return &Handler{cache: cacheService}
}

func (h *Handler) Init() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)
	router.Use(render.SetContentType(render.ContentTypeJSON))
	h.MountRoutes(router)
	return router
}

func (h *Handler) MountRoutes(router *chi.Mux) {
	router.Mount("/redis", h.initRedisRoutes())
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
}
