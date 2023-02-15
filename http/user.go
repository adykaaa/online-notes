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
	httplib "github.com/adykaaa/online-notes/lib/http"
	"github.com/adykaaa/online-notes/lib/password"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

func RegisterUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			httplib.JSON(w, msg{"error": "internal error decoding User struct"}, http.StatusInternalServerError)
			return
		}

		validate := validator.New()
		err = validate.Struct(&user)
		if err != nil {
			l.Error().Err(err).Msgf("error during User struct validation %v", err)
			httplib.JSON(w, msg{"error": "wrongly formatted or missing User parameter"}, http.StatusBadRequest)
			return
		}

		hashedPassword, err := password.Hash(user.Password)
		if err != nil {
			if errors.Is(err, password.ErrTooShort) {
				l.Error().Err(err).Msgf("The given password is too short%v", err)
				httplib.JSON(w, msg{"error": "password is too short"}, http.StatusBadRequest)
				return
			}
			l.Error().Err(err).Msgf("error during password hashing %v", err)
			httplib.JSON(w, msg{"error": "internal error during password hashing"}, http.StatusInternalServerError)
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
					httplib.JSON(w, msg{"error": "username or email already in use"}, http.StatusForbidden)
					l.Error().Err(err).Msgf("registration failed, username or email already in use for user %s", user.Username)
					return
				}
			}
			l.Error().Err(err).Msgf("Error during user registration to the DB! %v", err)
			httplib.JSON(w, msg{"error": "internal error during saving the user to the DB"}, http.StatusInternalServerError)
			return
		}

		httplib.JSON(w, msg{"success": "User registration successful!"}, http.StatusCreated)
		l.Info().Msgf("User registration for %s was successful!", uname)
	}

}

func LoginUser(q sqlc.Querier, c *PasetoCreator, tokenDuration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		var userRequest models.User

		err := json.NewDecoder(r.Body).Decode(&userRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			httplib.JSON(w, msg{"error": "internal error decoding User struct"}, http.StatusInternalServerError)
			return
		}

		user, err := q.GetUser(ctx, userRequest.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				l.Info().Err(err).Msgf("Requested user was not found in the database. %s", userRequest.Username)
				httplib.JSON(w, msg{"error": "User not found!"}, http.StatusNotFound)
				return
			}
			httplib.JSON(w, msg{"error": "interal server error while looking up user in the DB"}, http.StatusInternalServerError)
			return
		}

		err = password.Validate(user.Password, userRequest.Password)
		if err != nil {
			l.Info().Err(err).Msgf("Wrong password was provided for user %s", userRequest.Username)
			httplib.JSON(w, msg{"error": "wrong password was provided"}, http.StatusUnauthorized)
			return
		}

		token, payload, err := c.CreateToken(user.Username, tokenDuration)
		if err != nil {
			l.Info().Err(err).Msgf("Could not create PASETO for user. %v", err)
			httplib.JSON(w, msg{"error": "internal server error while creating the token"}, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "paseto",
			Value:    token,
			Expires:  payload.ExpiresAt,
			HttpOnly: true,
			Secure:   true,
		})
		httplib.JSON(w, msg{"success": "login successful"}, http.StatusOK)
		l.Info().Msgf("User login for %s was successful!", user.Username)
	}
}

func LogoutUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, _, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		username, err := io.ReadAll(r.Body)
		if err != nil {
			httplib.JSON(w, msg{"error": "couldn't decode request body"}, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "paseto",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   true,
		})
		httplib.JSON(w, msg{"success": "user successfully logged out"}, http.StatusOK)
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
