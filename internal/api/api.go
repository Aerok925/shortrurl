package api

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type service interface {
	GetUrl(id string) (string, error)
	CreateOrUpdateNewUrl(value string) (string, bool, error)
}

type Api struct {
	s        service
	r        *mux.Router
	hostname string
	port     string
}

func New(s service, hostname, port string) *Api {
	r := mux.NewRouter()

	return &Api{
		s:        s,
		r:        r,
		hostname: hostname,
		port:     port,
	}
}

func (a *Api) Rout() {
	a.r.Name("GetURL").
		Path("/{id}").
		HandlerFunc(a.handlerGetUrl).
		Methods(http.MethodGet)

	a.r.Name("createURL").
		Path("/").
		HandlerFunc(a.handlerCreateUrl).
		Methods(http.MethodPost)
}

func (a *Api) Start() error {
	return http.ListenAndServe(a.hostname+a.port, a.r)
}

func (a *Api) handlerGetUrl(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	id, ok := v["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url, err := a.s.GetUrl(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (a *Api) handlerCreateUrl(w http.ResponseWriter, r *http.Request) {
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := string(bodyData)
	if len(url) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newUrl, b, err := a.s.CreateOrUpdateNewUrl(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	if b {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(newUrl))

}
