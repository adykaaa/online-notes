package http

import (
	"context"
	"net/http/httptest"
	"testing"

	mockdb "github.com/adykaaa/online-notes/db/mock"
	db "github.com/adykaaa/online-notes/db/sqlc"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbmock := mockdb.NewMockQuerier(ctrl)
	jsonValidator := validator.New()
	ctx := context.Background()

	testCases := []struct {
		name          string
		body          *models.User
		validateJSON  func(v *validator.Validate, body *models.User)
		dbmock        func(mockdb *mockdb.MockQuerier, body *models.User)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "User registration OK",

			body: &models.User{
				Username: "user1",
				Password: "password1",
				Email:    "user1@user.com",
			},

			validateJSON: func(v *validator.Validate, body *models.User) {
				err := v.Struct(body)
				require.NoError(t, err)
			},

			dbmock: func(mockdb *mockdb.MockQuerier, body *models.User) {
				args := db.RegisterUserParams{
					Username: body.Username,
					Password: body.Password,
					Email:    body.Email,
				}

				mockdb.EXPECT().RegisterUser(ctx, &args).Return(args.Username, nil)
				user, err := mockdb.RegisterUser(ctx, &args)
				require.NoError(t, err)
				require.Equal(t, user, args.Username)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
			},
		},
	}
	for c := range testCases {

		testCases[c].validateJSON(jsonValidator, testCases[c].body)

	}

}
