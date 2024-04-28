package articles

import (
	"log/slog"
	"net/http"

	db "github.com/flareopti/deeptruth/internal/db/sqlc"
	"github.com/flareopti/deeptruth/internal/lib/api/resp"
	"github.com/flareopti/deeptruth/internal/lib/sl"
	"github.com/go-chi/render"
)

// Article create
// @Summary Create an article
// @Description Create an article
// @Tags articles
// @Accept json
// @Procuce json
// @Param article body db.CreateArticleParams true "Article to create"
// @Success 200 {object} db.Article
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/articles [post]
func Create(log *slog.Logger, q db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var article_params db.CreateArticleParams
		err := render.DecodeJSON(r.Body, &article_params)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("Failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("Failed to decode request"))
			return
		}
		article, err := q.CreateArticle(r.Context(), article_params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Failed to create article", sl.Err(err))
			render.JSON(w, r, resp.Error("Failed to create article"))
			return
		}
		render.JSON(w, r, article)
	}
}
