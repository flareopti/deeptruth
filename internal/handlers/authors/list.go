package authors

import (
	"log/slog"
	"net/http"
	"strconv"

	db "github.com/flareopti/deeptruth/internal/db/sqlc"
	"github.com/flareopti/deeptruth/internal/lib/api/resp"
	"github.com/go-chi/render"
)

// ListAuthors godoc
//
//	@Summary		List authors
//	@Description	Get authors with pagination
//	@Tags			authors
//	@Accept			json
//	@Produce		json
//	@Param			page query int true "Page number"
//	@Param			per_page query int true "Authors per page"
//	@Success		200	{array}		db.Author
//	@Failure		400	{object}	resp.Response
//	@Failure		404	{object}	resp.Response
//	@Failure		500	{object}	resp.Response
//	@Router			/api/authors [get]
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

		authors, err := q.ListAuthors(r.Context(), db.ListAuthorsParams{
			Limit:  int32(query_page*per_page + per_page),
			Offset: int32(query_page * per_page),
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to fetch authors"))
			return
		}
		if len(authors) == 0 {
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error("no authors found for this query"))
			return
		}
		render.JSON(w, r, authors)
	}
}
