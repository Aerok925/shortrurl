package api

import (
	"github.com/gorilla/mux"
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
	r        *mux.Router
	hostname string
	logger   *zap.Logger
}

func New(s service, hostname string, logger *zap.Logger) *API {
	r := mux.NewRouter()

	return &API{
		s:        s,
		r:        r,
		hostname: hostname,
		logger:   logger,
	}
}

func (a *API) Rout() {
	a.r.Name("GetURL").
		Path("/{id}").
		HandlerFunc(a.handlerGetURL).
		Methods(http.MethodGet)

	a.r.Name("createURL").
		Path("/").
		HandlerFunc(a.handlerCreateURL).
		Methods(http.MethodPost)
}

func (a *API) Start() error {
	a.logger.Info("Server start in " + a.hostname)
	return http.ListenAndServe(a.hostname, a.r)
}

func (a *API) handlerGetURL(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id, ok := v["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url, err := a.s.GetURL(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (a *API) handlerCreateURL(w http.ResponseWriter, r *http.Request) {
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
	newURL, b, err := a.s.CreateOrUpdateNewURL(url)
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
