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

type UpdateRatingParam struct {
	Rating int
}

// Article rating update
// @Summary Update rating of an article
// @Description Update rating of an article
// @Tags articles
// @Accept json
// @Procuce json
// @Param articleID path int true "Article ID"
// @Param article body UpdateRatingParam true "New rating"
// @Success 200 {object} db.Article
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/articles/{articleID} [patch]
func UpdateRating(log *slog.Logger, q db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		id := chi.URLParam(r, "articleID")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("no article id provided"))
			return
		}
		id_int, err := strconv.Atoi(id)
		if err != nil {
			log.Error("Failed to convert article id to int", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to convert article id to int"))
			return
		}
		var rating UpdateRatingParam
		err = render.DecodeJSON(r.Body, &rating)
		article_params := db.UpdateArticleRatingParams{
			ID:     int64(id_int),
			Rating: int32(rating.Rating),
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("Failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("Failed to decode request"))
			return
		}
		article, err := q.UpdateArticleRating(r.Context(), article_params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Failed to update article", sl.Err(err))
			render.JSON(w, r, resp.Error("Failed to update article"))
			return
		}
		render.JSON(w, r, article)
	}
}
