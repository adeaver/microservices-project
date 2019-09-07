package quotes

import (
	"time"

	"github.com/adeaver/microservices-project/util/alphavantage"
)

type Symbol struct {
	ID                   *string `db:"id,omitempty" json:"id,omitempty"`
	Name                 string  `db:"name" json:"name"`
	Symbol               string  `db:"symbol" json:"symbol"`
	MarketCapitalization int64   `db:"market_capitalization" json:"market_capitalization"`
	Sector               *string `db:"sector,omitempty" json:"sector,omitempty"`
	Industry             *string `db:"industry,omitempty" json:"industry,omitempty"`
	Exchange             string  `db:"exchange" json:"exchange"`
}

type Quote struct {
	Symbol          string    `db:"symbol" json:"symbol"`
	Time            time.Time `db:"time" json:"time"`
	OpenPriceCents  *int64    `db:"open_price_cents" json:"open_price_cents"`
	HighPriceCents  *int64    `db:"high_price_cents" json:"high_price_cents"`
	LowPriceCents   *int64    `db:"low_price_cents" json:"low_price_cents"`
	ClosePriceCents *int64    `db:"close_price_cents" json:"close_price_cents"`
	VolumeShares    *int64    `db:"volume_shares" json:"volume_shares"`
}

func quoteFromEquitySnapshot(s Symbol, es alphavantage.EquitySnapshot) Quote {
	return Quote{
		Symbol:          s.Symbol,
		Time:            es.Time,
		OpenPriceCents:  es.OpenPriceCents,
		HighPriceCents:  es.HighPriceCents,
		LowPriceCents:   es.LowPriceCents,
		ClosePriceCents: es.ClosePriceCents,
		VolumeShares:    es.VolumeShares,
	}
}
