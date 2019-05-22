package symbols

type Symbol struct {
	ID                   *string `db:"id,omitempty" json:"id,omitempty"`
	Name                 string  `db:"name" json:"name"`
	Symbol               string  `db:"symbol" json:"symbol"`
	MarketCapitalization int64   `db:"market_capitalization" json:"market_capitalization"`
	Sector               string  `db:"sector" json:"sector"`
	Industry             string  `db:"industry" json:"industry"`
	Exchange             string  `db:"exchange" json:"exchange"`
}
