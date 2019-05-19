package alphavantage

import "time"

type EquitySnapshot struct {
	Time            time.Time `json:"time"`
	OpenPriceCents  *int64    `json:"open_price_cents,omitempty"`
	HighPriceCents  *int64    `json:"high_price_cents,omitempty"`
	LowPriceCents   *int64    `json:"low_price_cents,omitempty"`
	ClosePriceCents *int64    `json:"close_price_cents,omitempty"`
	VolumeShares    *int64    `json:"volume_shares,omitempty"`
}
