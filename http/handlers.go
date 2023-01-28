package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/adykaaa/online-notes/utils"
	"github.com/go-chi/render"
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
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, "internal server error during JSON decoding")
			return
		}

		validate := validator.New()
		err = validate.Struct(&user)
		if err != nil {
			l.Error().Err(err).Msgf("error during User struct validation %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad or missing user parameter"))
			return
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			l.Error().Err(err).Msgf("error during password hashing %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, "internal server error during user password hashing")
			return
		}

		uname, err := q.RegisterUser(ctx, &sqlc.RegisterUserParams{
			Username: user.Username,
			Password: hashedPassword,
			Email:    user.Email,
		})
		if err != nil {
			l.Error().Err(err).Msgf("Error during user registration to the DB! %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, "internal server error during saving user to the database")
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, "User registration successful!")
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
