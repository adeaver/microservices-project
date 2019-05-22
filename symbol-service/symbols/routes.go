package symbols

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(r *mux.Router, db *sqlx.DB) {
	r.HandleFunc("/get_symbols_1", withDB(db, handleGetAllSymbols)).Methods("GET")
	r.HandleFunc("/insert_symbol_1", withDB(db, handleInsertSymbol)).Methods("POST")
}

// TODO: put this in a middleware
type withDBHandler func(db *sqlx.DB, w http.ResponseWriter, r *http.Request)

func withDB(db *sqlx.DB, f withDBHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(db, w, r)
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

func handleGetAllSymbols(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	var symbols []Symbol
	db.Select(&symbols, "SELECT * FROM symbols")
	makeJSONResponse(w, http.StatusOK, symbols)
}

func handleInsertSymbol(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	var s Symbol
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		makeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	fmt.Println("symbol %+v", s)
	tx, err := db.Beginx()
	if err != nil {
		makeJSONResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	// TODO: this needs a wrapper
	tx.MustExec("INSERT INTO symbols (name, symbol, market_capitalization, sector, industry, exchange) VALUES ($1, $2, $3, $4, $5, $6)", s.Name, s.Symbol, s.MarketCapitalization, s.Sector, s.Industry, s.Exchange)
	tx.Commit()
	makeJSONResponse(w, http.StatusOK, map[string]bool{"success": true})
}
