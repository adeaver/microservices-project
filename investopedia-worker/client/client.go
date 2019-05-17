package client

import (
	"fmt"
	"io/ioutil"
	"microservices-project/investopedia-worker/model"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	investopediaStockPricesURL = "https://www.investopedia.com/markets/api/partial/historical"

	investopediaQueryKeySymbol    = "Symbol"
	investopediaQueryKeyTimeFrame = "TimeFrame"
	investopediaQueryKeyStartDate = "StartDate"
	investopediaQueryKeyEndDate   = "EndDate"
	investopediaQueryKeyType      = "Type"

	investopediaDateFormat = "Jan 2, 2006"
)

type InvestopediaInput struct {
	EndDate   time.Time
	StartDate time.Time
	Symbol    string
	TimeFrame InvestopediaTimeFrame
	Type      InvestopediaDataType
}

type InvestopediaTimeFrame string

const (
	InvestopediaDailyTimeFrame InvestopediaTimeFrame = "Daily"
)

type InvestopediaDataType string

const (
	InvestopediaHistoricalPricesType InvestopediaDataType = "Historical Prices"
)

func GetStockPricesForSymbol(symbol string, startDate time.Time, endDate time.Time) ([]model.TradingDay, error) {
	req, err := makeStockPriceRequest(InvestopediaInput{
		EndDate:   endDate,
		StartDate: startDate,
		Symbol:    symbol,
		TimeFrame: InvestopediaDailyTimeFrame,
		Type:      InvestopediaHistoricalPricesType,
	})
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		statusMessage := fmt.Sprintf("Investopedia status code not okay, got %v", resp.StatusCode)
		return nil, fmt.Errorf(statusMessage)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	str := string(bodyBytes)
	return parseResponse(str)
}

func makeStockPriceRequest(input InvestopediaInput) (*http.Request, error) {
	u, err := url.Parse(investopediaStockPricesURL)
	if err != nil {
		return nil, err
	}
	query := makeQuery(input)
	u.RawQuery = query
	r, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func makeQuery(input InvestopediaInput) string {
	values := url.Values{}
	values.Set(investopediaQueryKeyType, string(input.Type))
	values.Set(investopediaQueryKeyEndDate, formatTime(input.EndDate))
	values.Set(investopediaQueryKeyStartDate, formatTime(input.StartDate))
	values.Set(investopediaQueryKeyTimeFrame, string(input.TimeFrame))
	values.Set(investopediaQueryKeySymbol, input.Symbol)
	return strings.Replace(values.Encode(), "%2C", ",", -1)
}

func formatTime(t time.Time) string {
	return t.Format(investopediaDateFormat)
}
