package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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
			utils.JSON(w, msg{"error": "internal error decoding Note struct"}, http.StatusInternalServerError)
			return
		}

		validate := validator.New()
		err = validate.Struct(&noteRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error during Note struct validation %v", err)
			utils.JSON(w, msg{"error": "wrongly formatted or missing Note parameter"}, http.StatusBadRequest)
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
					utils.JSON(w, msg{"error": "a Note with that title already exists! Titles must be unique."}, http.StatusForbidden)
					l.Error().Err(err).Msgf("Note creation failed, a note with that title already exists")
					return
				}
			}
			l.Error().Err(err).Msgf("Error during Note creation! %v", err)
			utils.JSON(w, msg{"error": "internal error during note creation"}, http.StatusInternalServerError)
			return
		}

		utils.JSON(w, msg{"success": "note creation successful!"}, http.StatusCreated)
		l.Info().Msgf("Note with ID %v has been created for user: %s", n.ID, n.Username)
	}
}

func GetAllNotesFromUser(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := utils.SetupHandler(w, r.Context())
		defer cancel()

		username := r.URL.Query().Get("username")
		if username == "" {
			l.Error().Msgf("error fetching username, the request parameter seems empty. %s", username)
			utils.JSON(w, msg{"error": "user not in request params"}, http.StatusInternalServerError)
			return
		}

		notes, err := q.GetAllNotesFromUser(ctx, username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				l.Info().Msgf("Requested user has no Notes!. %s", username)
				utils.JSON(w, msg{"error": "user has no notes!"}, http.StatusNotFound)
				return
			}
			l.Info().Err(err).Msgf("Could not retrieve Notes for user. %v", err)
			utils.JSON(w, msg{"error": "could not retrieve notes for user"}, http.StatusInternalServerError)
			return
		}

		l.Info().Msgf("Retriving user notes for %s was successful!", username)
		utils.JSON(w, notes, http.StatusOK)
	}
}

func DeleteNote(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := utils.SetupHandler(w, r.Context())
		defer cancel()

		reqUUID, err := uuid.Parse(strings.Split(r.URL.Path, "/")[2])
		if err != nil {
			l.Info().Msgf("Could not convert ID to UUID.")
			utils.JSON(w, msg{"error": "could not convert note id to uuid"}, http.StatusBadRequest)
			return
		}

		id, err := q.DeleteNote(ctx, reqUUID)
		if err != nil {
			l.Info().Msgf("Could not delete Note %v from the DB!", reqUUID)
			utils.JSON(w, msg{"error": "could not delete note from DB"}, http.StatusInternalServerError)
			return
		}

		l.Info().Msgf("Deleting note %v was successful!", id)
		utils.JSON(w, msg{"success": "note deleted"}, http.StatusOK)
	}

}

func UpdateNote(q sqlc.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ctx, cancel := utils.SetupHandler(w, r.Context())
		defer cancel()

		reqUUID, err := uuid.Parse(strings.Split(r.URL.Path, "/")[2])
		if err != nil {
			l.Error().Err(err).Msgf("Could not convert ID to UUID.")
			utils.JSON(w, msg{"error": "could not convert note id to uuid"}, http.StatusBadRequest)
			return
		}

		updateRequest := struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{}

		err = json.NewDecoder(r.Body).Decode(&updateRequest)
		if err != nil {
			l.Error().Err(err).Msgf("error decoding the Note into JSON during registration. %v", err)
			utils.JSON(w, msg{"error": "internal error decoding Note struct"}, http.StatusInternalServerError)
			return
		}

		isTitleEmpty := false
		isTextEmpty := false
		if updateRequest.Title == "" {
			isTitleEmpty = true
		}
		if updateRequest.Text == "" {
			isTextEmpty = true
		}

		_, err = q.UpdateNote(ctx, &sqlc.UpdateNoteParams{
			ID:        reqUUID,
			Title:     sql.NullString{String: updateRequest.Title, Valid: isTitleEmpty},
			Text:      sql.NullString{String: updateRequest.Text, Valid: isTextEmpty},
			UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			l.Info().Msgf("Could not update Note with id: %v in the DB!", reqUUID)
			utils.JSON(w, msg{"error": "could not update note in DB"}, http.StatusInternalServerError)
			return
		}

		utils.JSON(w, msg{"success": "note updated!"}, http.StatusOK)
		l.Info().Msgf("Note with id: %v successfully updated!", reqUUID)
	}
}
