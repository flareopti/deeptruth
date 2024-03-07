package root

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, struct {
			Result string
		}{
			Result: "Hello motherfuckers!",
		})
	}
}
