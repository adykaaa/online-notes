package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/adykaaa/online-notes/utils"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

func Home(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("home!"))
	}
}

func RegisterUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var user models.User

		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		l := zerolog.Ctx(ctx)

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			http.Error(w, "internal error decoding User struct", http.StatusInternalServerError)
			return
		}

		validate := validator.New()
		err = validate.Struct(&user)
		if err != nil {
			l.Error().Err(err).Msgf("error during User struct validation %v", err)
			http.Error(w, "invalid or missing User parameter", http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			l.Error().Err(err).Msgf("error during password hashing %v", err)
			http.Error(w, "internal error during password hashing", http.StatusInternalServerError)
			return
		}

		uname, err := q.RegisterUser(ctx, &sqlc.RegisterUserParams{
			Username: user.Username,
			Password: hashedPassword,
			Email:    user.Email,
		})
		if err != nil {
			l.Error().Err(err).Msgf("Error during user registration to the DB! %v", err)
			http.Error(w, "internal error during saving the user to the DB", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registration successful!"))
		l.Info().Msgf("User registration for %s was successful!", uname)
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
