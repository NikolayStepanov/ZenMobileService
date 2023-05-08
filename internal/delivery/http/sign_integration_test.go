package http

import (
	"ZenMobileService/internal/service"
	mock_service "ZenMobileService/internal/service/mocks"
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_SignMessage(t *testing.T) {
	type mockBehavior func(r *mock_service.MockSignatureServicer, inputRequest SignRequest)

	tests := []struct {
		name                 string
		inputBody            string
		inputRequest         SignRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ValidatedSuccess",
			inputBody: `{"text": "test", "key": "test123"}`,
			inputRequest: SignRequest{
				Text: "test",
				Key:  "test123",
			},
			mockBehavior: func(r *mock_service.MockSignatureServicer, inputRequest SignRequest) {
				r.EXPECT().GenerateSignature(gomock.Any(), inputRequest.Text, inputRequest.Key).Return("b70c5dd8a6cf30a36976fa2e2", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "b70c5dd8a6cf30a36976fa2e2",
		},
		{
			name:                 "InputBodyError",
			inputBody:            `{"text": "test", "key": test123}`,
			inputRequest:         SignRequest{},
			mockBehavior:         func(r *mock_service.MockSignatureServicer, inputRequest SignRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"invalid input body"}` + "\n",
		},
		{
			name:                 "KeyEmptyError",
			inputBody:            `{"text": "test"}`,
			inputRequest:         SignRequest{},
			mockBehavior:         func(r *mock_service.MockSignatureServicer, inputRequest SignRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"key can't be empty"}` + "\n",
		},
		{
			name:                 "TextEmptyError",
			inputBody:            `{"key": "test123"}`,
			inputRequest:         SignRequest{},
			mockBehavior:         func(r *mock_service.MockSignatureServicer, inputRequest SignRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"text can't be empty"}` + "\n",
		},
		{
			name:      "IncrementValueError",
			inputBody: `{"text": "test", "key": "test123"}`,
			inputRequest: SignRequest{
				Text: "test",
				Key:  "test123",
			},
			mockBehavior: func(r *mock_service.MockSignatureServicer, inputRequest SignRequest) {
				r.EXPECT().GenerateSignature(gomock.Any(), inputRequest.Text, inputRequest.Key).Return("", ErrSignGenerate)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"can't generate signature"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			signService := mock_service.NewMockSignatureServicer(ctl)
			services := &service.Services{SignService: signService}
			test.mockBehavior(signService, test.inputRequest)
			handler := NewHandler(services)
			router := handler.Init()

			router.Post(signRoute+hmacsha512Route, handler.SignMessage)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, signRoute+hmacsha512Route,
				bytes.NewBufferString(test.inputBody))

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
