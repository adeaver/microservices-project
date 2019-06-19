package main

import (
	"encoding/json"
	"fmt"
	"microservices-project/timeseries-worker/model"
	"net/http"
	"os"

	"github.com/adeaver/microservices-project/util/alphavantage"
)

const (
	API_KEY_ENVIRONMENT_VARIABLE = "ALPHA_VANTAGE_KEY"
	SYMBOL_SERVICE_HOST          = "SYMBOL_SERVICE_HOST"
	SYMBOL_SERVICE_PORT          = "SYMBOL_SERVICE_PORT"
)

func main() {

	client := alphavantage.NewClient(os.Getenv(API_KEY_ENVIRONMENT_VARIABLE))
	_, err := client.GetTimeSeries(alphavantage.GetTimeSeriesInput{
		Function:   alphavantage.FunctionTypeDaily,
		OutputSize: alphavantage.OutputSizeCompact,
		DataType:   alphavantage.DataTypeCSV,
		Symbol:     "GOOG",
	})
	if err != nil {
		panic(err)
	}
	symbols, err := getSymbols()
	if err != nil {
		panic(err)
	}
	fmt.Println(symbols)
}

func getSymbols() ([]model.Symbol, error) {
	symbolServiceHost := os.Getenv(SYMBOL_SERVICE_HOST)
	symbolServicePort := os.Getenv(SYMBOL_SERVICE_PORT)
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
	var symbols []model.Symbol
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&symbols); err != nil {
		return nil, err
	}
	return symbols, nil
}
