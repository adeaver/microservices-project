package quotes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/adeaver/microservices-project/util/httpservice"
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
)

func MakeRouteDefinitions(db *sqlx.DB, ch *amqp.Channel) []httpservice.Route {
	dataUtils := httpservice.InjectedDataUtils{
		Database: db,
		Channel:  ch,
	}
	return []httpservice.Route{
		{
			Endpoint: "/internal/get_all_quotes_for_top_symbols_1",
			Func:     httpservice.WithInjectedDataUtils(dataUtils, handleGetAllQuotesForTopSymbols),
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
	symbolServiceHost    = "SYMBOL_SERVICE_HOST"
	symbolServicePort    = "SYMBOL_SERVICE_PORT"
	rabbitMQSymbolsQueue = "RABBITMQ_SYMBOLS_QUEUE"
)

// There are all sorts of ways this could be faster
// 1) Parallelize the collection and the insertion
// 2) Dispatch these to workers
// 3) Batch insert into postgres
func handleGetAllQuotesForTopSymbols(dataUtils httpservice.InjectedDataUtils, w http.ResponseWriter, r *http.Request) (*httpservice.Response, error) {
	q, err := dataUtils.Channel.QueueDeclare(
		os.Getenv(rabbitMQSymbolsQueue),
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return injectSymbols(func(symbols []Symbol) (*httpservice.Response, error) {
		var errors []string
		for _, s := range symbols {
			msg, err := json.Marshal(s)
			if err != nil {
				errors = append(errors, fmt.Sprintf("error marshaling symbol: %s", err.Error()))
				continue
			}
			err = dataUtils.Channel.Publish(
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(msg),
				})
			if err != nil {
				errors = append(errors, fmt.Sprintf("error publishing message: %s", err.Error()))
				continue
			}
		}
		return httpservice.MakeOKResponse(map[string][]string{"errors": errors}), nil
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
	db.Select(&quotes, fmt.Sprintf("SELECT * FROM quotes WHERE symbol = '%s'", req.Symbol))
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
