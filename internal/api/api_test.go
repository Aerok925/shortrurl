package api

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestHandlerURL(t *testing.T) {
	serverAddres := "http://localhost:8080"
	httpc := resty.New().
		SetBaseURL(serverAddres).
		SetRedirectPolicy(resty.NoRedirectPolicy())

	t.Run("shorten", func(t *testing.T) {
		originalURL := "https://steamcommunity.com/market/listings/570/Elder%20The%20Defense%20Season%202%20War%20Dog?l=russian"
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		req := httpc.R().
			SetContext(ctx).
			SetBody(originalURL)
		resp, err := req.Post("/")
		assert.NoError(t, err, "ошибка при попытке создать сокращенную ссылку")
		shortenUrl := string(resp.Body())
		assert.Equal(t, http.StatusCreated, resp.StatusCode())
		_, err = url.Parse(shortenUrl)
		assert.NoError(t, err, "Не получилось распарсить url")
	})

	t.Run("create exist url", func(t *testing.T) {
		originalURL := "https://steamcommunity.com/market/listings/570/Elder%20The%20Defense%20Season%202%20War%20Dog?l=russian"
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		req := httpc.R().
			SetContext(ctx).
			SetBody(originalURL)
		resp, err := req.Post("/")
		assert.NoError(t, err, "ошибка при попытке создать сокращенную ссылку")
		shortenUrl := string(resp.Body())
		assert.Equal(t, http.StatusOK, resp.StatusCode(), "в ответе был указан не тот статус код")
		_, err = url.Parse(shortenUrl)
		assert.NoError(t, err, "Не получилось распарсить url")
	})

	t.Run("empty url", func(t *testing.T) {
		originalURL := ""
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		req := httpc.R().
			SetContext(ctx).
			SetBody(originalURL)
		resp, err := req.Post("/")
		assert.NoError(t, err, "ошибка при попытке создать сокращенную ссылку")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode(), "в ответе был указан не тот статус код")
	})

	t.Run("get", func(t *testing.T) {
		originalURL := "https://account.jetbrains.com/licenses"
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		req := httpc.R().
			SetContext(ctx).
			SetBody(originalURL)
		resp, err := req.Post("/")
		assert.NoError(t, err, "ошибка при попытке создать сокращенную ссылку")
		shortenUrl := string(resp.Body())
		assert.Equal(t, http.StatusCreated, resp.StatusCode())
		_, err = url.Parse(shortenUrl)
		assert.NoError(t, err, "Не получилось распарсить url")

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		req = httpc.R().
			SetContext(ctx)
		resp, err = req.Get(shortenUrl)
		assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode())
		assert.Equal(t, originalURL, resp.Header().Get("Location"))
	})
}
