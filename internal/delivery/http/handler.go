package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
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
	router.Get("/", h.handleConnection)
}

func (h *Handler) handleConnection(w http.ResponseWriter, r *http.Request) {
	render.HTML(w, r, "Hello")
}
