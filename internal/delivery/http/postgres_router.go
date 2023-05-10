package http

import (
	"ZenMobileService/internal/domain"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	usersRoute   = "/users"
	getUserRoute = "/{userID}"
)

var (
	ErrEmptyName = errors.New("name can't be empty")
)

type UserCreateRequest struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}

type UserCreateResponse struct {
	Id int `json:"id"`
}

type UserInformationResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (h *Handler) initPostgresRoutes() *chi.Mux {
	postgresRouter := chi.NewRouter()
	postgresRouter.Post(usersRoute, h.CreateUser)
	postgresRouter.Get(usersRoute+getUserRoute, h.GetUserInformation)
	return postgresRouter
}

func validateCreateUserReq(reqCreateUser *UserCreateRequest) error {
	if reqCreateUser.Name == "" {
		return ErrEmptyName
	}
	return nil
}

// @Summary CreateUser
// @Description Сreating a new user
// @Tags Postgres
// @Accept json
// @Produce json
// @Param input body UserCreateRequest true "json information user"
// @Success 200 {object} UserCreateResponse
// @Failure 400 {object} ErrResponse
// @Router /postgres/users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	userCreateRequest := &UserCreateRequest{}
	userCreateResponse := &UserCreateResponse{}
	user := domain.User{}
	userId := 0

	err := render.Decode(r, &userCreateRequest)
	if err != nil {
		log.Errorf("can't parse request: %s", err.Error())
		render.Render(w, r, ErrInvalidRequest(ErrInvalidInput))
		return
	}

	err = validateCreateUserReq(userCreateRequest)
	if err != nil {
		log.Errorf("bad request: %v: %s", userCreateRequest, err.Error())
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	user.SetName(userCreateRequest.Name)
	user.SetAge(userCreateRequest.Age)
	userId, err = h.services.UsersService.CreateUser(r.Context(), user)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	userCreateResponse.Id = userId
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, userCreateResponse)
}

// @Summary GetUserInformation
// @Description Getting information about the user
// @Tags Postgres
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserInformationResponse
// @Failure 400 {object} ErrResponse
// @Router /postgres/users/{id} [get]
func (h *Handler) GetUserInformation(w http.ResponseWriter, r *http.Request) {
	userInformationResponse := &UserInformationResponse{}

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	user, err := h.services.UsersService.GetUser(r.Context(), userID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	userInformationResponse.Id = user.ID()
	userInformationResponse.Age = int(user.Age())
	userInformationResponse.Name = user.Name()

	render.Status(r, http.StatusOK)
	render.JSON(w, r, userInformationResponse)
}
