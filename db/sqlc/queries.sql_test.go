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
	title := utils.NewRandomString(15)
	username := sql.NullString{String: utils.NewRandomString(10), Valid: true}
	text := sql.NullString{String: utils.NewRandomString(50), Valid: true}
	createdAt := sql.NullTime{Time: time.Now(), Valid: true}
	updatedAt := sql.NullTime{Time: time.Now(), Valid: true}

	ctrl := gomock.NewController(t)
	mockdb := mockdb.NewMockQuerier(ctrl)

	//TODO: test not happy cases as well
	t.Run("CreateNote OK", func(t *testing.T) {

		n := db.Note{
			ID:        id,
			Title:     title,
			Username:  username,
			Text:      text,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		args := db.CreateNoteParams{
			ID:        id,
			Title:     title,
			Username:  username,
			Text:      text,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		mockdb.EXPECT().CreateNote(ctx, args).Return(n, nil)
		retNote, err := mockdb.CreateNote(ctx, args)

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
			*utils.NewRandomNote(id),
			*utils.NewRandomNote(id),
		}

		randomNotes[0].Username = randomNotes[1].Username

		mockdb.EXPECT().GetAllNotesFromUser(ctx, username).Return(randomNotes, nil)
		notes, err := mockdb.GetAllNotesFromUser(ctx, username)

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
			Username: username,
			Title:    title,
		}

		mockdb.EXPECT().GetNoteByID(ctx, args).Return(id, nil)
		retID, err := mockdb.GetNoteByID(ctx, args)

		require.NoError(t, err)
		require.Equal(t, retID.ID, id)
	})
	t.Run("UpdateNoteText OK", func(t *testing.T) {
		args := db.UpdateNoteTextParams{
			ID:        id,
			Text:      text,
			UpdatedAt: updatedAt,
		}

		mockdb.EXPECT().UpdateNoteText(ctx, args).Return(id, nil)
		retID, err := mockdb.UpdateNoteText(ctx, args)

		require.NoError(t, err)
		require.Equal(t, retID.ID, id)
	})

}
