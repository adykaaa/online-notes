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

type MockTokenManager struct{}

func (m *MockTokenManager) CreateToken(username string, duration time.Duration) (string, *PasetoPayload, error) {
	return "testtoken",
		&PasetoPayload{},
		nil
}

func (m *MockTokenManager) VerifyToken(token string) (*PasetoPayload, error) {
	return &PasetoPayload{},
		nil
}

func TestRegisterUser(t *testing.T) {
	jsonValidator := validator.New()

	testCases := []struct {
		name             string
		body             *models.User
		validateJSON     func(t *testing.T, v *validator.Validate, user *models.User)
		dbmockCreateUser func(mockdb *mockdb.MockQuerier, user *models.User)
		checkResponse    func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request)
	}{
		{
			name: "user registration OK",

			body: &models.User{
				Username: "user1",
				Password: "password1",
				Email:    "user1@user.com",
			},

			validateJSON: func(t *testing.T, v *validator.Validate, user *models.User) {
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

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, recorder.Code, http.StatusCreated)
			},
		}, {
			name: "returns bad request because of short username",

			body: &models.User{
				Username: "u1",
				Password: "password1",
				Email:    "user1@user.com",
			},

			validateJSON: func(t *testing.T, v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.Error(t, err)
			},

			dbmockCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
				mockdb.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(0).Return("", nil)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
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

			validateJSON: func(t *testing.T, v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.Error(t, err)
			},

			dbmockCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
				mockdb.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(0).Return("", nil)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
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

			validateJSON: func(t *testing.T, v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.Error(t, err)
			},

			dbmockCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
				mockdb.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(0).Return("", nil)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
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

			validateJSON: func(t *testing.T, v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.NoError(t, err)
			},

			dbmockCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
				mockdb.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(1).Return("", &pq.Error{Code: "23505"})
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
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

			validateJSON: func(t *testing.T, v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.NoError(t, err)
			},

			dbmockCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
				mockdb.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(1).Return("", sql.ErrConnDone)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
			},
		},
	}

	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			dbmock := mockdb.NewMockQuerier(ctrl)

			tc.validateJSON(t, jsonValidator, tc.body)
			tc.dbmockCreateUser(dbmock, tc.body)

			b, err := json.Marshal(tc.body)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))

			handler := RegisterUser(dbmock)
			handler(rec, req)
			tc.checkResponse(t, rec, req)
		})
	}
}
func TestLoginUser(t *testing.T) {
	jsonValidator := validator.New()
	tm := &MockTokenManager{}

	testCases := []struct {
		name             string
		body             *LoginUserReq
		validateJSON     func(t *testing.T, v *validator.Validate, user *LoginUserReq)
		dbmockGetUser    func(t *testing.T, mockdb *mockdb.MockQuerier, user *LoginUserReq) string
		validatePassword func(t *testing.T, user *LoginUserReq, dbUserPassword string)
		createToken      func(t *testing.T, tm TokenManager, user *LoginUserReq, duration time.Duration) string
		checkResponse    func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request, token string)
	}{
		{
			name: "user login OK",

			body: &LoginUserReq{
				Username: "user1",
				Password: "password1",
			},

			validateJSON: func(t *testing.T, v *validator.Validate, user *LoginUserReq) {
				err := v.Struct(user)
				require.NoError(t, err)
			},

			dbmockGetUser: func(t *testing.T, mockdb *mockdb.MockQuerier, user *LoginUserReq) string {
				hashedPassword, err := password.Hash(user.Password)
				require.NoError(t, err)

				dbuser := db.User{
					Username: user.Username,
					Password: hashedPassword,
				}

				mockdb.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(dbuser, nil)
				return hashedPassword
			},

			validatePassword: func(t *testing.T, user *LoginUserReq, dbHashedPassword string) {
				err := password.Validate(dbHashedPassword, user.Password)
				require.NoError(t, err)
			},

			createToken: func(t *testing.T, tm TokenManager, user *LoginUserReq, duration time.Duration) string {
				token, _, err := tm.CreateToken(user.Username, duration)
				require.NoError(t, err)

				return token
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request, token string) {
				require.Equal(t, recorder.Result().Cookies()[0].Name, "paseto")
				require.Equal(t, recorder.Result().Cookies()[0].Value, token)
				require.Equal(t, recorder.Result().Cookies()[0].HttpOnly, true)
				require.Equal(t, recorder.Result().Cookies()[0].Secure, true)
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "returns internal server error because of JSON decode error",

			body: &LoginUserReq{
				Username: "user1",
				Password: "password1",
			},

			validateJSON: func(t *testing.T, v *validator.Validate, user *LoginUserReq) {
				err := v.Struct(user)
				require.NoError(t, err)
			},

			dbmockGetUser: func(t *testing.T, mockdb *mockdb.MockQuerier, user *LoginUserReq) string {
				hashedPassword, err := password.Hash(user.Password)
				require.NoError(t, err)

				dbuser := db.User{
					Username: user.Username,
					Password: hashedPassword,
				}

				mockdb.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(dbuser, nil)
				return hashedPassword
			},

			validatePassword: func(t *testing.T, user *LoginUserReq, dbHashedPassword string) {
				err := password.Validate(dbHashedPassword, user.Password)
				require.NoError(t, err)
			},

			createToken: func(t *testing.T, tm TokenManager, user *LoginUserReq, duration time.Duration) string {
				token, _, err := tm.CreateToken(user.Username, 300)
				require.NoError(t, err)
				return token
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request, token string) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			dbmock := mockdb.NewMockQuerier(ctrl)

			tc.validateJSON(t, jsonValidator, tc.body)
			pw := tc.dbmockGetUser(t, dbmock, tc.body)

			tc.validatePassword(t, tc.body, pw)
			token := tc.createToken(t, tm, tc.body, 300)

			b, err := json.Marshal(tc.body)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))

			handler := LoginUser(dbmock, tm, 300)
			handler(rec, req)
			tc.checkResponse(t, rec, req, token)
		})
	}
}
