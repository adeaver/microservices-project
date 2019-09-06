package quotes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/adeaver/microservices-project/util/alphavantage"
	"github.com/adeaver/microservices-project/util/httpservice"
	"github.com/jmoiron/sqlx"
)

func MakeRouteDefinitions(db *sqlx.DB) []httpservice.Route {
	return []httpservice.Route{
		{
			Endpoint: "/internal/get_all_quotes_for_top_symbols_1",
			Func:     httpservice.WithDB(db, handleGetAllQuotesForTopSymbols),
			Method:   httpservice.RouteMethodGET,
		},
		{
			Endpoint: "/get_all_quotes_for_symbol_1",
			Func:     httpservice.WithDB(db, handleGetAllQuotesForSymbol),
			Method:   httpservice.RouteMethodPOST,
		},
	}
}

const (
	alphavantageAPIKey = "ALPHA_VANTAGE_KEY"
	symbolServiceHost  = "SYMBOL_SERVICE_HOST"
	symbolServicePort  = "SYMBOL_SERVICE_PORT"
)

// There are all sorts of ways this could be faster
// 1) Parallelize the collection and the insertion
// 2) Dispatch these to workers
func handleGetAllQuotesForTopSymbols(db *sqlx.DB, w http.ResponseWriter, r *http.Request) (*httpservice.Response, error) {
	alphavantageAPIKey := os.Getenv(alphavantageAPIKey)
	client := alphavantage.NewClient(alphavantageAPIKey)
	return injectSymbols(func(symbols []Symbol) (*httpservice.Response, error) {
		for _, s := range symbols[:10] {
			data, err := client.GetTimeSeries(alphavantage.GetTimeSeriesInput{
				Function:   alphavantage.FunctionTypeDaily,
				OutputSize: alphavantage.OutputSizeCompact,
				DataType:   alphavantage.DataTypeCSV,
				Symbol:     s.Symbol,
			})
			if err != nil {
				return nil, err
			}
			var toInsert []Quote
			for _, d := range data {
				if d == nil {
					continue
				}
				toInsert = append(toInsert, quoteFromEquitySnapshot(s, *d))
			}
			tx, err := db.Beginx()
			if err != nil {
				return nil, err
			}
			tx.NamedExec("INSERT INTO quotes (symbol, time, open_price_cents, high_price_cents, low_price_cents, close_price_cents, volume_shares) VALUES (:symbol, :time, :open_price_cents, :high_price_cents, :low_price_cents, :close_price_cents, :volume_shares)", toInsert)
			tx.Commit()
		}
		return httpservice.MakeOKResponse(map[string]bool{"success": true}), nil
	})
}

type getAllQuotesForSymbolRequest struct {
	Symbol string `json"`
}

func handleGetAllQuotesForSymbol(db *sqlx.DB, w http.ResponseWriter, r *http.Request) (*httpservice.Response, error) {
	var req getAllQuotesForSymbolRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		return nil, err
	}
	var quotes []Quote
	db.Select(&quotes, fmt.Sprintf("SELECT * FROM quotes WHERE symbol='%s'", req.Symbol))
	return httpservice.MakeOKResponse(quotes), nil
}

func injectSymbols(f func(symbols []Symbol) (*httpservice.Response, error)) (*httpservice.Response, error) {
	symbolServiceHost := os.Getenv(symbolServiceHost)
	symbolServicePort := os.Getenv(symbolServicePort)
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%v:%v/get_top_symbols_1", symbolServiceHost, symbolServicePort), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status not okay")
	}
	var symbols []Symbol
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&symbols); err != nil {
		return nil, err
	}
	return f(symbols)
}
