package note

import (
	"context"
	"testing"

	mockdb "github.com/adykaaa/online-notes/db/mock"
	db "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/adykaaa/online-notes/lib/random"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	user := &db.User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@user.com",
	}

	testCases := []struct {
		name              string
		user              *db.User
		mockdbCreateUser  func(mockdb *mockdb.MockQuerier, args *db.RegisterUserParams)
		checkReturnValues func(t *testing.T, user *db.User, username string, err error)
	}{
		{
			name: "user registration OK",
			user: user,
			mockdbCreateUser: func(mockdb *mockdb.MockQuerier, args *db.RegisterUserParams) {
				mockdb.EXPECT().RegisterUser(gomock.Any(), args).Times(1).Return(args.Username, nil)
			},
			checkReturnValues: func(t *testing.T, user *db.User, username string, err error) {
				require.Equal(t, username, user.Username)
				require.Nil(t, err)
			},
		},
		{
			name: "user registration returns ErrUserAlreadyExists",
			user: user,
			mockdbCreateUser: func(mockdb *mockdb.MockQuerier, args *db.RegisterUserParams) {
				mockdb.EXPECT().RegisterUser(gomock.Any(), args).Times(1).Return("", ErrUserAlreadyExists)
			},
			checkReturnValues: func(t *testing.T, user *db.User, username string, err error) {
				require.ErrorIs(t, err, ErrUserAlreadyExists)
				require.Empty(t, username)
			},
		},
		{
			name: "user registration returns ErrDBInternal",
			user: user,
			mockdbCreateUser: func(mockdb *mockdb.MockQuerier, args *db.RegisterUserParams) {
				mockdb.EXPECT().RegisterUser(gomock.Any(), args).Times(1).Return("", ErrDBInternal)
			},
			checkReturnValues: func(t *testing.T, user *db.User, username string, err error) {
				require.ErrorIs(t, err, ErrDBInternal)
				require.Empty(t, username)
			},
		},
	}

	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockdb := mockdb.NewMockQuerier(ctrl)
			ns := NewService(mockdb)
			args := db.RegisterUserParams{
				Username: tc.user.Username,
				Password: tc.user.Password,
				Email:    tc.user.Email,
			}

			tc.mockdbCreateUser(mockdb, &args)

			u, err := ns.q.RegisterUser(context.Background(), &args)
			tc.checkReturnValues(t, tc.user, u, err)
		})
	}
}

func TestGetUser(t *testing.T) {
	const username = "user1"

	testCases := []struct {
		name              string
		username          string
		mockdbGetUser     func(mockdb *mockdb.MockQuerier)
		checkReturnValues func(t *testing.T, user db.User, err error)
	}{
		{
			name:     "getting user OK",
			username: username,
			mockdbGetUser: func(mockdb *mockdb.MockQuerier) {
				mockdb.EXPECT().GetUser(gomock.Any(), username).Times(1).Return(db.User{Username: username}, nil)
			},
			checkReturnValues: func(t *testing.T, user db.User, err error) {
				require.Equal(t, username, user.Username)
				require.Nil(t, err)
			},
		},
		{
			name:     "getting user returns ErrUserNotFound",
			username: username,
			mockdbGetUser: func(mockdb *mockdb.MockQuerier) {
				mockdb.EXPECT().GetUser(gomock.Any(), username).Times(1).Return(db.User{}, ErrUserNotFound)
			},
			checkReturnValues: func(t *testing.T, user db.User, err error) {
				require.Empty(t, user.Email)
				require.ErrorIs(t, err, ErrUserNotFound)
			},
		},
		{
			name:     "getting user returns ErrDBInternal",
			username: username,
			mockdbGetUser: func(mockdb *mockdb.MockQuerier) {
				mockdb.EXPECT().GetUser(gomock.Any(), username).Times(1).Return(db.User{}, ErrDBInternal)
			},
			checkReturnValues: func(t *testing.T, user db.User, err error) {
				require.Empty(t, user.Email)
				require.ErrorIs(t, err, ErrDBInternal)
			},
		},
	}

	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockdb := mockdb.NewMockQuerier(ctrl)
			ns := NewService(mockdb)

			tc.mockdbGetUser(mockdb)
			user, err := ns.q.GetUser(context.Background(), tc.username)
			tc.checkReturnValues(t, user, err)
		})
	}
}

func TestCreateNote(t *testing.T) {
	note := random.NewDBNote(uuid.New())

	testCases := []struct {
		name              string
		note              *db.Note
		mockdbCreateNote  func(mockdb *mockdb.MockQuerier, args *db.CreateNoteParams)
		checkReturnValues func(t *testing.T, note *db.Note, id uuid.UUID, err error)
	}{
		{
			name: "creating note OK",
			note: note,
			mockdbCreateNote: func(mockdb *mockdb.MockQuerier, args *db.CreateNoteParams) {
				mockdb.EXPECT().CreateNote(gomock.Any(), args).Times(1).Return(note.ID, nil)
			},
			checkReturnValues: func(t *testing.T, note *db.Note, id uuid.UUID, err error) {
				require.Equal(t, note.ID, id)
				require.Nil(t, err)
			},
		},
		{
			name: "creating note returns ErrAlreadyExist",
			note: note,
			mockdbCreateNote: func(mockdb *mockdb.MockQuerier, args *db.CreateNoteParams) {
				mockdb.EXPECT().CreateNote(gomock.Any(), args).Times(1).Return(uuid.Nil, ErrAlreadyExists)
			},
			checkReturnValues: func(t *testing.T, note *db.Note, id uuid.UUID, err error) {
				require.Equal(t, id, uuid.Nil)
				require.ErrorIs(t, err, ErrAlreadyExists)
			},
		},
		{
			name: "creating note returns ErrDBInternal",
			note: note,
			mockdbCreateNote: func(mockdb *mockdb.MockQuerier, args *db.CreateNoteParams) {
				mockdb.EXPECT().CreateNote(gomock.Any(), args).Times(1).Return(uuid.Nil, ErrDBInternal)
			},
			checkReturnValues: func(t *testing.T, note *db.Note, id uuid.UUID, err error) {
				require.Equal(t, id, uuid.Nil)
				require.ErrorIs(t, err, ErrDBInternal)
			},
		},
	}

	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockdb := mockdb.NewMockQuerier(ctrl)
			ns := NewService(mockdb)
			args := db.CreateNoteParams{
				ID:       tc.note.ID,
				Title:    tc.note.Title,
				Username: tc.note.Username,
				Text:     tc.note.Text,
			}

			tc.mockdbCreateNote(mockdb, &args)
			id, err := ns.q.CreateNote(context.Background(), &args)
			tc.checkReturnValues(t, tc.note, id, err)
		})
	}
}

func TestGetAllNotesFromUSer(t *testing.T) {
	const username = "user1"

	testCases := []struct {
		name                   string
		username               string
		mockdbGetNotesFromUser func(mockdb *mockdb.MockQuerier)
		checkReturnValues      func(t *testing.T, notes []db.Note, err error)
	}{
		{
			name:     "getting notes from user OK",
			username: username,
			mockdbGetNotesFromUser: func(mockdb *mockdb.MockQuerier) {
				mockdb.EXPECT().GetAllNotesFromUser(gomock.Any(), username).Times(1).Return([]db.Note{}, nil)
			},
			checkReturnValues: func(t *testing.T, notes []db.Note, err error) {
				require.Nil(t, err)
			},
		},
		{
			name:     "getting notes from user returns ErrDBInternal",
			username: username,
			mockdbGetNotesFromUser: func(mockdb *mockdb.MockQuerier) {
				mockdb.EXPECT().GetAllNotesFromUser(gomock.Any(), username).Times(1).Return([]db.Note{}, ErrDBInternal)
			},
			checkReturnValues: func(t *testing.T, notes []db.Note, err error) {
				require.ErrorIs(t, err, ErrDBInternal)
			},
		},
	}
	for c := range testCases {
		tc := testCases[c]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockdb := mockdb.NewMockQuerier(ctrl)
			ns := NewService(mockdb)

			tc.mockdbGetNotesFromUser(mockdb)
			notes, err := ns.q.GetAllNotesFromUser(context.Background(), tc.username)
			tc.checkReturnValues(t, notes, err)
		})
	}
}

func TestUpdateNote(t *testing.T) {
	note := random.NewDBNote(uuid.New())

	testCases := []struct {
		name              string
		note              *db.Note
		mockdbCreateNote  func(mockdb *mockdb.MockQuerier, args *db.CreateNoteParams)
		checkReturnValues func(t *testing.T, note *db.Note, id uuid.UUID, err error)
	}{
		{
			name: "creating note OK",
			note: note,
			mockdbCreateNote: func(mockdb *mockdb.MockQuerier, args *db.CreateNoteParams) {
				mockdb.EXPECT().CreateNote(gomock.Any(), args).Times(1).Return(note.ID, nil)
			},
			checkReturnValues: func(t *testing.T, note *db.Note, id uuid.UUID, err error) {
				require.Equal(t, note.ID, id)
				require.Nil(t, err)
			},
		},
	}
}
