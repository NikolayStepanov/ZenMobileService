package http

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
)

type IncrementRequest struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

type ValueIncrementResponse struct {
	Value int64 `json:"value"`
}

type ValueResponse struct {
	Value any `json:"value"`
}

type PutValueRequest struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (h *Handler) initRedisRoutes() *chi.Mux {
	redisRouter := chi.NewRouter()
	redisRouter.Post("/incr", h.IncrementValueByKey)
	redisRouter.Put("/", h.PutValueByKey)
	redisRouter.Get("/{key}", h.GetValueByKey)
	return redisRouter
}

func (h *Handler) IncrementValueByKey(w http.ResponseWriter, r *http.Request) {
	err := error(nil)
	incrRequest := IncrementRequest{}
	valueResponse := ValueIncrementResponse{}

	render.Decode(r, &incrRequest)
	valueResponse.Value, err = h.cache.IncrementValueByKey(r.Context(), incrRequest.Key, incrRequest.Value)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, valueResponse)
}

func (h *Handler) PutValueByKey(w http.ResponseWriter, r *http.Request) {
	err := error(nil)
	putValueRequest := PutValueRequest{}
	messageResponse := ""

	render.Decode(r, &putValueRequest)
	err = h.cache.SetValueByKey(r.Context(), putValueRequest.Key, putValueRequest.Value)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	messageResponse = fmt.Sprintf("Key = %s Value = %v put in Redis", putValueRequest.Key, putValueRequest.Value)

	render.Status(r, http.StatusOK)
	render.HTML(w, r, messageResponse)
}

func (h *Handler) GetValueByKey(w http.ResponseWriter, r *http.Request) {
	var value any
	err := error(nil)
	valueResponse := ValueResponse{}
	key := chi.URLParam(r, "key")
	if key != "" {
		value, err = h.cache.GetValueByKey(r.Context(), key)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		valueResponse.Value = value
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, valueResponse)
}
