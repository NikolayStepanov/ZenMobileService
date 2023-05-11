package http

import (
	"ZenMobileService/internal/domain"
	"ZenMobileService/internal/service"
	mock_service "ZenMobileService/internal/service/mocks"
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreateUser(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUsersServicer, inputRequest UserCreateRequest)

	tests := []struct {
		name                 string
		inputBody            string
		inputRequest         UserCreateRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ValidatedSuccess",
			inputBody: `{"name": "Alex", "age": 27}`,
			inputRequest: UserCreateRequest{
				Name: "Alex",
				Age:  27,
			},
			mockBehavior: func(r *mock_service.MockUsersServicer, inputRequest UserCreateRequest) {
				user := domain.NewUser(0, inputRequest.Name, inputRequest.Age)
				r.EXPECT().CreateUser(gomock.Any(), *user).Return(2, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":2}` + "\n",
		},
		{
			name:      "NameEmptyError",
			inputBody: `{"name": "", "age": 27}`,
			inputRequest: UserCreateRequest{
				Name: "",
				Age:  27,
			},
			mockBehavior:         func(r *mock_service.MockUsersServicer, inputRequest UserCreateRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"name can't be empty"}` + "\n",
		},
		{
			name:                 "NameIncorrectError",
			inputBody:            `{"name": 23213, "age": 27}`,
			inputRequest:         UserCreateRequest{},
			mockBehavior:         func(r *mock_service.MockUsersServicer, inputRequest UserCreateRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"invalid input body"}` + "\n",
		},
		{
			name:                 "AgeIncorrectError",
			inputBody:            `{"name":"Alex" , "age": d23}`,
			inputRequest:         UserCreateRequest{},
			mockBehavior:         func(r *mock_service.MockUsersServicer, inputRequest UserCreateRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"invalid input body"}` + "\n",
		},
		{
			name:      "CreateUserError",
			inputBody: `{"name": "Alex", "age": 27}`,
			inputRequest: UserCreateRequest{
				Name: "Alex",
				Age:  27,
			},
			mockBehavior: func(r *mock_service.MockUsersServicer, inputRequest UserCreateRequest) {
				user := domain.NewUser(0, inputRequest.Name, inputRequest.Age)
				r.EXPECT().CreateUser(gomock.Any(), *user).Return(0, ErrCreateUser)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"can't create user"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			usersService := mock_service.NewMockUsersServicer(ctl)
			services := &service.Services{UsersService: usersService}
			test.mockBehavior(usersService, test.inputRequest)
			handler := NewHandler(services)
			router := handler.Init()

			router.Post(postgresRoute+usersRoute, handler.CreateUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, postgresRoute+usersRoute,
				bytes.NewBufferString(test.inputBody))

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_GetUser(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUsersServicer, inputUserId int)

	tests := []struct {
		name                 string
		inputRequest         string
		inputUserId          int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "ValidatedSuccess",
			inputRequest: "20",
			inputUserId:  20,
			mockBehavior: func(r *mock_service.MockUsersServicer, inputUserId int) {
				r.EXPECT().GetUser(gomock.Any(), inputUserId).Return(*domain.NewUser(20, "Alex", 27), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":20,"name":"Alex","age":27}` + "\n",
		},
		{
			name:                 "UserIdIncorrectError",
			inputRequest:         "df20",
			inputUserId:          20,
			mockBehavior:         func(r *mock_service.MockUsersServicer, inputUserId int) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"can't parse the user ID"}` + "\n",
		},
		{
			name:         "GetUserError",
			inputRequest: "20",
			inputUserId:  20,
			mockBehavior: func(r *mock_service.MockUsersServicer, inputUserId int) {
				r.EXPECT().GetUser(gomock.Any(), inputUserId).Return(domain.User{}, ErrGetUser)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"can't get information user"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			usersService := mock_service.NewMockUsersServicer(ctl)
			services := &service.Services{UsersService: usersService}
			test.mockBehavior(usersService, test.inputUserId)
			handler := NewHandler(services)
			router := handler.Init()

			router.Get(postgresRoute+usersRoute+getUserRoute, handler.GetUserInformation)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, postgresRoute+usersRoute+"/"+test.inputRequest, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
