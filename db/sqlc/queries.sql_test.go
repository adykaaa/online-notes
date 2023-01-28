package db_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	mockdb "github.com/adykaaa/online-notes/db/mock"
	db "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/adykaaa/online-notes/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBMethods(t *testing.T) {

	ctx := context.Background()
	id := uuid.New()
	randTitle := utils.NewRandomString(15)
	randUsername := sql.NullString{String: utils.NewRandomString(10), Valid: true}
	randText := sql.NullString{String: utils.NewRandomString(50), Valid: true}
	randCreatedAt := sql.NullTime{Time: time.Now(), Valid: true}
	randUpdatedAt := sql.NullTime{Time: time.Now(), Valid: true}
	randPassword := utils.NewRandomString(15)
	randEmail := "random@random.com"

	ctrl := gomock.NewController(t)
	mockdb := mockdb.NewMockQuerier(ctrl)

	//TODO: test not happy cases as well
	t.Run("CreateNote OK", func(t *testing.T) {

		n := db.Note{
			ID:        id,
			Title:     randTitle,
			Username:  randUsername,
			Text:      randText,
			CreatedAt: randCreatedAt,
			UpdatedAt: randUpdatedAt,
		}

		args := db.CreateNoteParams{
			ID:        id,
			Title:     randTitle,
			Username:  randUsername,
			Text:      randText,
			CreatedAt: randCreatedAt,
			UpdatedAt: randUpdatedAt,
		}

		mockdb.EXPECT().CreateNote(ctx, &args).Return(n, nil)
		retNote, err := mockdb.CreateNote(ctx, &args)

		assert.NoError(t, err)
		assert.NotNil(t, id)
		assert.Equal(t, args.ID, retNote.ID)
		assert.Equal(t, args.Title, retNote.Title)
		assert.Equal(t, args.Username, retNote.Username)
		assert.Equal(t, args.Text, retNote.Text)
		assert.Equal(t, args.CreatedAt, retNote.CreatedAt)
		assert.Equal(t, args.UpdatedAt, retNote.UpdatedAt)

	})
	t.Run("DeleteNote OK", func(t *testing.T) {

		mockdb.EXPECT().DeleteNote(ctx, id).Return(id, nil)
		retID, err := mockdb.DeleteNote(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, id)
		assert.Equal(t, id, retID)
	})
	t.Run("GetAllNotesFromUser OK", func(t *testing.T) {
		randomNotes := []db.Note{
			*utils.NewRandomDBNote(id),
			*utils.NewRandomDBNote(id),
		}

		randomNotes[0].Username = randomNotes[1].Username

		mockdb.EXPECT().GetAllNotesFromUser(ctx, randUsername).Return(randomNotes, nil)
		notes, err := mockdb.GetAllNotesFromUser(ctx, randUsername)

		assert.NoError(t, err)
		assert.NotEmpty(t, notes[0].Title)
		assert.NotEmpty(t, notes[1].Title)
		assert.NotEmpty(t, notes[0].Username)
		assert.NotEmpty(t, notes[1].Username)
		assert.Equal(t, notes[0].Username, notes[1].Username)
		assert.NotEqual(t, notes[0].CreatedAt, notes[1].CreatedAt)
	})

	t.Run("GetNoteByID OK", func(t *testing.T) {
		args := db.GetNoteByIDParams{
			Username: randUsername,
			Title:    randTitle,
		}

		mockdb.EXPECT().GetNoteByID(ctx, &args).Return(id, nil)
		retID, err := mockdb.GetNoteByID(ctx, &args)

		require.NoError(t, err)
		require.Equal(t, retID, id)
	})
	t.Run("UpdateNoteText OK", func(t *testing.T) {
		args := db.UpdateNoteTextParams{
			ID:        id,
			Text:      randText,
			UpdatedAt: randUpdatedAt,
		}

		mockdb.EXPECT().UpdateNoteText(ctx, &args).Return(db.Note{
			ID:        id,
			Text:      randText,
			UpdatedAt: randUpdatedAt,
		}, nil)

		retNote, err := mockdb.UpdateNoteText(ctx, &args)

		require.NoError(t, err)
		require.Equal(t, retNote.ID, args.ID)
		require.Equal(t, retNote.Text, args.Text)
		require.Equal(t, retNote.UpdatedAt, args.UpdatedAt)
	})
	t.Run("UpdateNoteTitle OK", func(t *testing.T) {
		args := db.UpdateNoteTitleParams{
			ID:        id,
			Title:     randTitle,
			UpdatedAt: randUpdatedAt,
		}

		mockdb.EXPECT().UpdateNoteTitle(ctx, &args).Return(db.Note{
			ID:        id,
			Title:     randTitle,
			UpdatedAt: randUpdatedAt,
		}, nil)

		retNote, err := mockdb.UpdateNoteTitle(ctx, &args)

		require.NoError(t, err)
		require.Equal(t, retNote.ID, args.ID)
		require.Equal(t, retNote.Title, args.Title)
		require.Equal(t, retNote.UpdatedAt, args.UpdatedAt)
	})
	t.Run("UpdateNoteTitle OK", func(t *testing.T) {
		args := db.UpdateNoteTitleParams{
			ID:        id,
			Title:     randTitle,
			UpdatedAt: randUpdatedAt,
		}

		mockdb.EXPECT().UpdateNoteTitle(ctx, &args).Return(db.Note{
			ID:        id,
			Title:     randTitle,
			UpdatedAt: randUpdatedAt,
		}, nil)

		retNote, err := mockdb.UpdateNoteTitle(ctx, &args)

		require.NoError(t, err)
		require.Equal(t, retNote.ID, args.ID)
		require.Equal(t, retNote.Title, args.Title)
		require.Equal(t, retNote.UpdatedAt, args.UpdatedAt)
	})
	t.Run("RegisterUser OK", func(t *testing.T) {
		args := db.RegisterUserParams{
			Username: randUsername.String,
			Password: randPassword,
			Email:    randEmail,
		}
		mockdb.EXPECT().RegisterUser(ctx, &args).Return(args.Username, nil)
		user, err := mockdb.RegisterUser(ctx, &args)
		require.NoError(t, err)
		require.Equal(t, user, args.Username)
	})
	t.Run("DeleteUser OK", func(t *testing.T) {

		mockdb.EXPECT().DeleteUser(ctx, randUsername.String).Return(randUsername.String, nil)
		username, err := mockdb.DeleteUser(ctx, randUsername.String)
		require.NoError(t, err)
		require.Equal(t, username, randUsername.String)
	})
}
