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
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type createNoteArgs db.CreateNoteParams

func (a *createNoteArgs) Matches(x interface{}) bool {
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

func (a *createNoteArgs) String() string {
	return fmt.Sprintf("Username: %s, Email: %s", a.Username, a.Email)
}

func TestCreateNote(t *testing.T) {
	jsonValidator := validator.New()

	testNote := &models.Note{
		ID:        uuid.New(),
		Title:     "test1",
		User:      "testuser1",
		Text:      "test1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testCases := []struct {
		name             string
		body             *models.Note
		validateJSON     func(t *testing.T, v *validator.Validate, note *models.Note)
		dbmockCreateNote func(mockdb *mockdb.MockQuerier, note *models.Note)
		checkResponse    func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request)
	}{
		{
			name: "Note creation OK",

			body: testNote,

			validateJSON: func(t *testing.T, v *validator.Validate, note *models.Note) {
				err := v.Struct(note)
				require.NoError(t, err)
			},

			dbmockCreateNote: func(mockdb *mockdb.MockQuerier, note *models.Note) {
				args := db.CreateNoteParams{
					ID:        testNote.ID,
					Title:     testNote.Title,
					Username:  testNote.User,
					Text:      sql.NullString{String: testNote.Text, Valid: true},
					CreatedAt: testNote.CreatedAt,
					UpdatedAt: testNote.UpdatedAt,
				}

				mockdb.EXPECT().CreateNote(gomock.Any(), &args).Times(1).Return(args.ID, nil)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}

	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			dbmock := mockdb.NewMockQuerier(ctrl)

			tc.validateJSON(t, jsonValidator, tc.body)
			tc.dbmockCreateNote(dbmock, tc.body)

			b, err := json.Marshal(tc.body)
			require.NoError(t, err)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/notes/create", bytes.NewReader(b))

			handler := CreateNote(dbmock)
			handler(rec, req)
			tc.checkResponse(t, rec, req)
		})
	}

}
