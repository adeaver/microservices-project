package alphavantage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const alphaVantageURL string = "https://www.alphavantage.co/query"

type FunctionType string

const (
	FunctionTypeDaily FunctionType = "TIME_SERIES_DAILY"
)

func (f *FunctionType) String() string {
	return string(*f)
}

type OutputSize string

const (
	OutputSizeCompact OutputSize = "compact"
)

func (o *OutputSize) String() string {
	return string(*o)
}

type DataType string

const (
	DataTypeJSON DataType = "json"
	DataTypeCSV  DataType = "csv"
)

func (d *DataType) String() string {
	return string(*d)
}

type alphaVantageClient struct {
	apiKey string
}

func NewClient(apiKey string) *alphaVantageClient {
	return &alphaVantageClient{
		apiKey: apiKey,
	}
}

type GetTimeSeriesInput struct {
	Function   FunctionType
	OutputSize OutputSize
	DataType   DataType
	Symbol     string
}

func (c *alphaVantageClient) GetTimeSeries(input GetTimeSeriesInput) ([]*EquitySnapshot, error) {
	url, err := c.makeURLFromInput(input)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", url.String(), nil)
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
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	str := string(bodyBytes)
	switch input.DataType {
	case DataTypeCSV:
		fmt.Println(str)
		return parseResponseCSV(str)
	case DataTypeJSON:
		return nil, fmt.Errorf("json type not implemented")
	default:
		return nil, fmt.Errorf("unsupported data type")
	}
}

func (c *alphaVantageClient) makeURLFromInput(input GetTimeSeriesInput) (*url.URL, error) {
	u, err := url.Parse(alphaVantageURL)
	if err != nil {
		return nil, err
	}
	values := url.Values{}
	values.Set("function", input.Function.String())
	values.Set("symbol", input.Symbol)
	values.Set("outputsize", input.OutputSize.String())
	values.Set("datatype", input.DataType.String())
	values.Set("apikey", c.apiKey)
	u.RawQuery = values.Encode()
	return u, nil
}
