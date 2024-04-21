package articles

import (
	"log/slog"
	"net/http"
	"strconv"

	db "github.com/flareopti/deeptruth/internal/db/sqlc"
	"github.com/flareopti/deeptruth/internal/lib/api/resp"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Get one article
// @Summary Get one article
// @Description Get one article using URL
// @Tags articles
// @Procuce json
// @Param articleID path int true "Article ID"
// @Success 200 {object} db.Article
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/articles/{articleID} [get]
func Get(log *slog.Logger, q db.Querier) http.HandlerFunc {
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
			log.Debug("Error:", err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to convert article id to int"))
			return
		}
		article, err := q.GetArticle(r.Context(), int64(id_int))
		if err != nil {
			log.Error("Failed to fetch article")
			log.Debug("Error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to fetch article"))
			return
		}
		render.JSON(w, r, article)
	}
}
