package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/adykaaa/online-notes/http/models"
)

func MockHandler(w http.ResponseWriter, r *http.Request) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Test server started!")
	}))
	defer svr.Close()

}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		body          *models.User
		handler       func(t *testing.T, w http.ResponseWriter, r *http.Request)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request)
	}{}

	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/notes/create", nil)
			w := httptest.NewRecorder()
		})
	}
}
