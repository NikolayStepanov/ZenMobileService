package http

import (
	"errors"
	"github.com/go-chi/render"
	"net/http"
)

var (
	ErrInvalidInput = errors.New("invalid input body")
	ErrEmptyKey     = errors.New("key can't be empty")
	ErrEmptyValue   = errors.New("value can't be empty")
	ErrEmptyText    = errors.New("text can't be empty")
)

type ErrResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
	ErrorText      string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}
