package http

import (
	"encoding/json"
	"net/http"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/adykaaa/online-notes/utils"
	"github.com/go-playground/validator/v10"
)

func RegisterUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := SetupHandler(w, r.Context())
		defer cancel()

		var user *models.User

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
		l, ctx, cancel := SetupHandler(w, r.Context())
		defer cancel()

		var user *models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			http.Error(w, "internal error decoding User struct", http.StatusInternalServerError)
			return
		}

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
