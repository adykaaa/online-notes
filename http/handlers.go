package http

import (
	"encoding/json"
	"net/http"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/adykaaa/online-notes/utils"
	"github.com/rs/zerolog"
)

func Home(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("you hit the server!"))
	}
}

func RegisterUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		l := zerolog.Ctx(r.Context())

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration!", err)
			return
		}

		hashedPassword, err := utils.HashUserPassword(user.Password)
		if err != nil {
			l.Error().Err(err).Msgf("error during password hashing", err)
			return
		}

		err = q.RegisterUser(r.Context(), sqlc.RegisterUserParams{
			Username: user.Username,
			Password: hashedPassword,
			Email:    user.Email,
		})
		if err != nil {
			l.Error().Err(err).Msgf("Error during user registration to DB! %v", err)
			return
		}

		l.Info().Msgf("User registration for %v was successful!", err)
	}
}

func LoginUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func ListUsers(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func DeleteUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func CreateNote(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func GetNoteByID(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func GetAllNotesFromUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func DeleteNote(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
