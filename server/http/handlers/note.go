package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	httplib "github.com/adykaaa/online-notes/lib/http"
	models "github.com/adykaaa/online-notes/server/http/models"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func CreateNote(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		var noteRequest models.Note

		err := json.NewDecoder(r.Body).Decode(&noteRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the Note into httplib.JSON during registration. %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error decoding Note struct"}, http.StatusInternalServerError)
			return
		}

		validate := validator.New()
		err = validate.Struct(&noteRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error during Note struct validation %v", err)
			httplib.JSON(w, httplib.Msg{"error": "wrongly formatted or missing Note parameter"}, http.StatusBadRequest)
			return
		}

		retID, err := q.CreateNote(ctx, &sqlc.CreateNoteParams{
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
					httplib.JSON(w, httplib.Msg{"error": "a Note with that title already exists! Titles must be unique."}, http.StatusForbidden)
					l.Error().Err(err).Msgf("Note creation failed, a note with that title already exists")
					return
				}
			}
			l.Error().Err(err).Msgf("Error during Note creation! %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error during note creation"}, http.StatusInternalServerError)
			return
		}

		httplib.JSON(w, httplib.Msg{"success": "note creation successful!"}, http.StatusCreated)
		l.Info().Msgf("Note with ID %v has been created for user: %s", retID, noteRequest.User)
	}
}

func GetAllNotesFromUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		username := r.URL.Query().Get("username")
		if username == "" {
			l.Error().Msgf("error fetching username, the request parameter is empty. %s", username)
			httplib.JSON(w, httplib.Msg{"error": "user not in request params"}, http.StatusBadRequest)
			return
		}

		notes, err := q.GetAllNotesFromUser(ctx, username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				l.Info().Msgf("Requested user has no Notes!. %s", username)
			}
			l.Info().Err(err).Msgf("Could not retrieve Notes for user. %v", err)
			httplib.JSON(w, httplib.Msg{"error": "could not retrieve notes for user"}, http.StatusInternalServerError)
			return
		}

		l.Info().Msgf("Retriving user notes for %s was successful!", username)
		httplib.JSON(w, notes, http.StatusOK)
	}
}

func DeleteNote(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		reqUUID, err := uuid.Parse(strings.Split(r.URL.Path, "/")[2])
		if err != nil {
			l.Info().Msgf("Could not convert ID to UUID.")
			httplib.JSON(w, httplib.Msg{"error": "could not convert note id to uuid"}, http.StatusBadRequest)
			return
		}

		id, err := q.DeleteNote(ctx, reqUUID)
		if err != nil {
			l.Info().Msgf("Could not delete Note %v from the DB!", reqUUID)
			httplib.JSON(w, httplib.Msg{"error": "could not delete note from DB"}, http.StatusInternalServerError)
			return
		}

		l.Info().Msgf("Deleting note %v was successful!", id)
		httplib.JSON(w, httplib.Msg{"success": "note deleted"}, http.StatusOK)
	}

}

func UpdateNote(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := httplib.SetupHandler(w, r.Context())
		defer cancel()

		var isTextValid bool = true

		reqUUID, err := uuid.Parse(strings.Split(r.URL.Path, "/")[2])
		if err != nil {
			l.Error().Err(err).Msgf("Could not convert ID to UUID.")
			httplib.JSON(w, httplib.Msg{"error": "could not convert note id to uuid"}, http.StatusBadRequest)
			return
		}

		updateRequest := struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{}

		err = json.NewDecoder(r.Body).Decode(&updateRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the Note into httplib.JSON during registration. %v", err)
			httplib.JSON(w, httplib.Msg{"error": "internal error decoding Note struct"}, http.StatusInternalServerError)
			return
		}

		if updateRequest.Title == "" {
			l.Error().Err(err).Msgf("Title cannot be empty!")
			httplib.JSON(w, httplib.Msg{"error": "title of a note cannot be empty"}, http.StatusBadRequest)
			return
		}

		if updateRequest.Text == "" {
			isTextValid = false
		}

		id, err := q.UpdateNote(ctx, &sqlc.UpdateNoteParams{
			ID:        reqUUID,
			Title:     sql.NullString{String: updateRequest.Title, Valid: true},
			Text:      sql.NullString{String: updateRequest.Text, Valid: isTextValid},
			UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			l.Info().Msgf("Could not update Note with id: %v in the DB!", reqUUID)
			httplib.JSON(w, httplib.Msg{"error": "could not update note in DB"}, http.StatusInternalServerError)
			return
		}

		httplib.JSON(w, httplib.Msg{"success": "note updated!"}, http.StatusOK)
		l.Info().Msgf("Note with id: %v successfully updated!", id)
	}
}
