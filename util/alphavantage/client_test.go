package alphavantage

import (
	"testing"
)

func TestMakeURL(t *testing.T) {
	expectedURL := "https://www.alphavantage.co/query?apikey=testkey&datatype=json&function=TIME_SERIES_DAILY&outputsize=compact&symbol=TEST"
	testClient := NewClient("testkey")
	url, err := testClient.makeURLFromInput(GetTimeSeriesInput{
		Function:   FunctionTypeDaily,
		OutputSize: OutputSizeCompact,
		DataType:   DataTypeJSON,
		Symbol:     "TEST",
	})
	if err != nil {
		t.Errorf("Expecting nil err, got %v", err)
	}
	if url.String() != expectedURL {
		t.Errorf("Expected URL %v, but got %v", expectedURL, url.String())
	}
}
