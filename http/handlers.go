package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/adykaaa/online-notes/domain"
	"github.com/rs/zerolog"
)

func Home(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("you hit the server!"))
	}
}

func RegisterUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user domain.User
		l := zerolog.Ctx(r.Context())

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			l.Info().Err(err).Msg("Error decoding the User into JSON during registration!")
		}

		err = q.RegisterUser(r.Context(), sqlc.RegisterUserParams{
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
		})
		if err != nil {
			fmt.Errorf("Error during user registration! %v", err)
		}

		log.Printf("User registration for %s was successful!", user.Username)
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
