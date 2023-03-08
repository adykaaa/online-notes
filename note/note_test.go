package note

import (
	"context"
	"testing"

	mockdb "github.com/adykaaa/online-notes/db/mock"
	db "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/adykaaa/online-notes/server/http/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	user := &models.User{
		Username: "user1",
		Password: "password1",
		Email:    "user1@user.com",
	}

	testCases := []struct {
		name              string
		user              *models.User
		mockdbCreateUser  func(mockdb *mockdb.MockQuerier, user *models.User)
		checkReturnValues func(t *testing.T, user *models.User, username string, err error)
	}{
		{
			name: "user registration OK",
			user: user,
			mockdbCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
				args := db.RegisterUserParams{
					Username: user.Username,
					Password: user.Password,
					Email:    user.Email,
				}
				mockdb.EXPECT().RegisterUser(gomock.Any(), &args).Times(1).Return(args.Username, nil)
			},
			checkReturnValues: func(t *testing.T, user *models.User, username string, err error) {
				require.Equal(t, username, user.Username)
				require.Nil(t, err)
			},
		},
		{
			name: "user registration returns ErrUserAlreadyExists",
			user: user,
			mockdbCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
				args := db.RegisterUserParams{
					Username: user.Username,
					Password: user.Password,
					Email:    user.Email,
				}
				mockdb.EXPECT().RegisterUser(gomock.Any(), &args).Times(1).Return("", ErrUserAlreadyExists)
			},
			checkReturnValues: func(t *testing.T, user *models.User, username string, err error) {
				require.ErrorIs(t, err, ErrUserAlreadyExists)
				require.Empty(t, username)
			},
		},
		{
			name: "user registration returns ErrDBInternal",
			user: user,
			mockdbCreateUser: func(mockdb *mockdb.MockQuerier, user *models.User) {
				args := db.RegisterUserParams{
					Username: user.Username,
					Password: user.Password,
					Email:    user.Email,
				}
				mockdb.EXPECT().RegisterUser(gomock.Any(), &args).Times(1).Return("", ErrDBInternal)
			},
			checkReturnValues: func(t *testing.T, user *models.User, username string, err error) {
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

			tc.mockdbCreateUser(mockdb, tc.user)

			u, err := ns.q.RegisterUser(context.Background(), &db.RegisterUserParams{
				Username: tc.user.Username,
				Password: tc.user.Password,
				Email:    tc.user.Email,
			})
			tc.checkReturnValues(t, tc.user, u, err)
		})
	}
}
