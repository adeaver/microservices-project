package httpservice

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	Routes     []Route
	ListenPort string
}

func (s *Service) Start() {
	r := mux.NewRouter()
	for _, route := range s.Routes {
		f := registerRoute(route.Func)
		r.HandleFunc(route.Endpoint, f).Methods(route.Method.Str())
	}
	http.ListenAndServe(s.ListenPort, r)
}

type RouteMethod string

// TODO: add more route methods if I ever need them
const (
	RouteMethodGET  RouteMethod = "GET"
	RouteMethodPOST RouteMethod = "POST"
)

func (r RouteMethod) Str() string {
	switch r {
	case RouteMethodGET:
		return "GET"
	case RouteMethodPOST:
		return "POST"
	default:
		panic("unrecongized route method: %s", r)
	}
}

type Route struct {
	Endpoint string
	Method   RouteMethod
	Func     func(w http.ResponseWriter, r *http.Request) (interface{}, error)
}

type WithDBHandler func(db *sqlx.DB, w http.ResponseWriter, r *http.Request) (interface{}, error)

func WithDB(db *sqlx.DB, f WithDBHandler) func(http.ResponseWriter, *http.Request) (interface{}, error) {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		return f(db, w, r)
	}
}

func registerRoute(f func(http.ResponseWriter, *http.Request) (interface{}, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := f(w, r)
		if err != nil {
			makeErrorResponse(w, err)
			return
		}
		makeJSONResponse(w, http.StatusOK, payload)
	}
}

func makeJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func makeErrorResponse(w http.ResponseWriter, err error) {
	makeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}
