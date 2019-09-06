package quotes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/adeaver/microservices-project/util/alphavantage"
	"github.com/adeaver/microservices-project/util/httpservice"
)

func MakeRouteDefinitions() []httpservice.Route {
	return []httpservice.Route{
		{
			Endpoint: "/internal/get_quotes_for_top_symbols_1",
			Func:     handleGetQuotesForTopSymbols,
			Method:   httpservice.RouteMethodGET,
		},
	}
}

const (
	alphavantageAPIKey = "ALPHA_VANTAGE_KEY"
	symbolServiceHost  = "SYMBOL_SERVICE_HOST"
	symbolServicePort  = "SYMBOL_SERVICE_PORT"
)

func handleGetQuotesForTopSymbols(w http.ResponseWriter, r *http.Request) (*httpservice.Response, error) {
	alphavantageAPIKey := os.Getenv(alphavantageAPIKey)
	client := alphavantage.NewClient(alphavantageAPIKey)
	_, err := client.GetTimeSeries(alphavantage.GetTimeSeriesInput{
		Function:   alphavantage.FunctionTypeDaily,
		OutputSize: alphavantage.OutputSizeCompact,
		DataType:   alphavantage.DataTypeCSV,
		Symbol:     "GOOG",
	})
	if err != nil {
		return nil, err
	}
	symbols, err := getSymbols()
	if err != nil {
		return nil, err
	}
	return httpservice.MakeOKResponse(symbols), nil
}

func getSymbols() ([]Symbol, error) {
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
	return symbols, nil
}
