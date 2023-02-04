package http

import (
	"database/sql"
	"encoding/json"
	"net/http"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	models "github.com/adykaaa/online-notes/http/models"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)

func CreateNote(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := SetupHandler(w, r.Context())
		defer cancel()

		var noteRequest models.Note

		err := json.NewDecoder(r.Body).Decode(&noteRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the Note into JSON during registration. %v", err)
			http.Error(w, "internal error decoding Note struct", http.StatusInternalServerError)
			return
		}

		validate := validator.New()
		err = validate.Struct(&noteRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error during Note struct validation %v", err)
			http.Error(w, "wrongly formatted or missing Note parameter", http.StatusBadRequest)
			return
		}

		n, err := q.CreateNote(ctx, &sqlc.CreateNoteParams{
			ID:        noteRequest.ID,
			Title:     noteRequest.Title,
			Username:  sql.NullString{String: noteRequest.User, Valid: true},
			Text:      sql.NullString{String: noteRequest.Text, Valid: true},
			CreatedAt: sql.NullTime{Time: noteRequest.CreatedAt, Valid: true},
			UpdatedAt: sql.NullTime{Time: noteRequest.UpdatedAt, Valid: true},
		})
		if err != nil {
			if postgreError, ok := err.(*pq.Error); ok {
				if postgreError.Code.Name() == "unique_violation" {
					http.Error(w, "A Note with that title already exists! Titles must be unique.", http.StatusBadRequest)
					l.Error().Err(err).Msgf("Note creation failed, a note with that title already exists")
					return
				}
			}
			l.Error().Err(err).Msgf("Error during Note creation! %v", err)
			http.Error(w, "internal error during note creation", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Note creation successful!"))
		l.Info().Msgf("Note with ID %v has been created for user: %s", n.ID, n.Username.String)
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
