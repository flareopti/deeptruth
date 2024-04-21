package static

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func New(log *slog.Logger, directory string) http.HandlerFunc {
	workDir, _ := os.Getwd()
	fullPath := http.Dir(filepath.Join(workDir, directory))

	return func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(fullPath))
		fs.ServeHTTP(w, r)
	}
}
