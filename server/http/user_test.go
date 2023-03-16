package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	db "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/adykaaa/online-notes/lib/password"
	"github.com/adykaaa/online-notes/note"
	mocksvc "github.com/adykaaa/online-notes/note/mock"
	auth "github.com/adykaaa/online-notes/server/http/auth"
	models "github.com/adykaaa/online-notes/server/http/models"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type regUserMatcher db.RegisterUserParams

func (m *regUserMatcher) Matches(x interface{}) bool {
	reflectedValue := reflect.ValueOf(x).Elem()
	if m.Username != reflectedValue.FieldByName("Username").String() {
		return false
	}
	if m.Email != reflectedValue.FieldByName("Email").String() {
		return false
	}
	err := password.Validate(reflectedValue.FieldByName("Password").String(), m.Password)
	if err != nil {
		return false
	}

	return true
}

func (m *regUserMatcher) String() string {
	return fmt.Sprintf("Username: %s, Email: %s", m.Username, m.Email)
}

func TestRegisterUser(t *testing.T) {
	jsonValidator := validator.New()

	testCases := []struct {
		name          string
		body          *models.User
		validateJSON  func(t *testing.T, v *validator.Validate, u *models.User)
		mockSvcCall   func(mocksvc *mocksvc.MockNoteService, u *models.User)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request)
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

			mockSvcCall: func(mocksvc *mocksvc.MockNoteService, u *models.User) {
				mocksvc.EXPECT().RegisterUser(gomock.Any(), &regUserMatcher{
					Username: u.Username,
					Password: u.Password,
					Email:    u.Email,
				}).Times(1).Return(u.Username, nil)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		}, {
			name: "returns bad request - short username",

			body: &models.User{
				Username: "u",
				Password: "password1",
				Email:    "user1@user.com",
			},

			validateJSON: func(t *testing.T, v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.Error(t, err)
			},

			mockSvcCall: func(mocksvc *mocksvc.MockNoteService, u *models.User) {
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "returns bad request - short password",

			body: &models.User{
				Username: "username1",
				Password: "pw1",
				Email:    "user1@user.com",
			},

			validateJSON: func(t *testing.T, v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.Error(t, err)
			},

			mockSvcCall: func(mocksvc *mocksvc.MockNoteService, u *models.User) {
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "returns bad request - malformatted email",

			body: &models.User{
				Username: "username1",
				Password: "password1",
				Email:    "wrongemail@",
			},

			validateJSON: func(t *testing.T, v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.Error(t, err)
			},

			mockSvcCall: func(mocksvc *mocksvc.MockNoteService, u *models.User) {
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "returns forbidden - duplicate username",

			body: &models.User{
				Username: "username1",
				Password: "password1",
				Email:    "user1@user.com",
			},

			validateJSON: func(t *testing.T, v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.NoError(t, err)
			},

			mockSvcCall: func(mocksvc *mocksvc.MockNoteService, u *models.User) {
				mocksvc.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(1).Return("", note.ErrAlreadyExists)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "returns internal error - DB failure",

			body: &models.User{
				Username: "username1",
				Password: "password1",
				Email:    "user1@user.com",
			},

			validateJSON: func(t *testing.T, v *validator.Validate, user *models.User) {
				err := v.Struct(user)
				require.NoError(t, err)
			},

			mockSvcCall: func(mocksvc *mocksvc.MockNoteService, u *models.User) {
				mocksvc.EXPECT().RegisterUser(gomock.Any(), gomock.Any()).Times(1).Return("", note.ErrDBInternal)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mocksvc := mocksvc.NewMockNoteService(ctrl)

			tc.validateJSON(t, jsonValidator, tc.body)
			tc.mockSvcCall(mocksvc, tc.body)

			b, err := json.Marshal(tc.body)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))

			handler := RegisterUser(mocksvc)
			handler(rec, req)
			tc.checkResponse(t, rec, req)
		})
	}
}
func TestLoginUser(t *testing.T) {

	type loginUserReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	tm := &auth.MockTokenManager{}

	testCases := []struct {
		name             string
		body             *loginUserReq
		mockSvcCall      func(t *testing.T, mocksvc *mocksvc.MockNoteService, user *loginUserReq) string
		validatePassword func(t *testing.T, user *loginUserReq, dbUserPassword string)
		createToken      func(t *testing.T, tm auth.TokenManager, user *loginUserReq, duration time.Duration) string
		checkResponse    func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request, token string)
	}{
		{
			name: "user login OK",

			body: &loginUserReq{
				Username: "user1",
				Password: "password1",
			},

			mockSvcCall: func(t *testing.T, mocksvc *mocksvc.MockNoteService, user *loginUserReq) string {
				hashedPassword, err := password.Hash(user.Password)
				require.NoError(t, err)

				dbuser := db.User{
					Username: user.Username,
					Password: hashedPassword,
				}

				mocksvc.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(dbuser, nil)
				return hashedPassword
			},

			validatePassword: func(t *testing.T, user *loginUserReq, dbHashedPassword string) {
				err := password.Validate(dbHashedPassword, user.Password)
				require.NoError(t, err)
			},

			createToken: func(t *testing.T, tm auth.TokenManager, user *loginUserReq, duration time.Duration) string {
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
			name: "returns not found - user not registered",

			body: &loginUserReq{
				Username: "user1",
				Password: "password1",
			},

			mockSvcCall: func(t *testing.T, mocksvc *mocksvc.MockNoteService, user *loginUserReq) string {
				hashedPassword, err := password.Hash(user.Password)
				require.NoError(t, err)

				mocksvc.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(db.User{}, note.ErrNotFound)
				return hashedPassword
			},

			validatePassword: func(t *testing.T, user *loginUserReq, dbHashedPassword string) {
			},

			createToken: func(t *testing.T, tm auth.TokenManager, user *loginUserReq, duration time.Duration) string {
				return ""
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request, token string) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "returns unauthorized - wrong password",

			body: &loginUserReq{
				Username: "user1",
				Password: "password1",
			},

			mockSvcCall: func(t *testing.T, mocksvc *mocksvc.MockNoteService, user *loginUserReq) string {
				hashedPassword, err := password.Hash("wrongpassword")
				require.NoError(t, err)

				dbuser := db.User{
					Username: user.Username,
					Password: hashedPassword,
				}

				mocksvc.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(dbuser, nil)
				return hashedPassword
			},

			validatePassword: func(t *testing.T, user *loginUserReq, dbHashedPassword string) {
				err := password.Validate(dbHashedPassword, user.Password)
				require.Error(t, err)
			},

			createToken: func(t *testing.T, tm auth.TokenManager, user *loginUserReq, duration time.Duration) string {
				return ""
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request, token string) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "returns internal server error - DB error",

			body: &loginUserReq{
				Username: "user1",
				Password: "password1",
			},

			mockSvcCall: func(t *testing.T, mocksvc *mocksvc.MockNoteService, user *loginUserReq) string {
				mocksvc.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(db.User{}, note.ErrDBInternal)
				return ""
			},

			validatePassword: func(t *testing.T, user *loginUserReq, dbHashedPassword string) {
			},

			createToken: func(t *testing.T, tm auth.TokenManager, user *loginUserReq, duration time.Duration) string {
				return ""
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request, token string) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}
	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mocksvc := mocksvc.NewMockNoteService(ctrl)

			b, err := json.Marshal(tc.body)
			require.NoError(t, err)

			pw := tc.mockSvcCall(t, mocksvc, tc.body)
			tc.validatePassword(t, tc.body, pw)
			token := tc.createToken(t, tm, tc.body, 300)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))

			handler := LoginUser(mocksvc, tm, 300)
			handler(rec, req)

			tc.checkResponse(t, rec, req, token)
		})
	}
}
