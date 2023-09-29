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
	service := app.New(cache, r, logger, "http://localhost:8080")
	a := New(service, "localhost:8080", logger)
	a.Rout()

	serverAddres := "http://localhost:8080"

	t.Run("shorten", func(t *testing.T) {
		w := httptest.NewRecorder()
		h := http.HandlerFunc(a.handlerCreateURL)

		originalURL := "https://steamcommunity.com/market/listings/570/Elder%20The%20Defense%20Season%202%20War%20Dog?l=russian"

		req, err := http.NewRequest(http.MethodPost, serverAddres+"/", bytes.NewReader([]byte(originalURL)))
		h(w, req)
		resp := w.Result()
		assert.NoError(t, err, "ошибка при попытке создать сокращенную ссылку")
		body, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		shortenURL := string(body)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		_, err = url.Parse(shortenURL)
		assert.NoError(t, err, "Не получилось распарсить url")
	})

	t.Run("create exist url", func(t *testing.T) {
		w := httptest.NewRecorder()
		h := http.HandlerFunc(a.handlerCreateURL)
		originalURL := "https://steamcommunity.com/market/listings/570/Elder%20The%20Defense%20Season%202%20War%20Dog?l=russian"
		req, err := http.NewRequest(http.MethodPost, serverAddres+"/", bytes.NewReader([]byte(originalURL)))
		assert.NoError(t, err)
		h(w, req)
		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		shortenURL := string(body)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "в ответе был указан не тот статус код")
		_, err = url.Parse(shortenURL)
		assert.NoError(t, err, "Не получилось распарсить url")
	})

	t.Run("empty url", func(t *testing.T) {
		w := httptest.NewRecorder()
		h := http.HandlerFunc(a.handlerCreateURL)
		originalURL := ""
		req, err := http.NewRequest(http.MethodPost, serverAddres+"/", bytes.NewReader([]byte(originalURL)))
		h(w, req)
		resp := w.Result()
		defer resp.Body.Close()
		assert.NoError(t, err, "ошибка при попытке создать сокращенную ссылку")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "в ответе был указан не тот статус код")
	})

}
