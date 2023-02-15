package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	mockdb "github.com/adykaaa/online-notes/db/mock"
	db "github.com/adykaaa/online-notes/db/sqlc"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/adykaaa/online-notes/lib/password"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
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
		name             string
		body             *models.User
		validateJSON     func(v *validator.Validate, user *models.User)
		dbmockCreateUser func(mockdb *mockdb.MockQuerier, user *models.User)
		checkResponse    func(recorder *httptest.ResponseRecorder, request *http.Request)
	}{
		{
			name: "user registration OK",

			body: &models.User{
				Username: "user1",
				Password: "password1",
				Email:    "user1@user.com",
			},

			validateJSON: func(v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.NoError(t, err)
			},

			dbmockCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
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
		}, {
			name: "returns bad request because of short username",

			body: &models.User{
				Username: "u1",
				Password: "password1",
				Email:    "user1@user.com",
			},

			validateJSON: func(v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.Error(t, err)
			},

			dbmockCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
				mockdb.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(0).Return("", nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "returns bad request because of short password",

			body: &models.User{
				Username: "username1",
				Password: "pw1",
				Email:    "user1@user.com",
			},

			validateJSON: func(v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.Error(t, err)
			},

			dbmockCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
				mockdb.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(0).Return("", nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "returns bad request because malformatted email",

			body: &models.User{
				Username: "username1",
				Password: "password1",
				Email:    "wrongemail@",
			},

			validateJSON: func(v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.Error(t, err)
			},

			dbmockCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
				mockdb.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(0).Return("", nil)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "fails because of duplicate username",

			body: &models.User{
				Username: "username1",
				Password: "password1",
				Email:    "user1@user.com",
			},

			validateJSON: func(v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.NoError(t, err)
			},

			dbmockCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
				mockdb.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(1).Return("", &pq.Error{Code: "23505"})
			},

			checkResponse: func(recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, recorder.Code, http.StatusForbidden)
			},
		},
		{
			name: "returns internal error because of DB failure",

			body: &models.User{
				Username: "username1",
				Password: "password1",
				Email:    "user1@user.com",
			},

			validateJSON: func(v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.NoError(t, err)
			},

			dbmockCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
				mockdb.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(1).Return("", sql.ErrConnDone)
			},

			checkResponse: func(recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
			},
		},
	}

	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {
			tc.validateJSON(jsonValidator, tc.body)
			tc.dbmockCreateUser(dbmock, tc.body)

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
func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	dbmock := mockdb.NewMockQuerier(ctrl)
	jsonValidator := validator.New()

	testCases := []struct {
		name             string
		body             *models.User
		validateJSON     func(v *validator.Validate, user *models.User)
		dbmockGetUser    func(mockdb *mockdb.MockQuerier, user *models.User) string
		validatePassword func(user *models.User, dbUserPassword string)
		createToken      func(user *models.User, duration time.Duration)
		checkResponse    func(recorder *httptest.ResponseRecorder, request *http.Request)
	}{
		{
			name: "user login OK",
			//recorder.Result().Cookies()[0].Value
			body: &models.User{
				Username: "user1",
				Password: "password1",
				Email:    "user1@user.com",
			},

			validateJSON: func(v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.NoError(t, err)
			},

			dbmockGetUser: func(mockdb *mockdb.MockQuerier, user *models.User) string {
				dbuser := db.User{
					Username: user.Username,
					Password: user.Password,
					Email:    user.Email,
				}
				mockdb.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(dbuser, nil)
				return dbuser.Password
			},

			validatePassword: func(user *models.User, dbUserPassword string) {
				err := password.Validate(user.Password, dbUserPassword)
				require.NoError(t, err)
			},
		},
	}
}
