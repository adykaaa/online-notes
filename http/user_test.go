package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/adykaaa/online-notes/db/mock"
	db "github.com/adykaaa/online-notes/db/sqlc"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/adykaaa/online-notes/lib/password"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// needed for the custom mocker.
type regUserArgs db.RegisterUserParams

func (a *regUserArgs) Matches(x interface{}) bool {
	reflectedValue := reflect.ValueOf(x).Elem()
	if a.Username != reflectedValue.FieldByName("Username").String() {
		return false
	}
	if a.Email != reflectedValue.FieldByName("Email").String() {
		return false
	}
	err := password.Validate(reflectedValue.FieldByName("Password").String(), a.Password)
	if err != nil {
		return false
	}
	return true
}

func (a *regUserArgs) String() string {
	return fmt.Sprintf("Username: %s, Email: %s", a.Username, a.Email)
}

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbmock := mockdb.NewMockQuerier(ctrl)
	jsonValidator := validator.New()

	testCases := []struct {
		name          string
		body          *models.User
		validateJSON  func(v *validator.Validate, user *models.User)
		dbmock        func(mockdb *mockdb.MockQuerier, user *models.User)
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

			dbmock: func(mockdb *mockdb.MockQuerier, user *models.User) {
				args := regUserArgs{
					Username: user.Username,
					Password: user.Password,
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

		t.Run(tc.name, func(t *testing.T) {
			tc.validateJSON(jsonValidator, tc.body)
			_, err := password.Hash(tc.body.Password)
			require.NoError(t, err)

			tc.dbmock(dbmock, tc.body)

			b, err := json.Marshal(tc.body)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))

			handler := RegisterUser(dbmock)
			handler(rec, req)

			tc.checkResponse(rec, req)
		})
	}

}
