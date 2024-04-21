package authors

import (
	"log/slog"
	"net/http"

	db "github.com/flareopti/deeptruth/internal/db/sqlc"
	"github.com/flareopti/deeptruth/internal/lib/api/resp"
	"github.com/go-chi/render"
)

// Author create
// @Summary Create an author
// @Description Create an author
// @Tags authors
// @Accept json
// @Procuce json
// @Param author body db.CreateAuthorParams true "Author to create"
// @Success 200 {object} db.Author
// @Router /api/authors [post]
func Create(log *slog.Logger, q db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authorParams db.CreateAuthorParams
		err := render.DecodeJSON(r.Body, &authorParams)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("Failed to decode request body")
			log.Debug("Error:", err)
			render.JSON(w, r, resp.Error("Failed to decode request"))
			return
		}
		author, err := q.CreateAuthor(r.Context(), authorParams)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Failed to create author")
			log.Debug("Error:", err)
			render.JSON(w, r, resp.Error("Failed to create author"))
			return
		}
		render.JSON(w, r, author)
	}
}
