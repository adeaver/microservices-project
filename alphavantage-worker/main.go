package main

import (
	"fmt"
	"os"

	"microservices-project/alphavantage-worker/timeseries"
)

const API_KEY_ENVIRONMENT_VARIABLE = "ALPHA_VANTAGE_KEY"

func main() {
	client := timeseries.NewClient(os.Getenv(API_KEY_ENVIRONMENT_VARIABLE))
	resp, err := client.GetTimeSeries(timeseries.GetTimeSeriesInput{
		Function:   timeseries.FunctionTypeDaily,
		OutputSize: timeseries.OutputSizeCompact,
		DataType:   timeseries.DataTypeCSV,
		Symbol:     "GOOG",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp)
}
