package main

import (
	"fmt"
	"os"

	"github.com/adeaver/microservices-project/alphavantage"
)

const API_KEY_ENVIRONMENT_VARIABLE = "ALPHA_VANTAGE_KEY"

func main() {
	client := alphavantage.NewClient(os.Getenv(API_KEY_ENVIRONMENT_VARIABLE))
	resp, err := client.GetTimeSeries(alphavantage.GetTimeSeriesInput{
		Function:   alphavantage.FunctionTypeDaily,
		OutputSize: alphavantage.OutputSizeCompact,
		DataType:   alphavantage.DataTypeCSV,
		Symbol:     "GOOG",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(*resp)
}
