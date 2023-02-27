package http

/* import (
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/adykaaa/online-notes/db/mock"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateNote(t *testing.T) {
	jsonValidator := validator.New()

	testCases := []struct {
		name             string
		body             *models.Note
		validateJSON     func(t *testing.T, v *validator.Validate, user *models.Note)
		dbmockCreateNote func(mockdb *mockdb.MockQuerier, user *models.User)
		checkResponse    func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request)
	}{
		{
			name: "Note creation OK",

			body: &models.Note{
				Username: "user1",
				Password: "password1",
				Email:    "user1@user.com",
			},

			validateJSON: func(t *testing.T, v *validator.Validate, user *models.Note) {
				err := v.Struct(user)
				require.NoError(t, err)
			},

			dbmockCreateNote: func(mockdb *mockdb.MockQuerier, user *models.User) {
				args := regUserArgs{
					Username: user.Username,
					Password: user.Password,
					Email:    user.Email,
				}
				mockdb.EXPECT().RegisterUser(gomock.Any(), &args).Times(1).Return(args.Username, nil)
			},

			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, request *http.Request) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}
}
*/
