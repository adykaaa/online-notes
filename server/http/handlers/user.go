package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	httplib "github.com/adykaaa/online-notes/lib/http"
	"github.com/adykaaa/online-notes/lib/password"
	auth "github.com/adykaaa/online-notes/server/http/auth"
	models "github.com/adykaaa/online-notes/server/http/models"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

func RegisterUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		var userRequest models.User
		err := json.NewDecoder(r.Body).Decode(&userRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error decoding User struct"}, http.StatusInternalServerError)
			return
		}

		validate := validator.New()
		err = validate.Struct(&userRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error during User struct validation %v", err)
			httplib.JSON(w, httplib.Msg{"error": "wrongly formatted or missing User parameter"}, http.StatusBadRequest)
			return
		}

		hashedPassword, err := password.Hash(userRequest.Password)
		if err != nil {
			if errors.Is(err, password.ErrTooShort) {
				l.Error().Err(err).Msgf("The given password is too short%v", err)
				httplib.JSON(w, httplib.Msg{"error": "password is too short"}, http.StatusBadRequest)
				return
			}
			l.Error().Err(err).Msgf("error during password hashing %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error during password hashing"}, http.StatusInternalServerError)
			return
		}

		uname, err := q.RegisterUser(ctx, &sqlc.RegisterUserParams{
			Username: userRequest.Username,
			Password: hashedPassword,
			Email:    userRequest.Email,
		})
		if err != nil {
			if postgreError, ok := err.(*pq.Error); ok {
				if postgreError.Code.Name() == "unique_violation" {
					httplib.JSON(w, httplib.Msg{"error": "username or email already in use"}, http.StatusForbidden)
					l.Error().Err(err).Msgf("registration failed, username or email already in use for user %s", userRequest.Username)
					return
				}
			}
			l.Error().Err(err).Msgf("Error during user registration to the DB! %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error during saving the user to the DB"}, http.StatusInternalServerError)
			return
		}

		httplib.JSON(w, httplib.Msg{"success": "User registration successful!"}, http.StatusCreated)
		l.Info().Msgf("User registration for %s was successful!", uname)
	}

}

func LoginUser(q sqlc.Querier, t auth.TokenManager, tokenDuration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		userRequest := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&userRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error decoding User struct"}, http.StatusInternalServerError)
			return
		}

		dbuser, err := q.GetUser(ctx, userRequest.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				l.Info().Err(err).Msgf("Requested user was not found in the database. %s", userRequest.Username)
				httplib.JSON(w, httplib.Msg{"error": "User not found!"}, http.StatusNotFound)
				return
			}
			httplib.JSON(w, httplib.Msg{"error": "interal server error while looking up user in the DB"}, http.StatusInternalServerError)
			return
		}

		err = password.Validate(dbuser.Password, userRequest.Password)
		if err != nil {
			l.Info().Err(err).Msgf("Wrong password was provided for user %s", userRequest.Username)
			httplib.JSON(w, httplib.Msg{"error": "wrong password was provided"}, http.StatusUnauthorized)
			return
		}

		token, payload, err := t.CreateToken(userRequest.Username, tokenDuration)
		if err != nil {
			l.Info().Err(err).Msgf("Could not create PASETO for user. %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal server error while creating the token"}, http.StatusInternalServerError)
			return
		}

		httplib.SetCookie(w, "paseto", token, payload.ExpiresAt)
		httplib.JSON(w, httplib.Msg{"success": "login successful"}, http.StatusOK)
		l.Info().Msgf("User login for %s was successful!", userRequest.Username)
	}
}

func LogoutUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, _, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		username, err := io.ReadAll(r.Body)
		if err != nil {
			httplib.JSON(w, httplib.Msg{"error": "couldn't decode request body"}, http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "paseto",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   true,
		})
		httplib.JSON(w, httplib.Msg{"success": "user successfully logged out"}, http.StatusOK)
		l.Info().Msgf("User logout for %s was successful!", string(username))
	}
}

// TODO: implement these in the future
func ListUsers(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func DeleteUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}