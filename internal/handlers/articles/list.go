package articles

import (
	"log/slog"
	"net/http"
	"strconv"

	db "github.com/flareopti/deeptruth/internal/db/sqlc"
	"github.com/flareopti/deeptruth/internal/lib/api/resp"
	"github.com/go-chi/render"
)

// ListAccounts godoc
//
//	@Summary		List articles
//	@Description	Get articles with pagination
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Param			page query int true "Page number"
//	@Param			per_page query int true "Articles per page"
//	@Success		200	{array}		db.Article
//	@Failure		400	{object}	resp.Response
//	@Failure		500	{object}	resp.Response
//	@Router			/api/articles [get]
func List(log *slog.Logger, q db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query_page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("no 'page' query parameter declared"))
			return
		}
		per_page, err := strconv.Atoi(r.URL.Query().Get("per_page"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("no 'per_page' query parameter declared"))
			return
		}

		articles, err := q.ListArticles(r.Context(), db.ListArticlesParams{
			Limit:  int32(query_page*per_page + per_page),
			Offset: int32(query_page * per_page),
		})
		if len(articles) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("no articles for you buddy"))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to fetch articles"))
			return
		}
		render.JSON(w, r, articles)
	}
}
