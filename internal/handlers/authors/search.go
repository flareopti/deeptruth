package authors

import (
	"log/slog"
	"net/http"

	db "github.com/flareopti/deeptruth/internal/db/sqlc"
)

func Search(log *slog.Logger, q db.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
