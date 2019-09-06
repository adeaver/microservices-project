package symbols

import (
	"encoding/json"
	"net/http"

	"github.com/adeaver/microservices-project/util/httpservice"

	"github.com/jmoiron/sqlx"
)

func MakeRouteDefinitions(db *sqlx.DB) []httpservice.Route {
	return []httpservice.Route{
		{
			Endpoint: "/get_symbols_1",
			Func:     httpservice.WithDB(db, handleGetAllSymbols),
			Method:   httpservice.RouteMethodGET,
		},
		{
			Endpoint: "/get_top_symbols_1",
			Func:     httpservice.WithDB(db, handleGetTopSymbolsByMarketCap),
			Method:   httpservice.RouteMethodGET,
		},
		{
			Endpoint: "/insert_symbol_1",
			Func:     httpservice.WithDB(db, handleInsertSymbol),
			Method:   httpservice.RouteMethodPOST,
		},
	}
}

func handleGetAllSymbols(db *sqlx.DB, w http.ResponseWriter, r *http.Request) (*httpservice.Response, error) {
	var symbols []Symbol
	db.Select(&symbols, "SELECT * FROM symbols")
	return httpservice.MakeOKResponse(symbols), nil
}

func handleInsertSymbol(db *sqlx.DB, w http.ResponseWriter, r *http.Request) (*httpservice.Response, error) {
	var s Symbol
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&s); err != nil {
		return nil, err
	}
	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	// TODO: this needs a wrapper
	tx.MustExec("INSERT INTO symbols (name, symbol, market_capitalization, sector, industry, exchange) VALUES ($1, $2, $3, $4, $5, $6)", s.Name, s.Symbol, s.MarketCapitalization, s.Sector, s.Industry, s.Exchange)
	tx.Commit()
	return httpservice.MakeOKResponse(map[string]bool{"success": true}), nil
}
func handleGetTopSymbolsByMarketCap(db *sqlx.DB, w http.ResponseWriter, r *http.Request) (*httpservice.Response, error) {
	var symbols []Symbol
	db.Select(&symbols, "SELECT * FROM symbols ORDER BY market_capitalization DESC LIMIT 500")
	return httpservice.MakeOKResponse(symbols), nil
}
