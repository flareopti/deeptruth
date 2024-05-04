package articles

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

	db "github.com/flareopti/deeptruth/internal/db/sqlc"
	"github.com/flareopti/deeptruth/internal/lib/api/resp"
	"github.com/flareopti/deeptruth/internal/lib/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Estimate fuckness
// @Summary Estimate fuckness
// @Description Estimate fuckness
// @Tags articles
// @Procuce json
// @Param articleID path int true "Article ID"
// @Success 200 {object} db.UpdateArticleRatingParams
// @Failure 400 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Router /api/articles/{articleID} [post]
func Estimate(log *slog.Logger, q db.Querier) http.HandlerFunc {
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
		article, err := q.GetArticle(r.Context(), int64(id_int))
		if err != nil {
			log.Error("Failed to get article", sl.Err(err))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error("failed to get this article"))
			return
		}
		newRating := getOpenaiRating(log, article)
		updatedParams := db.UpdateArticleRatingParams{
			ID:     article.ID,
			Rating: newRating,
		}
		_, err = q.UpdateArticleRating(r.Context(), updatedParams)
		if err != nil {
			log.Error("Failed to update article rating", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to update article rating"))
			return
		}
		render.JSON(w, r, updatedParams)
	}
}

func getOpenaiRating(log *slog.Logger, article db.Article) int32 {
	query := fmt.Sprintf(`[{"role": "system", "content": "Ты пытаешься понять является ли следующая информация правдой или ложью, выдай оценку этой новости от 0 до 3, где 0 - ужасающе неправдивая информация, 3 - гарантированная правда. Ответ выдавай в четком формате rating:{число}, например rating:3"}, {"role": "user", "content":"%s\n%s\n"}]`, article.Title, article.Content)
	body := fmt.Sprintf(`{
		  "model":       "gpt-3.5-turbo",
		  "messages":    %s,
		  "max_tokens":  1024,
		  "temperature": 0.8,
		  "stream":      false
		  }`, query)

	openai_r, _ := http.NewRequest("POST", "https://api.naga.ac/v1/chat/completions", bytes.NewBuffer([]byte(body)))
	openai_r.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer dZvWPb8Pxbv5P7zAQHegqVyA5yBiCDXIjMLrDLKZVQ4"},
	}

	openai_response, err := http.DefaultClient.Do(openai_r)
	if err != nil {
		log.Error("Openai service fucked up!", sl.Err(err))
		return -1
	}
	defer openai_response.Body.Close()
	response, err := io.ReadAll(openai_response.Body)
	if err != nil {
		log.Error("Failed to read response", sl.Err(err))
		return -1
	}
	log.Debug(string(response))
	re, err := regexp.Compile(`rating:(\d+)`)
	if err != nil {
		log.Error("Failed to find rating", sl.Err(err))
		return -1
	}
	rating_string := re.FindStringSubmatch(string(response))
	rating, err := strconv.Atoi(rating_string[1])
	if err != nil {
		log.Error("Failed to convert", sl.Err(err))
		return -1
	}

	return int32(rating)
}
