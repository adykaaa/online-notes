package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/adykaaa/online-notes/db/mock"
	db "github.com/adykaaa/online-notes/db/sqlc"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/adykaaa/online-notes/utils"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbmock := mockdb.NewMockQuerier(ctrl)
	jsonValidator := validator.New()

	testCases := []struct {
		name          string
		body          *models.User
		validateJSON  func(v *validator.Validate, user *models.User)
		hashPassword  func(user *models.User) string
		dbmock        func(mockdb *mockdb.MockQuerier, user *models.User, hashedPassword string)
		checkResponse func(recorder *httptest.ResponseRecorder, request *http.Request)
	}{
		{
			name: "User registration OK",

			body: &models.User{
				Username: "user1",
				Password: "password1",
				Email:    "user1@user.com",
			},

			validateJSON: func(v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.NoError(t, err)
			},

			hashPassword: func(user *models.User) string {
				hp, err := utils.HashPassword(user.Password)
				require.NoError(t, err)
				return hp
			},

			dbmock: func(mockdb *mockdb.MockQuerier, user *models.User, hashedPassword string) {
				args := db.RegisterUserParams{
					Username: user.Username,
					Password: hashedPassword,
					Email:    user.Email,
				}
				mockdb.EXPECT().RegisterUser(gomock.Any(), &args).Times(1).Return(args.Username, nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, recorder.Code, http.StatusCreated)
			},
		},
	}

	for c := range testCases {
		tc := testCases[c]

		tc.validateJSON(jsonValidator, tc.body)
		hp := tc.hashPassword(tc.body)
		tc.dbmock(dbmock, tc.body, hp)

		b, err := json.Marshal(tc.body)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))

		handler := RegisterUser(dbmock)
		handler(rec, req)

		tc.checkResponse(rec, req)
	}

}
