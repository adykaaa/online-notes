package http

import (
	"context"
	"net/http/httptest"
	"testing"

	mockdb "github.com/adykaaa/online-notes/db/mock"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/golang/mock/gomock"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockdb := mockdb.NewMockQuerier(ctrl)
	ctx := context.Background()

	tc := []struct {
		name          string
		body          models.User
		dbmock        func(db *mockdb.MockQuerier)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{}
}
