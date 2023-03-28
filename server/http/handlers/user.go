package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	db "github.com/adykaaa/online-notes/db/sqlc"
	httplib "github.com/adykaaa/online-notes/lib/http"
	"github.com/adykaaa/online-notes/lib/password"
	"github.com/adykaaa/online-notes/note"
	auth "github.com/adykaaa/online-notes/server/http/auth"
	models "github.com/adykaaa/online-notes/server/http/models"
	"github.com/go-playground/validator/v10"
)

func RegisterUser(s NoteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		var req models.User
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error decoding User struct"}, http.StatusInternalServerError)
			return
		}

		validate := validator.New()
		err = validate.Struct(&req)
		if err != nil {
			l.Error().Err(err).Msgf("error during User struct validation %v", err)
			httplib.JSON(w, httplib.Msg{"error": "wrongly formatted or missing User parameter"}, http.StatusBadRequest)
			return
		}

		hashedPw, err := password.Hash(req.Password)
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

		uname, err := s.RegisterUser(ctx, &db.RegisterUserParams{
			Username: req.Username,
			Password: hashedPw,
			Email:    req.Email,
		})

		switch {
		case errors.Is(err, note.ErrAlreadyExists):
			l.Error().Err(err).Msgf("registration failed, username or email already in use for user %s", req.Username)
			httplib.JSON(w, httplib.Msg{"error": "username or email already in use"}, http.StatusForbidden)
			return
		case errors.Is(err, note.ErrDBInternal):
			l.Error().Err(err).Msgf("Error during User registration! %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error during user registration"}, http.StatusInternalServerError)
			return
		default:
			httplib.JSON(w, httplib.Msg{"success": "User registration successful!"}, http.StatusCreated)
			l.Info().Msgf("User registration for %s was successful!", uname)
		}
	}
}

func LoginUser(s NoteService, t auth.TokenManager, tokenDuration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		req := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the User into JSON during registration. %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error decoding User struct"}, http.StatusInternalServerError)
			return
		}

		user, err := s.GetUser(ctx, req.Username)
		switch {
		case errors.Is(err, note.ErrNotFound):
			l.Error().Err(err).Msgf("user: %s is not found", req.Username)
			httplib.JSON(w, httplib.Msg{"error": "user is not found"}, http.StatusForbidden)
			return
		case errors.Is(err, note.ErrDBInternal):
			l.Error().Err(err).Msgf("Error during user lookup! %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error during user lookup!"}, http.StatusInternalServerError)
			return
		}

		err = password.Validate(user.Password, req.Password)
		if err != nil {
			l.Info().Err(err).Msgf("Wrong password was provided for user %s", req.Username)
			httplib.JSON(w, httplib.Msg{"error": "wrong password was provided"}, http.StatusUnauthorized)
			return
		}

		token, payload, err := t.CreateToken(req.Username, tokenDuration)
		if err != nil {
			l.Info().Err(err).Msgf("Could not create PASETO for user. %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal server error while creating the token"}, http.StatusInternalServerError)
			return
		}

		httplib.SetCookie(w, "paseto", token, payload.ExpiresAt)
		httplib.JSON(w, httplib.Msg{"success": "login successful"}, http.StatusOK)
		l.Info().Msgf("User login for %s was successful!", req.Username)
	}
}

func LogoutUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, _, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		uname, err := io.ReadAll(r.Body)
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
		l.Info().Msgf("User logout for %s was successful!", string(uname))
	}
}
