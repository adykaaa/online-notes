package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	httplib "github.com/adykaaa/online-notes/lib/http"
	"github.com/adykaaa/online-notes/lib/password"
	auth "github.com/adykaaa/online-notes/server/http/auth"
	models "github.com/adykaaa/online-notes/server/http/models"
	"github.com/adykaaa/online-notes/user"
	"github.com/go-playground/validator/v10"
)

func RegisterUser(us user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		var request models.User
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error decoding User struct"}, http.StatusInternalServerError)
			return
		}

		validate := validator.New()
		err = validate.Struct(&request)
		if err != nil {
			l.Error().Err(err).Msgf("error during User struct validation %v", err)
			httplib.JSON(w, httplib.Msg{"error": "wrongly formatted or missing User parameter"}, http.StatusBadRequest)
			return
		}

		hashedPassword, err := password.Hash(request.Password)
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

		uname, err := us.RegisterUser(ctx, request.Username, hashedPassword, request.Email)
		switch {
		case errors.Is(err, user.ErrAlreadyExists):
			l.Error().Err(err).Msgf("registration failed, username or email already in use for user %s", request.Username)
			httplib.JSON(w, httplib.Msg{"error": "username or email already in use"}, http.StatusForbidden)
			return
		case errors.Is(err, user.ErrDBInternal):
			l.Error().Err(err).Msgf("Error during User registration! %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error during user registration"}, http.StatusInternalServerError)
			return
		default:
			httplib.JSON(w, httplib.Msg{"success": "User registration successful!"}, http.StatusCreated)
			l.Info().Msgf("User registration for %s was successful!", uname)
		}
	}
}

func LoginUser(us user.UserService, t auth.TokenManager, tokenDuration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		request := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error decoding User struct"}, http.StatusInternalServerError)
			return
		}

		dbuser, err := us.GetUser(ctx, request.Username)
		switch {
		case errors.Is(err, user.ErrNotFound):
			l.Error().Err(err).Msgf("user: %s is not found", request.Username)
			httplib.JSON(w, httplib.Msg{"error": "user is not found"}, http.StatusForbidden)
			return
		case errors.Is(err, user.ErrDBInternal):
			l.Error().Err(err).Msgf("Error during user lookup! %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error during user lookup!"}, http.StatusInternalServerError)
			return
		}

		err = password.Validate(dbuser.Password, request.Password)
		if err != nil {
			l.Info().Err(err).Msgf("Wrong password was provided for user %s", request.Username)
			httplib.JSON(w, httplib.Msg{"error": "wrong password was provided"}, http.StatusUnauthorized)
			return
		}

		token, payload, err := t.CreateToken(request.Username, tokenDuration)
		if err != nil {
			l.Info().Err(err).Msgf("Could not create PASETO for user. %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal server error while creating the token"}, http.StatusInternalServerError)
			return
		}

		httplib.SetCookie(w, "paseto", token, payload.ExpiresAt)
		httplib.JSON(w, httplib.Msg{"success": "login successful"}, http.StatusOK)
		l.Info().Msgf("User login for %s was successful!", request.Username)
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
