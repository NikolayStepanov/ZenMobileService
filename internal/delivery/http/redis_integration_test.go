package http

import (
	service "ZenMobileService/internal/service"
	mock_service "ZenMobileService/internal/service/mocks"
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_IncrementValueByKey(t *testing.T) {
	type mockBehavior func(r *mock_service.MockCacheServicer, inputRequest IncrementRequest)

	tests := []struct {
		name                 string
		inputBody            string
		inputRequest         IncrementRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ValidatedSuccess",
			inputBody: `{"key": "age", "value": 10}`,
			inputRequest: IncrementRequest{
				Key:   "age",
				Value: 10,
			},
			mockBehavior: func(r *mock_service.MockCacheServicer, inputRequest IncrementRequest) {
				r.EXPECT().IncrementValueByKey(gomock.Any(), inputRequest.Key, inputRequest.Value).Return(int64(37), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"value":37}` + "\n",
		},
		{
			name:                 "InputBodyError",
			inputBody:            `{"key": "age", "value": "10"}`,
			inputRequest:         IncrementRequest{},
			mockBehavior:         func(r *mock_service.MockCacheServicer, inputRequest IncrementRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"invalid input body"}` + "\n",
		},
		{
			name:                 "KeyEmptyError",
			inputBody:            `{"key": "", "value": 10}`,
			inputRequest:         IncrementRequest{},
			mockBehavior:         func(r *mock_service.MockCacheServicer, inputRequest IncrementRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"key can't be empty"}` + "\n",
		},
		{
			name:                 "ValueEmptyError",
			inputBody:            `{"key": "age", "value": 0}`,
			inputRequest:         IncrementRequest{},
			mockBehavior:         func(r *mock_service.MockCacheServicer, inputRequest IncrementRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"value can't be empty"}` + "\n",
		},
		{
			name:      "IncrementValueError",
			inputBody: `{"key": "age", "value": 10}`,
			inputRequest: IncrementRequest{
				Key:   "age",
				Value: 10,
			},
			mockBehavior: func(r *mock_service.MockCacheServicer, inputRequest IncrementRequest) {
				r.EXPECT().IncrementValueByKey(gomock.Any(), inputRequest.Key, inputRequest.Value).Return(int64(0), ErrIncrementValue)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"can't increment value by key"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			cacheService := mock_service.NewMockCacheServicer(ctl)
			services := &service.Services{CacheService: cacheService}
			test.mockBehavior(cacheService, test.inputRequest)
			handler := NewHandler(services)
			router := handler.Init()

			router.Post(redisRoute+incrRoute, handler.IncrementValueByKey)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, redisRoute+incrRoute,
				bytes.NewBufferString(test.inputBody))

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_SaveValueByKey(t *testing.T) {
	type mockBehavior func(r *mock_service.MockCacheServicer, inputRequest SaveValueRequest)

	tests := []struct {
		name                 string
		inputBody            string
		inputRequest         SaveValueRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ValidatedSuccess",
			inputBody: `{"key": "age", "value": "27"}`,
			inputRequest: SaveValueRequest{
				Key:   "age",
				Value: "27",
			},
			mockBehavior: func(r *mock_service.MockCacheServicer, inputRequest SaveValueRequest) {
				r.EXPECT().SetValueByKey(gomock.Any(), inputRequest.Key, inputRequest.Value).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `Key = age Value = 27 saved in Redis`,
		},
		{
			name:                 "InputBodyError",
			inputBody:            `{"key": age, "value": 10}`,
			inputRequest:         SaveValueRequest{},
			mockBehavior:         func(r *mock_service.MockCacheServicer, inputRequest SaveValueRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"invalid input body"}` + "\n",
		},
		{
			name:                 "KeyEmptyError",
			inputBody:            `{"key": "", "value": 10}`,
			inputRequest:         SaveValueRequest{},
			mockBehavior:         func(r *mock_service.MockCacheServicer, inputRequest SaveValueRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"key can't be empty"}` + "\n",
		},
		{
			name:      "SaveError",
			inputBody: `{"key": "age", "value": "10"}`,
			inputRequest: SaveValueRequest{
				Key:   "age",
				Value: "10",
			},
			mockBehavior: func(r *mock_service.MockCacheServicer, inputRequest SaveValueRequest) {
				r.EXPECT().SetValueByKey(gomock.Any(), inputRequest.Key, inputRequest.Value).Return(ErrSave)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"can't save value by key"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			cacheService := mock_service.NewMockCacheServicer(ctl)
			services := &service.Services{CacheService: cacheService}
			test.mockBehavior(cacheService, test.inputRequest)
			handler := NewHandler(services)
			router := handler.Init()

			router.Post(redisRoute+slash, handler.SaveValueByKey)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, redisRoute+slash,
				bytes.NewBufferString(test.inputBody))

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_ReadValueByKey(t *testing.T) {
	type mockBehavior func(r *mock_service.MockCacheServicer, inputRequest string)

	tests := []struct {
		name                 string
		inputRequest         string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "ValidatedSuccess",
			inputRequest: "age",
			mockBehavior: func(r *mock_service.MockCacheServicer, inputRequest string) {
				r.EXPECT().GetValueByKey(gomock.Any(), inputRequest).Return(27, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"value":27}` + "\n",
		},
		{
			name:         "KeyIncorrectError",
			inputRequest: "a",
			mockBehavior: func(r *mock_service.MockCacheServicer, inputRequest string) {
				r.EXPECT().GetValueByKey(gomock.Any(), inputRequest).Return(0, ErrRead)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"status":"Invalid request.","error":"can't read value by key"}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			cacheService := mock_service.NewMockCacheServicer(ctl)
			services := &service.Services{CacheService: cacheService}
			test.mockBehavior(cacheService, test.inputRequest)
			handler := NewHandler(services)
			router := handler.Init()

			router.Get(redisRoute+getKeyRoute, handler.ReadValueByKey)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, redisRoute+slash+test.inputRequest, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
