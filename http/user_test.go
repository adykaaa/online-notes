package http

import (
	"context"
	"encoding/json"
	"net/http"
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
		checkResponse func(user *models.User, path string)
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
				mockdb.EXPECT().RegisterUser(ctx, &args).Times(1).Return(args.Username, nil)
			},

			checkResponse: func(user *models.User, path string) {

				b, err := json.Marshal(user)
				require.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, "/register", b)
				res := httptest.NewRecorder()
				RegisterUser(req, res)

			},
		},
	}

	for c := range testCases {
		tc := testCases[c]

		tc.validateJSON(jsonValidator, tc.body)
		tc.dbmock(dbmock, tc.body)

	}

}
