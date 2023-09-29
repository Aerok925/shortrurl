package api

import (
	"encoding/json"
	"github.com/Aerok925/shortrurl/internal/entities"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type service interface {
	GetURL(id string) (entities.ShortURL, error)
	CreateOrUpdateNewURL(shortURL entities.UnprocessedURL) (entities.ShortURL, error)
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
	api.r.Post("/api/shorten", api.handlerCreateURLJSON)
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
	w.Header().Set("Location", url.URL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (api *API) handlerCreateURL(w http.ResponseWriter, r *http.Request) {
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := entities.UnprocessedURL{
		URL: string(bodyData),
	}
	if len(url.URL) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newURL, err := api.s.CreateOrUpdateNewURL(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	if newURL.Create {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(newURL.URL))
}

func (api *API) handlerCreateURLJSON(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	var r entities.UnprocessedURL
	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newURL, err := api.s.CreateOrUpdateNewURL(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if newURL.Create {
		w.WriteHeader(http.StatusCreated)
	}
	resp, err := json.Marshal(newURL)
	if err != nil {
		return
	}
	w.Write(resp)
}
