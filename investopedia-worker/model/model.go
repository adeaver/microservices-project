package model

import "time"

type TradingDay struct {
	DateUTC         time.Time
	OpenPriceCents  int64
	ClosePriceCents int64
	HighPriceCents  int64
	LowPriceCents   int64
	VolumeShares    int64
}
