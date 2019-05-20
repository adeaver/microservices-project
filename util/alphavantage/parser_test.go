package alphavantage

import (
	"fmt"
	"strings"
	"testing"
)

func makePtr(v int64) *int64 {
	return &v
}

func testEquality(snapshot, expectedSnapshot EquitySnapshot) string {
	var errs []string
	// TODO: add date
	if *snapshot.OpenPriceCents != *expectedSnapshot.OpenPriceCents {
		errs = append(errs, fmt.Sprintf("open price is not equal. expected %v, got %v", *expectedSnapshot.OpenPriceCents, *snapshot.OpenPriceCents))
	}
	if *snapshot.HighPriceCents != *expectedSnapshot.HighPriceCents {
		errs = append(errs, fmt.Sprintf("high price is not equal. expected %v, got %v", *expectedSnapshot.HighPriceCents, *snapshot.HighPriceCents))
	}
	if *snapshot.LowPriceCents != *expectedSnapshot.LowPriceCents {
		errs = append(errs, fmt.Sprintf("low price is not equal. expected %v, got %v", *expectedSnapshot.LowPriceCents, *snapshot.LowPriceCents))
	}
	if *snapshot.ClosePriceCents != *expectedSnapshot.ClosePriceCents {
		errs = append(errs, fmt.Sprintf("close price is not equal. expected %v, got %v", *expectedSnapshot.ClosePriceCents, *snapshot.ClosePriceCents))
	}
	if *snapshot.VolumeShares != *expectedSnapshot.VolumeShares {
		errs = append(errs, fmt.Sprintf("volume shares is not equal. expected %v, got %v", *expectedSnapshot.VolumeShares, *snapshot.VolumeShares))
	}
	return strings.Join(errs, ", ")
}

func TestParseCSV(t *testing.T) {
	expectedSnapshot := EquitySnapshot{
		OpenPriceCents:  makePtr(12831),
		HighPriceCents:  makePtr(13046),
		LowPriceCents:   makePtr(12792),
		ClosePriceCents: makePtr(12807),
		VolumeShares:    makePtr(25051168),
	}
	response := `timestamp,open,high,low,close,volume
2019-05-17,128.3050,130.4600,127.9200,128.0700,25051168
2019-05-16,126.7500,129.3800,126.4600,128.9300,30112216
2019-05-15,124.2600,126.7100,123.7000,126.0200,24722708
2019-05-14,123.8700,125.8800,123.7000,124.7300,25266315`
	output, err := parseResponseCSV(response)
	if err != nil {
		t.Errorf("Expecting nil err, got %v", err)
	}
	if len(output) != 4 {
		t.Errorf("Expected 4 data points, got %v", len(output))
	}
	if errs := testEquality(*output[0], expectedSnapshot); len(errs) != 0 {
		t.Errorf(errs)
	}
}

func TestParseCSVChangeOrder(t *testing.T) {
	expectedSnapshot := EquitySnapshot{
		OpenPriceCents:  makePtr(12807),
		HighPriceCents:  makePtr(12831),
		LowPriceCents:   makePtr(13046),
		ClosePriceCents: makePtr(12792),
		VolumeShares:    makePtr(25051168),
	}
	response := `timestamp,high,low,close,open,volume
2019-05-17,128.3050,130.4600,127.9200,128.0700,25051168
2019-05-16,126.7500,129.3800,126.4600,128.9300,30112216
2019-05-15,124.2600,126.7100,123.7000,126.0200,24722708
2019-05-14,123.8700,125.8800,123.7000,124.7300,25266315`
	output, err := parseResponseCSV(response)
	if err != nil {
		t.Errorf("Expecting nil err, got %v", err)
	}
	if len(output) != 4 {
		t.Errorf("Expected 4 data points, got %v", len(output))
	}
	if errs := testEquality(*output[0], expectedSnapshot); len(errs) != 0 {
		t.Errorf(errs)
	}
}
