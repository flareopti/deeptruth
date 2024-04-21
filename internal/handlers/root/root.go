package root

import (
	"log/slog"
	"net/http"

	"github.com/flareopti/deeptruth/internal/lib/api/resp"
	"github.com/go-chi/render"
)

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, resp.OK())
	}
}
