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
	ErrSignGenerate = errors.New("can't generate signature")
)

type SignRequest struct {
	Text string `json:"text"`
	Key  string `json:"key"`
}

func (h *Handler) initSignRoutes() *chi.Mux {
	signRouter := chi.NewRouter()
	signRouter.Post(hmacsha512, h.SignMessage)
	return signRouter
}

func validateSignReq(signReq *SignRequest) error {
	if signReq.Text == "" {
		return ErrEmptyText
	}
	if signReq.Key == "" {
		return ErrEmptyKey
	}
	return nil
}

// @Summary SignMessage
// @Description Signature message
// @Tags Signature
// @Accept json
// @Produce html
// @Param input body SignRequest true "json request: signature text, key"
// @Success 200 {string} string
// @Failure 400 {object} ErrResponse
// @Router /sign/hmacsha512 [post]
func (h *Handler) SignMessage(w http.ResponseWriter, r *http.Request) {
	err := error(nil)
	signRequest := &SignRequest{}
	valueResponse := ""

	err = render.Decode(r, signRequest)
	if err != nil {
		log.Errorf("can't parse request: %s", err.Error())
		render.Render(w, r, ErrInvalidRequest(ErrInvalidInput))
		return
	}

	err = validateSignReq(signRequest)
	if err != nil {
		log.Errorf("bad request: %v: %s", signRequest, err.Error())
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	valueResponse, err = h.services.SignService.GenerateSignature(r.Context(), signRequest.Text, signRequest.Key)
	if err != nil {
		log.Errorf("can't generate signature: %v: %s", signRequest, err.Error())
		render.Render(w, r, ErrInvalidRequest(ErrSignGenerate))
		return
	}

	render.Status(r, http.StatusOK)
	render.HTML(w, r, valueResponse)
}
