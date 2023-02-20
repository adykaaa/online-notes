package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/adykaaa/online-notes/http/models"
	"github.com/rs/zerolog"
)

func TestAuthMiddleware(t *testing.T) {
	l := zerolog.New(io.Discard)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from test handler!"))
	})
	ts := httptest.NewServer(testHandler)
	defer ts.Close()

	tm := MockTokenManager{}
	testCases := []struct {
		name          string
		body          *models.User
		handler       func(t *testing.T, w http.ResponseWriter, r *http.Request)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request)
	}{}

	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, ts.URL, nil)
			rec := httptest.NewRecorder()
			asd := AuthMiddleware(tm, "testkey", &l)(testHandler)
			asd.ServeHTTP(rec, req)

		})
	}
}
