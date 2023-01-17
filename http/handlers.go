package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/adykaaa/online-notes/db"
	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/adykaaa/online-notes/domain"
)

func Home(repo *db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("you hit the server!"))
	}
}

func RegisterUser(repo *db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user domain.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			fmt.Errorf("Could not decode response body into User! %v", err)
		}

		err = repo.RegisterUser(r.Context(), sqlc.RegisterUserParams{
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

func LoginUser(repo *db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func ListUsers(repo *db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func DeleteUser(repo *db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func CreateNote(repo *db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func GetNoteByID(repo *db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func GetAllNotesFromUser(repo *db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func DeleteNote(repo *db.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
