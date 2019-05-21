package symbols

type Symbol struct {
	ID                   string `db:"id"`
	Name                 string `db:"name"`
	Symbol               string `db:"symbol"`
	MarketCapitalization int64  `db:"market_capitalization"`
	Sector               string `db:"sector"`
	Industry             string `db:"industry"`
	Exchange             string `db:"exchange"`
}
