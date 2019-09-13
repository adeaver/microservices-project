package httpservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
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
		panic(fmt.Sprintf("unrecongized route method: %s", r))
	}
}

type Route struct {
	Endpoint string
	Method   RouteMethod
	Func     func(w http.ResponseWriter, r *http.Request) (*Response, error)
}

type Response struct {
	Payload    []byte
	StatusCode int
}

type InjectedDataUtils struct {
	Database *sqlx.DB
	Channel  *amqp.Channel
}

type withInjectedDataUtilsHandler func(dataUtils InjectedDataUtils, w http.ResponseWriter, r *http.Request) (*Response, error)

func withInjectedDataUtils(dataUtils InjectedDataUtils, f withInjectedDataUtilsHandler) func(http.ResponseWriter, *http.Request) (*Response, error) {
	return func(w http.ResponseWriter, r *http.Request) (*Response, error) {
		return f(dataUtils, w, r)
	}
}

type withDBHandler func(db *sqlx.DB, w http.ResponseWriter, r *http.Request) (*Response, error)

func WithDB(db *sqlx.DB, f withDBHandler) func(http.ResponseWriter, *http.Request) (*Response, error) {
	return func(w http.ResponseWriter, r *http.Request) (*Response, error) {
		return f(db, w, r)
	}
}

func registerRoute(f func(http.ResponseWriter, *http.Request) (*Response, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := f(w, r)
		if err != nil {
			sendErrorResponse(w, err)
			return
		}
		if payload != nil {
			sendJSONResponse(w, *payload)
		}
	}
}

func MakeOKResponse(v interface{}) *Response {
	response, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return &Response{
		Payload:    response,
		StatusCode: http.StatusOK,
	}
}

func sendJSONResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	w.Write(response.Payload)
}

func sendErrorResponse(w http.ResponseWriter, err error) {
	errorMap := map[string]string{"error": err.Error()}
	response, err := json.Marshal(errorMap)
	if err != nil {
		panic(err)
	}
	sendJSONResponse(w, Response{
		Payload:    response,
		StatusCode: http.StatusBadRequest,
	})
}
