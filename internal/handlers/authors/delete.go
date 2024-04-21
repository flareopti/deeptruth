package authors

import (
	"log/slog"
	"net/http"
	"strconv"

	db "github.com/flareopti/deeptruth/internal/db/sqlc"
	"github.com/flareopti/deeptruth/internal/lib/api/resp"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Delete author
// @Summary Delete an author
// @Description Delete an author using their ID
// @Tags authors
// @Produce json
// @Param authorID path int true "Author ID"
// @Success 200 {object} resp.Response
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/authors/{authorID} [delete]
func Delete(log *slog.Logger, q db.Querier) http.HandlerFunc {
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
			log.Debug("Error:", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to convert author id to int"))
			return
		}
		err = q.DeleteAuthor(r.Context(), int64(id_int))
		if err != nil {
			log.Error("Failed to delete author")
			log.Debug("Error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete author"))
			return
		}
		render.JSON(w, r, resp.OK())
	}
}
