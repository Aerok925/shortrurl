package api

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type service interface {
	GetURL(id string) (string, error)
	CreateOrUpdateNewURL(value string) (string, bool, error)
}

type API struct {
	s        service
	r        *chi.Mux
	hostname string
	logger   *zap.Logger
}

func New(s service, hostname string, logger *zap.Logger) *API {
	r := chi.NewMux()

	return &API{
		s:        s,
		r:        r,
		hostname: hostname,
		logger:   logger,
	}
}

func (api *API) Rout() {
	api.r.Use(api.logging)
	api.r.Get("/{id}", api.handlerGetURL)

	api.r.Post("/", api.handlerCreateURL)
}

func (api *API) Start() error {
	api.logger.Info("Server start in " + api.hostname)
	return http.ListenAndServe(api.hostname, api.r)
}

func (api *API) handlerGetURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url, err := api.s.GetURL(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (api *API) handlerCreateURL(w http.ResponseWriter, r *http.Request) {
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := string(bodyData)
	if len(url) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newURL, b, err := api.s.CreateOrUpdateNewURL(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	if b {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(newURL))
}
