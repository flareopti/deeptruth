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

// Get one author
// @Summary Get one author
// @Description Get one author using ID
// @Tags authors
// @Procuce json
// @Param authorID path int true "Author ID"
// @Success 200 {object} db.Author
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/authors/{authorID} [get]
func Get(log *slog.Logger, q db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "authorID")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("no author id provided"))
			return
		}
		idInt, err := strconv.Atoi(id)
		if err != nil {
			log.Error("Failed to convert author id to int")
			log.Debug("Error", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to convert author id to int"))
			return
		}
		author, err := q.GetAuthor(r.Context(), int64(idInt))
		if err != nil {
			log.Error("Failed to fetch author")
			log.Debug("Error", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to fetch author"))
			return
		}
		render.JSON(w, r, author)
	}
}
