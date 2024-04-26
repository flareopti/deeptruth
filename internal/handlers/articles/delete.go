package articles

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

// Delete article
// @Summary Delete an article
// @Description Delete an article using its ID
// @Tags articles
// @Produce json
// @Param articleID path int true "Article ID"
// @Success 200 {object} resp.Response
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/articles/{articleID} [delete]
func Delete(log *slog.Logger, q db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "articleID")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("no article id provided"))
			return
		}
		id_int, err := strconv.Atoi(id)
		if err != nil {
			log.Error("Failed to convert article id to int")
			log.Debug("Error", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to convert article id to int"))
			return
		}
		err = q.DeleteArticle(r.Context(), int64(id_int))
		if err != nil {
			log.Error("Failed to delete article")
			log.Debug("Error", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to delete article"))
			return
		}
		render.JSON(w, r, resp.OK())
	}
}
