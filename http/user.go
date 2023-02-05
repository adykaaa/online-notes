package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/adykaaa/online-notes/utils"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

func RegisterUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := utils.SetupHandler(w, r.Context())
		defer cancel()

		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			utils.JSONresponse(w, map[string]string{"error": "internal error decoding User struct"}, http.StatusInternalServerError)
			return
		}

		validate := validator.New()
		err = validate.Struct(&user)
		if err != nil {
			l.Error().Err(err).Msgf("error during User struct validation %v", err)
			utils.JSONresponse(w, map[string]string{"error": "wrongly formatted or missing User parameter"}, http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			l.Error().Err(err).Msgf("error during password hashing %v", err)
			utils.JSONresponse(w, map[string]string{"error": "internal error during password hashing"}, http.StatusInternalServerError)
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
					utils.JSONresponse(w, map[string]string{"error": "username or email already in use"}, http.StatusForbidden)
					l.Error().Err(err).Msgf("registration failed, username or email already in use for us %s", user.Username)
					return
				}
			}
			l.Error().Err(err).Msgf("Error during user registration to the DB! %v", err)
			utils.JSONresponse(w, map[string]string{"error": "internal error during saving the user to the DB"}, http.StatusInternalServerError)
			return
		}

		utils.JSONresponse(w, map[string]string{"success": "User registration successful!"}, http.StatusCreated)
		l.Info().Msgf("User registration for %s was successful!", uname)
	}

}

func LoginUser(q sqlc.Querier, c *PasetoCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := utils.SetupHandler(w, r.Context())
		defer cancel()

		var userRequest models.User

		err := json.NewDecoder(r.Body).Decode(&userRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			utils.JSONresponse(w, map[string]string{"error": "internal error decoding User struct"}, http.StatusInternalServerError)
			return
		}

		user, err := q.GetUser(ctx, userRequest.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				l.Info().Err(err).Msgf("Requested user was not found in the database. %s", userRequest.Username)
				utils.JSONresponse(w, map[string]string{"error": "User not found!"}, http.StatusNotFound)
				return
			}
			utils.JSONresponse(w, map[string]string{"error": "interal server error while looking up user in the DB"}, http.StatusInternalServerError)
			return
		}

		err = utils.ValidatePassword(user.Password, userRequest.Password)
		if err != nil {
			l.Info().Err(err).Msgf("Wrong password was provided for user %s", userRequest.Username)
			utils.JSONresponse(w, map[string]string{"error": "wrong password was provided"}, http.StatusUnauthorized)
			return
		}

		token, payload, err := c.CreateToken(user.Username)
		if err != nil {
			l.Info().Err(err).Msgf("Could not create PASETO for user. %v", err)
			utils.JSONresponse(w, map[string]string{"error": "internal server error while creating the token"}, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "paseto",
			Value:    token,
			Expires:  payload.ExpiresAt,
			HttpOnly: true,
			Secure:   true,
		})
		utils.JSONresponse(w, map[string]string{"success": "login successful"}, http.StatusOK)
		l.Info().Msgf("User login for %s was successful!", user.Username)
	}
}

func LogoutUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, _, cancel := utils.SetupHandler(w, r.Context())
		defer cancel()

		username, err := io.ReadAll(r.Body)
		if err != nil {
			utils.JSONresponse(w, map[string]string{"error": "couldn't decode request body"}, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "paseto",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   true,
		})
		utils.JSONresponse(w, map[string]string{"success": "user successfully logged out"}, http.StatusOK)
		l.Info().Msgf("User logout for %s was successful!", string(username))
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
