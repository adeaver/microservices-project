package main

import (
	"microservices-project/quote-service/quotes"

	"github.com/adeaver/microservices-project/util/httpservice"
)

func main() {
	s := httpservice.Service{
		Routes:     quotes.MakeRouteDefinitions(),
		ListenPort: ":3050",
	}
	s.Start()
}
