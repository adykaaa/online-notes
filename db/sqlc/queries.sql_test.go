package db_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	mockdb "github.com/adykaaa/online-notes/db/mock"
	db "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/adykaaa/online-notes/lib/random"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: test not happy cases as well, transform to table-driven tests, etc.
func TestDBMethods(t *testing.T) {

	ctx := context.Background()
	id := uuid.New()
	randTitle := random.NewString(15)
	randUsername := random.NewString(10)
	randText := sql.NullString{String: random.NewString(50), Valid: true}
	randCreatedAt := time.Now()
	randUpdatedAt := time.Now()
	randPassword := random.NewString(15)
	randEmail := "random@random.com"

	ctrl := gomock.NewController(t)
	mockdb := mockdb.NewMockQuerier(ctrl)

	t.Run("CreateNote OK", func(t *testing.T) {
		args := db.CreateNoteParams{
			ID:        id,
			Title:     randTitle,
			Username:  randUsername,
			Text:      randText,
			CreatedAt: randCreatedAt,
			UpdatedAt: randUpdatedAt,
		}

		mockdb.EXPECT().CreateNote(ctx, &args).Return(args.ID, nil)
		retID, err := mockdb.CreateNote(ctx, &args)

		assert.NoError(t, err)
		assert.NotNil(t, id)
		assert.Equal(t, args.ID, retID)
	})
	t.Run("UpdateNote OK", func(t *testing.T) {
		updatedAt := time.Date(2021, 8, 15, 14, 30, 45, 100, time.Local)
		args := &db.UpdateNoteParams{
			ID:        id,
			Title:     sql.NullString{String: "updated title", Valid: true},
			Text:      sql.NullString{String: "updated text", Valid: true},
			UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
		}

		mockdb.EXPECT().UpdateNote(ctx, args).Return(args.ID, nil)

		retID, err := mockdb.UpdateNote(ctx, args)
		require.NoError(t, err)
		require.Equal(t, args.ID, retID)
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
			*random.NewDBNote(id),
			*random.NewDBNote(id),
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

	t.Run("RegisterUser OK", func(t *testing.T) {
		args := db.RegisterUserParams{
			Username: randUsername,
			Password: randPassword,
			Email:    randEmail,
		}
		mockdb.EXPECT().RegisterUser(ctx, &args).Return(args.Username, nil)
		user, err := mockdb.RegisterUser(ctx, &args)
		require.NoError(t, err)
		require.Equal(t, user, args.Username)
	})
}
