package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	db "github.com/adykaaa/online-notes/db/sqlc"
	note "github.com/adykaaa/online-notes/note"
	mocksvc "github.com/adykaaa/online-notes/note/mock"
	models "github.com/adykaaa/online-notes/server/http/models"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateNote(t *testing.T) {
	jsonValidator := validator.New()

	testNote := &models.Note{
		ID:        uuid.New(),
		Title:     "testtitle",
		User:      "testuser1",
		Text:      "testtext",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testCases := []struct {
		name          string
		body          *models.Note
		validateJSON  func(t *testing.T, v *validator.Validate, note *models.Note)
		mockSvcCall   func(mocksvc *mocksvc.MockNoteService, note *models.Note)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request)
	}{
		{
			name: "note creation OK",

			body: testNote,

			validateJSON: func(t *testing.T, v *validator.Validate, note *models.Note) {
				err := v.Struct(note)
				require.NoError(t, err)
			},

			mockSvcCall: func(mocksvc *mocksvc.MockNoteService, n *models.Note) {
				mocksvc.EXPECT().CreateNote(gomock.Any(), n.Title, n.User, n.Text).Times(1).Return(n.ID, nil)
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

			mockSvcCall: func(mocksvc *mocksvc.MockNoteService, n *models.Note) {
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "returns forbidden - duplicate note title",

			body: testNote,

			validateJSON: func(t *testing.T, v *validator.Validate, n *models.Note) {
				err := v.Struct(n)
				require.NoError(t, err)
			},

			mockSvcCall: func(mocksvc *mocksvc.MockNoteService, n *models.Note) {
				mocksvc.EXPECT().CreateNote(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(n.ID, note.ErrAlreadyExists)
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

			mockSvcCall: func(mocksvc *mocksvc.MockNoteService, n *models.Note) {
				mocksvc.EXPECT().CreateNote(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(n.ID, note.ErrDBInternal)
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
			req := httptest.NewRequest(http.MethodPost, "/notes/create", bytes.NewReader(b))

			handler := CreateNote(mocksvc)
			handler(rec, req)
			tc.checkResponse(t, rec, req)
		})
	}

}

func TestGetAllNotesFromUser(t *testing.T) {
	testCases := []struct {
		name          string
		addQuery      func(t *testing.T, r *http.Request)
		mockSvcCall   func(svcmock *mocksvc.MockNoteService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request)
	}{
		{
			name: "gettings notes from user OK",

			addQuery: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				q.Add("username", "testuser1")
				r.URL.RawQuery = q.Encode()
			},

			mockSvcCall: func(mocksvc *mocksvc.MockNoteService) {
				mocksvc.EXPECT().GetAllNotesFromUser(gomock.Any(), "testuser1").Times(1).Return([]db.Note{}, nil)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "returns bad request - missing url param",

			addQuery: func(t *testing.T, r *http.Request) {
			},

			mockSvcCall: func(mocksvc *mocksvc.MockNoteService) {
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "returns internal server error - db error",

			addQuery: func(t *testing.T, r *http.Request) {
				q := r.URL.Query()
				q.Add("username", "testuser1")
				r.URL.RawQuery = q.Encode()
			},

			mockSvcCall: func(mocksvc *mocksvc.MockNoteService) {
				mocksvc.EXPECT().GetAllNotesFromUser(gomock.Any(), "testuser1").Times(1).Return(nil, note.ErrDBInternal)
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

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/notes", nil)

			tc.addQuery(t, req)
			tc.mockSvcCall(mocksvc)

			handler := GetAllNotesFromUser(mocksvc)
			handler(rec, req)
			tc.checkResponse(t, rec, req)
		})
	}

}

func TestDeleteNote(t *testing.T) {

	id := uuid.New()
	testCases := []struct {
		name          string
		path          string
		mockSvcCall   func(svcmock *mocksvc.MockNoteService)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request)
	}{
		{
			name: "deleting note OK",
			path: id.String(),
			mockSvcCall: func(mocksvc *mocksvc.MockNoteService) {
				mocksvc.EXPECT().DeleteNote(gomock.Any(), id).Times(1).Return(id, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "deleting note returns bad request - malformed uuid",
			path: "invalid",
			mockSvcCall: func(mocksvc *mocksvc.MockNoteService) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "deleting note returns bad request - malformed uuid",
			path: id.String(),
			mockSvcCall: func(mocksvc *mocksvc.MockNoteService) {
				mocksvc.EXPECT().DeleteNote(gomock.Any(), id).Times(1).Return(uuid.Nil, note.ErrDBInternal)
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

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/notes/"+tc.path, nil)

			tc.mockSvcCall(mocksvc)

			handler := DeleteNote(mocksvc)
			handler(rec, req)
			tc.checkResponse(t, rec, req)
		})
	}
}

func TestUpdateNote(t *testing.T) {

	id := uuid.New()
	updateRequest := struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}{}

	testCases := []struct {
		name          string
		path          string
		body          *models.Note
		validateJSON  func(t *testing.T, v *validator.Validate, note *models.Note)
		mockSvcCall   func(mocksvc *mocksvc.MockNoteService, note *models.Note)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request)
	}{
		{
			name: "updating note OK",
			path: "invalid",
			mockSvcCall: func(mocksvc *mocksvc.MockNoteService, note *models.Note) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/notes/"+tc.path, nil)

			tc.mockSvcCall(mocksvc)

			handler := UpdateNote(mocksvc)
			handler(rec, req)
			tc.checkResponse(t, rec, req)
		})
	}
}
