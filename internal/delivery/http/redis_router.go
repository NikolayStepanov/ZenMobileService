package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type IncrementRequest struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

type ValueResponse struct {
	Value int64 `json:"value"`
}

func (h *Handler) initRedisRoutes() *chi.Mux {
	redisRouter := chi.NewRouter()
	redisRouter.Post("/incr", h.IncrementValueByKey)
	return redisRouter
}

func (h *Handler) IncrementValueByKey(w http.ResponseWriter, r *http.Request) {
	err := error(nil)
	incrRequest := IncrementRequest{}
	valueResponse := ValueResponse{}

	render.Decode(r, &incrRequest)
	valueResponse.Value, err = h.cache.IncrementValueByKey(r.Context(), incrRequest.Key, incrRequest.Value)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, valueResponse)
}
