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
