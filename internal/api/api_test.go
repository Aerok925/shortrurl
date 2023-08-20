package api

import (
	"bytes"
	"github.com/Aerok925/shortrurl/internal/app"
	"github.com/Aerok925/shortrurl/internal/inmemory"
	"github.com/Aerok925/shortrurl/internal/reducing"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHandlerURL(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	cache := inmemory.New()
	r := reducing.New()
	service := app.New(cache, r, logger, "localhost:8080")
	a := New(service, "localhost", ":8080")
	a.Rout()

	serverAddres := "http://localhost:8080"

	t.Run("shorten", func(t *testing.T) {
		w := httptest.NewRecorder()
		h := http.HandlerFunc(a.handlerCreateUrl)

		originalURL := "https://steamcommunity.com/market/listings/570/Elder%20The%20Defense%20Season%202%20War%20Dog?l=russian"

		req, err := http.NewRequest(http.MethodPost, serverAddres+"/", bytes.NewReader([]byte(originalURL)))
		h(w, req)
		resp := w.Result()

		assert.NoError(t, err, "ошибка при попытке создать сокращенную ссылку")
		body, _ := io.ReadAll(resp.Body)
		shortenUrl := string(body)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		_, err = url.Parse(shortenUrl)
		assert.NoError(t, err, "Не получилось распарсить url")
	})

	t.Run("create exist url", func(t *testing.T) {
		w := httptest.NewRecorder()
		h := http.HandlerFunc(a.handlerCreateUrl)
		originalURL := "https://steamcommunity.com/market/listings/570/Elder%20The%20Defense%20Season%202%20War%20Dog?l=russian"
		req, err := http.NewRequest(http.MethodPost, serverAddres+"/", bytes.NewReader([]byte(originalURL)))
		h(w, req)
		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		shortenUrl := string(body)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "в ответе был указан не тот статус код")
		_, err = url.Parse(shortenUrl)
		assert.NoError(t, err, "Не получилось распарсить url")
	})

	t.Run("empty url", func(t *testing.T) {
		w := httptest.NewRecorder()
		h := http.HandlerFunc(a.handlerCreateUrl)
		originalURL := ""
		req, err := http.NewRequest(http.MethodPost, serverAddres+"/", bytes.NewReader([]byte(originalURL)))
		h(w, req)
		resp := w.Result()
		assert.NoError(t, err, "ошибка при попытке создать сокращенную ссылку")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "в ответе был указан не тот статус код")
	})

}
