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

type SaveValueRequest struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (h *Handler) initRedisRoutes() *chi.Mux {
	redisRouter := chi.NewRouter()
	redisRouter.Post("/incr", h.IncrementValueByKey)
	redisRouter.Post("/", h.SaveValueByKey)
	redisRouter.Get("/{key}", h.ReadValueByKey)
	return redisRouter
}

// @Summary IncrementValueByKey
// @Description Increment value by key if value is stored in redis
// @Tags Redis
// @Accept json
// @Produce json
// @Param input body IncrementRequest true "json request: increment value by key"
// @Success 200 {object} ValueIncrementResponse
// @Failure 400 {object} ErrResponse
// @Router /redis/incr [post]
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

// @Summary SaveValueByKey
// @Description Saving a new key with a value
// @Tags Redis
// @Accept json
// @Produce html
// @Param input body SaveValueRequest true "json request: save value"
// @Success 200 {string} string
// @Failure 400 {object} ErrResponse
// @Router /redis/ [post]
func (h *Handler) SaveValueByKey(w http.ResponseWriter, r *http.Request) {
	err := error(nil)
	createValueRequest := SaveValueRequest{}
	messageResponse := ""

	render.Decode(r, &createValueRequest)
	err = h.cache.SetValueByKey(r.Context(), createValueRequest.Key, createValueRequest.Value)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	messageResponse = fmt.Sprintf("Key = %s Value = %v saved in Redis", createValueRequest.Key, createValueRequest.Value)

	render.Status(r, http.StatusOK)
	render.HTML(w, r, messageResponse)
}

// @Summary ReadValueByKey
// @Description Getting value by key
// @Tags Redis
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} ValueResponse
// @Failure 400 {object} ErrResponse
// @Router /redis/{key} [get]
func (h *Handler) ReadValueByKey(w http.ResponseWriter, r *http.Request) {
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
