package static

import (
	"net/http"
	"os"
	"path/filepath"
)

type staticHandler struct {
	staticPath string
	indexPath  string
}

func (h staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticPath, r.URL.Path)

	stat, err := os.Stat(path)
	if os.IsNotExist(err) || stat.IsDir() {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func New(staticPath, indexPath string) staticHandler {
	return staticHandler{
		staticPath: staticPath,
		indexPath:  indexPath,
	}
}
