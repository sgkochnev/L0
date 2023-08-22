package get_test

import (
	"LO/internal/http-server/handlers/get"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"LO/internal/http-server/handlers/get/mocks"
)

func TestGetOrderHandler(t *testing.T) {
	cases := []struct {
		name       string
		orderUID   string
		fieldName  string
		expected   []byte
		statusCode int
		mockError  error
	}{
		{
			name:       "Post Ok",
			fieldName:  "order_uid",
			orderUID:   "order123",
			expected:   []byte("it's json data"),
			statusCode: http.StatusOK,
		},
		{
			name:       "Not Found 1",
			fieldName:  "order_uid",
			orderUID:   "order123",
			statusCode: http.StatusNotFound,
			mockError:  errors.New("Entry not found"),
		},
		{
			name:       "Not Found 2",
			fieldName:  "order_uid",
			statusCode: http.StatusNotFound,
			mockError:  errors.New("Entry not found"),
		},
	}

	for _, tc := range cases {

		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockGetter := mocks.NewGetter(t)

			mockGetter.
				On("Get", tc.orderUID).
				Return(tc.expected, tc.mockError)

			body := fmt.Sprintf(`{"%s": "%s"}`, tc.fieldName, tc.orderUID)

			handler := get.New(mockGetter)
			req, err := http.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte(body)))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			assert.Equal(t, rr.Code, tc.statusCode)

			assert.Equal(t, rr.Body.String(), string(tc.expected))
		})
	}
}

func TestGetOrderBadJson(t *testing.T) {

	cases := []struct {
		name       string
		body       string
		expected   []byte
		statusCode int
	}{
		{
			name:       "Bad request 1",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Bad request 2",
			body:       `{order123}`,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		tc := tc

		mockGetter := mocks.NewGetter(t)

		handler := get.New(mockGetter)
		req, err := http.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte(tc.body)))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, tc.statusCode)

		assert.Equal(t, rr.Body.String(), string(tc.expected))
	}
}

func TestGetOrderGet(t *testing.T) {

	cases := []struct {
		name       string
		orderUID   string
		expected   []byte
		statusCode int
		mockError  error
	}{
		{
			name:       "Get Ok",
			orderUID:   "order123",
			expected:   []byte("it's json data"),
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range cases {
		tc := tc

		mockGetter := mocks.NewGetter(t)

		mockGetter.
			On("Get", tc.orderUID).
			Return(tc.expected, tc.mockError)

		r := mux.NewRouter()
		r.Handle("/{order_uid}", get.New(mockGetter)).Methods(http.MethodGet)

		ts := httptest.NewServer(r)
		defer ts.Close()

		newURL := fmt.Sprintf("%s/%s", ts.URL, tc.orderUID)

		req, err := http.NewRequest(http.MethodGet, newURL, nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, tc.statusCode)

		assert.Equal(t, rr.Body.String(), string(tc.expected))
	}
}
