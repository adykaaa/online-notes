package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	httplib "github.com/adykaaa/online-notes/lib/http"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	l := zerolog.New(io.Discard)

	testCases := []struct {
		name            string
		newMockTokenMgr func() *MockTokenManager
		setCookie       func(w http.ResponseWriter)
		checkResponse   func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request)
	}{
		{
			name: "auth OK",

			newMockTokenMgr: func() *MockTokenManager {
				return &MockTokenManager{
					ReturnInvalidToken: false,
					ReturnExpiredToken: false,
				}
			},

			setCookie: func(w http.ResponseWriter) {
				httplib.SetCookie(w, "paseto", "testtoken", time.Now().Add(30*time.Minute))
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {
			tm := tc.newMockTokenMgr()

			r := chi.NewRouter()
			r.Use(AuthMiddleware(tm, &l))

			r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
				httplib.JSON(w, "msg from test handler", http.StatusOK)
			})

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			tc.setCookie(rec)

			r.ServeHTTP(rec, req)
			tc.checkResponse(t, rec, req)

		})
	}
}
