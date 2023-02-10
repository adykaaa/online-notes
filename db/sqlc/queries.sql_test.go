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

//TODO: test UpdateNote

func TestDBMethods(t *testing.T) {

	ctx := context.Background()
	id := uuid.New()
	randTitle := utils.NewRandomString(15)
	randUsername := utils.NewRandomString(10)
	randText := sql.NullString{String: utils.NewRandomString(50), Valid: true}
	randCreatedAt := time.Now()
	randUpdatedAt := time.Now()
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
	t.Run("DeleteUser OK", func(t *testing.T) {

		mockdb.EXPECT().DeleteUser(ctx, randUsername).Return(randUsername, nil)
		username, err := mockdb.DeleteUser(ctx, randUsername)
		require.NoError(t, err)
		require.Equal(t, username, randUsername)
	})
}
