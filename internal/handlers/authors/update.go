package authors

import (
	"log/slog"
	"net/http"
	"strconv"

	db "github.com/flareopti/deeptruth/internal/db/sqlc"
	"github.com/flareopti/deeptruth/internal/lib/api/resp"
	"github.com/flareopti/deeptruth/internal/lib/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UpdateRatingParam struct {
	Rating int
}

// Author rating update
// @Summary Update rating of an author
// @Description Update rating of an author
// @Tags authors
// @Accept json
// @Procuce json
// @Param authorID path int true "Author ID"
// @Param rating body UpdateRatingParam true "New rating"
// @Success 200 {object} db.Author
// @Router /api/authors/{authorID} [patch]
func UpdateRating(log *slog.Logger, q db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "authorID")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("no author id provided"))
			return
		}
		id_int, err := strconv.Atoi(id)
		if err != nil {
			log.Error("Failed to convert author id to int")
			log.Debug("Error", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to convert author id to int"))
			return
		}
		var rating UpdateRatingParam
		err = render.DecodeJSON(r.Body, &rating)
		author_params := db.UpdateAuthorRatingParams{
			ID:     int64(id_int),
			Rating: int32(rating.Rating),
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("Failed to decode request body")
			log.Debug("Error", sl.Err(err))
			render.JSON(w, r, resp.Error("Failed to decode request"))
			return
		}
		author, err := q.UpdateAuthorRating(r.Context(), author_params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Failed to update author")
			log.Debug("Error:", err)
			render.JSON(w, r, resp.Error("Failed to update author"))
			return
		}
		render.JSON(w, r, author)
	}
}
