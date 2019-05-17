package client

import (
	"fmt"
	"io"
	"microservices-project/investopedia-worker/model"
	"microservices-project/investopedia-worker/util"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func parseResponse(resp io.Reader) ([]model.TradingDay, error) {
	z := html.NewTokenizer(resp)
	collectTokens := false
	out := make([]model.TradingDay, 0)
	curr := make([]string, 0)
	for {
		t := z.Next()
		switch t {
		case html.StartTagToken:
			for _, attr := range z.Token().Attr {
				if attr.Key == "class" && attr.Val == "in-the-money" {
					collectTokens = true
					curr = make([]string, 0)
				}
			}
		case html.EndTagToken:
			if collectTokens && z.Token().Data == "tr" {
				tradingDay, err := makeTradingDay(curr)
				if err != nil {
					return nil, err
				}
				out = append(out, *tradingDay)
				collectTokens = false
			}
		case html.TextToken:
			text := strings.Trim(string(z.Text()), " \n\r")
			if collectTokens && len(text) > 0 {
				if strings.Contains(text, "Dividend") {
					curr = make([]string, 0)
					collectTokens = false
					continue
				}
				curr = append(curr, text)
			}
		case html.ErrorToken:
			if err := z.Err(); err != io.EOF {
				return nil, err
			}
			return out, nil
		}
	}
	return out, nil
}

func makeTradingDay(values []string) (*model.TradingDay, error) {
	if len(values) != 6 {
		return nil, fmt.Errorf("Incorrect length")
	}
	t, err := time.Parse(investopediaDateFormat, values[0])
	if err != nil {
		return nil, err
	}
	v, err := strconv.ParseInt(strings.Replace(values[5], ",", "", -1), 10, 64)
	if err != nil {
		return nil, err
	}
	return &TradingDay{
		Date:   t,
		Open:   util.MustParseCentsFromFloatString(values[1]),
		Close:  util.MustParseCentsFromFloatString(values[2]),
		High:   util.MustParseCentsFromFloatString(values[3]),
		Low:    util.MustParseCentsFromFloatString(values[4]),
		Volume: v,
	}, nil
}
