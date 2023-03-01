package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	mockdb "github.com/adykaaa/online-notes/db/mock"
	db "github.com/adykaaa/online-notes/db/sqlc"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

type createNoteArgs db.CreateNoteParams

func (a *createNoteArgs) Matches(x interface{}) bool {
	reflectedValue := reflect.ValueOf(x).Elem()
	if a.Username != reflectedValue.FieldByName("Username").String() {
		return false
	}
	if a.Title != reflectedValue.FieldByName("Title").String() {
		return false
	}
	if a.Text.String != reflectedValue.FieldByName("Text").FieldByName("String").String() {
		return false
	}

	return true
}

func (a *createNoteArgs) String() string {
	return fmt.Sprintf("ID: %v, Title: %s, Username: %s, Text: %s", a.ID, a.Title, a.Username, a.Text.String)
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
			name: "note creation OK",

			body: testNote,

			validateJSON: func(t *testing.T, v *validator.Validate, note *models.Note) {
				err := v.Struct(note)
				require.NoError(t, err)
			},

			dbmockCreateNote: func(mockdb *mockdb.MockQuerier, note *models.Note) {
				args := createNoteArgs{
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
		{
			name: "returns bad request - wrongly formatted note param",

			body: &models.Note{
				ID:        uuid.New(),
				Title:     "",
				User:      "testuser1",
				Text:      "test1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},

			validateJSON: func(t *testing.T, v *validator.Validate, note *models.Note) {
				err := v.Struct(note)
				require.Error(t, err)
			},

			dbmockCreateNote: func(mockdb *mockdb.MockQuerier, note *models.Note) {
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "returns forbidden - duplicate note title",

			body: testNote,

			validateJSON: func(t *testing.T, v *validator.Validate, note *models.Note) {
				err := v.Struct(note)
				require.NoError(t, err)
			},

			dbmockCreateNote: func(mockdb *mockdb.MockQuerier, note *models.Note) {
				mockdb.EXPECT().CreateNote(gomock.Any(), gomock.Any()).Times(1).Return(note.ID, &pq.Error{Code: "23505"})
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "returns internal server error - db error",

			body: testNote,

			validateJSON: func(t *testing.T, v *validator.Validate, note *models.Note) {
				err := v.Struct(note)
				require.NoError(t, err)
			},

			dbmockCreateNote: func(mockdb *mockdb.MockQuerier, note *models.Note) {
				mockdb.EXPECT().CreateNote(gomock.Any(), gomock.Any()).Times(1).Return(note.ID, errors.New("internal error"))
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

func TestGetAllNotesFromUser(t *testing.T) {
	testCases := []struct {
		name                      string
		addQuery                  func(t *testing.T, r *http.Request)
		dbmockGetAllNotesFromUser func(mockdb *mockdb.MockQuerier)
		checkResponse             func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request)
	}{
		{
			name: "gettings notes from user OK",

			addQuery: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				q.Add("username", "testuser1")
				r.URL.RawQuery = q.Encode()
			},

			dbmockGetAllNotesFromUser: func(mockdb *mockdb.MockQuerier) {
				mockdb.EXPECT().GetAllNotesFromUser(gomock.Any(), "testuser1").Times(1).Return([]db.Note{}, nil)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			dbmock := mockdb.NewMockQuerier(ctrl)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/notes", nil)

			tc.addQuery(t, req)
			tc.dbmockGetAllNotesFromUser(dbmock)

			handler := GetAllNotesFromUser(dbmock)
			handler(rec, req)
			tc.checkResponse(t, rec, req)
		})
	}

}

func TestDeleteNote(t *testing.T) {

}

func TestUpdateNote(t *testing.T) {

}
