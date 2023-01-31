package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/adykaaa/online-notes/utils"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

func RegisterUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := SetupHandler(w, r.Context())
		defer cancel()

		var user models.User

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
			http.Error(w, "wrongly formatted or missing User parameter", http.StatusBadRequest)
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
			if postgreError, ok := err.(*pq.Error); ok {
				if postgreError.Code.Name() == "unique_violation" {
					http.Error(w, "username or email already in use", http.StatusForbidden)
					l.Error().Err(err).Msgf("registration failed, username or email already in use for us %s", user.Username)
					return
				}
			}
			l.Error().Err(err).Msgf("Error during user registration to the DB! %v", err)
			http.Error(w, "internal error during saving the user to the DB", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registration successful!"))
		l.Info().Msgf("User registration for %s was successful!", uname)
	}

}

func LoginUser(q sqlc.Querier, c *PasetoCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := SetupHandler(w, r.Context())
		defer cancel()

		var userRequest *models.User

		err := json.NewDecoder(r.Body).Decode(&userRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			http.Error(w, "internal error decoding User struct", http.StatusInternalServerError)
			return
		}

		validate := validator.New()
		err = validate.Struct(&userRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error during User struct validation %v", err)
			http.Error(w, "wrongly formatted or missing User parameter", http.StatusBadRequest)
			return
		}

		user, err := q.GetUser(ctx, userRequest.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				l.Info().Err(err).Msgf("Requested user was not found in the database. %s", userRequest.Username)
				http.Error(w, "User not found!", http.StatusNotFound)
				return
			}
			http.Error(w, "interal server error while looking up user in the DB", http.StatusInternalServerError)
			return
		}

		err = utils.ValidatePassword(user.Password, userRequest.Password)
		if err != nil {
			l.Info().Err(err).Msgf("Wrong password was provided for user %s", userRequest.Username)
			http.Error(w, "wrong password was provided", http.StatusUnauthorized)
			return
		}

		token, payload, err := c.CreateToken(user.Username)
		if err != nil {
			l.Info().Err(err).Msgf("Could not create PASETO for user. %v", err)
			http.Error(w, "internal server error while creating the token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "authentication",
			Value:    token,
			Expires:  payload.ExpiresAt,
			HttpOnly: true,
			Secure:   true,
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successful login!"))
		l.Info().Msgf("User login for %s was successful!", user.Username)

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
