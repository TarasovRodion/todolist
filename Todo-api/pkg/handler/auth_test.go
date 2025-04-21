package handler

import (
	todo "Todo-api"
	"Todo-api/pkg/service"
	mock_service "Todo-api/pkg/service/mocks"
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_singUp(t *testing.T) {
	type mockbehavior func(s *mock_service.MockAuthorization, user todo.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            todo.User
		mockBehavior         mockbehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"test","username":"test","password":"qwerty"}`,
			inputUser: todo.User{Name: "test", Username: "test", Password: "qwerty"},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Empty Fields",
			inputBody: `{"name":"","username":"","password":""}`,
			inputUser: todo.User{Name: "", Username: "", Password: ""},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(0, nil)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"empty fields"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			//Test Server
			r := gin.Default()
			r.POST("/sign-up", handler.singUp)

			//Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))

			//perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
