package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	httplib "github.com/adykaaa/online-notes/lib/http"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
	l := zerolog.New(io.Discard)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httplib.JSON(w, msg{"success": "test handler success"}, http.StatusOK)
	})

	ts := httptest.NewServer(testHandler)
	defer ts.Close()

	tm := &MockTokenManager{}
	testCases := []struct {
		name          string
		testHandler   func(t *testing.T, h func(w http.ResponseWriter, r *http.Request))
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request)
	}{
		{
			name: "auth OK",
			testHandler: func(t *testing.T, h func(w http.ResponseWriter, r *http.Request)) {
				httplib.SetCookie(w, "paseto", "szevasz", 360)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
	}

	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPost, ts.URL, nil)
			rec := httptest.NewRecorder()
			asd := AuthMiddleware(tm, "testkey", &l)(testHandler)
			asd.ServeHTTP(rec, req)
			tc.checkResponse(t, rec, req)

		})
	}
}
