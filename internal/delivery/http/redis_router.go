package http

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	slash       = "/"
	incrRoute   = "/incr"
	getKeyRoute = "/{key}"
)

var (
	ErrIncrementValue = errors.New("can't increment value by key")
	ErrSave           = errors.New("can't save value by key")
	ErrRead           = errors.New("can't read value by key")
)

type IncrementRequest struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

type SaveValueRequest struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type ValueIncrementResponse struct {
	Value int64 `json:"value"`
}

type ValueResponse struct {
	Value any `json:"value"`
}

func validateSaveReq(reqSave *SaveValueRequest) error {
	err := validateKeyParam(reqSave.Key)
	if err != nil {
		return err
	}
	return nil
}

func validateIncrementReq(reqIncr *IncrementRequest) error {
	err := validateKeyParam(reqIncr.Key)
	if err != nil {
		return err
	}
	if reqIncr.Value == 0 {
		return ErrEmptyValue
	}
	return nil
}

func validateKeyParam(key string) error {
	if key == "" {
		return ErrEmptyKey
	}
	return nil
}

func (h *Handler) initRedisRoutes() *chi.Mux {
	redisRouter := chi.NewRouter()
	redisRouter.Post(incrRoute, h.IncrementValueByKey)
	redisRouter.Post(slash, h.SaveValueByKey)
	redisRouter.Get(getKeyRoute, h.ReadValueByKey)
	return redisRouter
}

// @Summary IncrementValueByKey
// @Description Increment value by key if value is stored input redis
// @Tags Redis
// @Accept json
// @Produce json
// @Param input body IncrementRequest true "json request: increment value by key"
// @Success 200 {object} ValueIncrementResponse
// @Failure 400 {object} ErrResponse
// @Router /redis/incr [post]
func (h *Handler) IncrementValueByKey(w http.ResponseWriter, r *http.Request) {
	err := error(nil)
	incrRequest := &IncrementRequest{}
	valueResponse := &ValueIncrementResponse{}

	err = render.Decode(r, incrRequest)
	if err != nil {
		log.Errorf("can't parse request: %s", err.Error())
		render.Render(w, r, ErrInvalidRequest(ErrInvalidInput))
		return
	}

	err = validateIncrementReq(incrRequest)
	if err != nil {
		log.Errorf("bad request: %v: %s", incrRequest, err.Error())
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	valueResponse.Value, err = h.services.CacheService.IncrementValueByKey(r.Context(), incrRequest.Key, incrRequest.Value)
	if err != nil {
		log.Errorf("redis: can't increment value by key: %v: %s", incrRequest, err.Error())
		render.Render(w, r, ErrInvalidRequest(ErrIncrementValue))
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
	saveValueRequest := &SaveValueRequest{}
	messageResponse := ""

	err = render.Decode(r, &saveValueRequest)
	if err != nil {
		log.Errorf("can't parse req: %s", err.Error())
		render.Render(w, r, ErrInvalidRequest(ErrInvalidInput))
		return
	}

	err = validateSaveReq(saveValueRequest)
	if err != nil {
		log.Errorf("bad req: %v: %s", saveValueRequest, err.Error())
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	err = h.services.CacheService.SetValueByKey(r.Context(), saveValueRequest.Key, saveValueRequest.Value)
	if err != nil {
		log.Errorf("redis: can't save value by key: %v: %s", saveValueRequest, err.Error())
		render.Render(w, r, ErrInvalidRequest(ErrSave))
		return
	}
	messageResponse = fmt.Sprintf("Key = %s Value = %v saved in Redis", saveValueRequest.Key, saveValueRequest.Value)

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
	value := any(0)
	err := error(nil)
	valueResponse := &ValueResponse{}
	key := chi.URLParam(r, "key")

	err = validateKeyParam(key)
	if err != nil {
		log.Errorf("bad req: %s: %s", key, err.Error())
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	value, err = h.services.CacheService.GetValueByKey(r.Context(), key)
	if err != nil {
		log.Errorf("redis: can't read value by key: %s: %s", key, err.Error())
		render.Render(w, r, ErrInvalidRequest(ErrRead))
		return
	}
	valueResponse.Value = value

	render.Status(r, http.StatusOK)
	render.JSON(w, r, valueResponse)
}
