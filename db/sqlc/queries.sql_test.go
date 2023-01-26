package db_test

import (
	"database/sql"
	"testing"

	mockdb "github.com/adykaaa/online-notes/db/mock"
	db "github.com/adykaaa/online-notes/db/sqlc"
	http "github.com/adykaaa/online-notes/http/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetNoteById(t *testing.T) {
	n, err := http.NewNote("faketitle", "", "fakeuser")
	require.NoError(t, err)

	mockParams := db.GetNoteByIDParams{
		Username: sql.NullString{String: n.User, Valid: true},
		Title:    n.Title,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockdb := mockdb.NewMockQuerier(ctrl)
	mockdb.EXPECT().GetNoteByID(gomock.Any(), mockParams).Return(n.ID, nil)
}
