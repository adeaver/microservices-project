package symbols

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r *mux.Router, db *sqlx.DB) {
	r.HandleFunc("/symbols", handleGetAllSymbols(db))
}

func handleGetAllSymbols(db *sqlx.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var symbols []Symbol
		db.Select(&symbols, "SELECT * FROM symbols")
		bytes, err := json.Marshal(symbols)
		if err != nil {
			panic(err)
		}
		w.Write(bytes)
	}
}
