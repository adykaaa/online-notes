package db_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	mockdb "github.com/adykaaa/online-notes/db/mock"
	db "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateNote(t *testing.T) {

	ctx := context.Background()
	id := uuid.New()

	args := db.CreateNoteParams{
		ID:        id,
		Title:     "Fake Note Title",
		Username:  sql.NullString{String: "fake_username123", Valid: true},
		Text:      sql.NullString{String: "fake text inside note", Valid: true},
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	ctrl := gomock.NewController(t)
	mockdb := mockdb.NewMockQuerier(ctrl)

	mockdb.EXPECT().CreateNote(ctx, args).Return(id, nil)
	retID, err := mockdb.CreateNote(ctx, args)

	assert.NoError(t, err)
	assert.NotNil(t, id)
	assert.Equal(t, id, retID)

}
