package http

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	hmacsha512 = "/hmacsha512"
)

var (
	ErrSignatureGenerate = errors.New("can't generate signature")
)

type SignatureRequest struct {
	Text string `json:"text"`
	Key  string `json:"key"`
}

func (h *Handler) initSignRoutes() *chi.Mux {
	signRouter := chi.NewRouter()
	signRouter.Post(hmacsha512, h.SignatureMessage)
	return signRouter
}

func validateSignatureReq(signReq *SignatureRequest) error {
	if signReq.Text == "" {
		return ErrEmptyText
	}
	if signReq.Key == "" {
		return ErrEmptyKey
	}
	return nil
}

func (h *Handler) SignatureMessage(w http.ResponseWriter, r *http.Request) {
	err := error(nil)
	signRequest := &SignatureRequest{}
	valueResponse := ""

	err = render.Decode(r, signRequest)
	if err != nil {
		log.Errorf("can't parse request: %s", err.Error())
		render.Render(w, r, ErrInvalidRequest(ErrInvalidInput))
		return
	}

	err = validateSignatureReq(signRequest)
	if err != nil {
		log.Errorf("bad request: %v: %s", signRequest, err.Error())
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	valueResponse, err = h.sign.GenerateSignature(r.Context(), signRequest.Text, signRequest.Key)
	if err != nil {
		log.Errorf("can't generate signature: %v: %s", signRequest, err.Error())
		render.Render(w, r, ErrInvalidRequest(ErrSignatureGenerate))
		return
	}

	render.Status(r, http.StatusOK)
	render.HTML(w, r, valueResponse)
}
