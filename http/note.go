package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/adykaaa/online-notes/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func CreateNote(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := utils.SetupHandler(w, r.Context())
		defer cancel()

		var noteRequest models.Note

		err := json.NewDecoder(r.Body).Decode(&noteRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the Note into JSON during registration. %v", err)
			utils.JSONresponse(w, map[string]string{"error": "internal error decoding Note struct"}, http.StatusInternalServerError)
			return
		}

		validate := validator.New()
		err = validate.Struct(&noteRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error during Note struct validation %v", err)
			utils.JSONresponse(w, map[string]string{"error": "wrongly formatted or missing Note parameter"}, http.StatusBadRequest)
			return
		}

		n, err := q.CreateNote(ctx, &sqlc.CreateNoteParams{
			ID:        uuid.New(),
			Title:     noteRequest.Title,
			Username:  noteRequest.User,
			Text:      sql.NullString{String: noteRequest.Text, Valid: true},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			if postgreError, ok := err.(*pq.Error); ok {
				if postgreError.Code.Name() == "unique_violation" {
					utils.JSONresponse(w, map[string]string{"error": "a Note with that title already exists! Titles must be unique."}, http.StatusForbidden)
					l.Error().Err(err).Msgf("Note creation failed, a note with that title already exists")
					return
				}
			}
			l.Error().Err(err).Msgf("Error during Note creation! %v", err)
			utils.JSONresponse(w, map[string]string{"error": "internal error during note creation"}, http.StatusInternalServerError)
			return
		}

		utils.JSONresponse(w, map[string]string{"success": "note creation successful!"}, http.StatusCreated)
		l.Info().Msgf("Note with ID %v has been created for user: %s", n.ID, n.Username)
	}
}

func GetNoteByID(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func GetAllNotesFromUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := utils.SetupHandler(w, r.Context())
		defer cancel()

		getNotesRequest := struct {
			Username string `json:"user"`
		}{}

		err := json.NewDecoder(r.Body).Decode(&getNotesRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding request. %v", err)
			utils.JSONresponse(w, map[string]string{"error": "internal error decoding request"}, http.StatusInternalServerError)
			return
		}

		notes, err := q.GetAllNotesFromUser(ctx, getNotesRequest.Username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				l.Info().Msgf("Requested user has no Notes!. %s", getNotesRequest.Username)
				utils.JSONresponse(w, map[string]string{"error": "user has no notes!"}, http.StatusNotFound)
				return
			}
			l.Info().Err(err).Msgf("Could not retrieve Notes for user. %v", err)
			utils.JSONresponse(w, map[string]string{"error": "could not retrieve notes for user"}, http.StatusInternalServerError)
			return
		}

		l.Info().Msgf("Retriving user notes for %s was successful!", getNotesRequest.Username)
		utils.JSONresponse(w, notes, http.StatusOK)
	}
}

func DeleteNote(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func UpdateNote(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
