package alphavantage

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type csvKey string

const (
	csvKeyOpen      csvKey = "open"
	csvKeyClose     csvKey = "close"
	csvKeyLow       csvKey = "low"
	csvKeyHigh      csvKey = "high"
	csvKeyVolume    csvKey = "volume"
	csvKeyTimestamp csvKey = "timestamp"
)

func (k *csvKey) String() string {
	return string(*k)
}

func parseResponseCSV(body string) ([]*EquitySnapshot, error) {
	lines := strings.SplitN(body, "\n", -1)
	keys := strings.SplitN(lines[0], ",", -1)
	var out []*EquitySnapshot
	for _, data := range lines[1:] {
		parts := strings.SplitN(data, ", ", -1)
		dataMap := make(map[string]string)
		for i, p := range parts {
			dataMap[keys[i]] = p
		}
		snapshot, err := csvDataMapToEquitySnapshot(dataMap)
		// This should maybe continue
		if err != nil {
			return nil, err
		}
		out = append(out, snapshot)
	}
	return out, nil
}

func csvDataMapToEquitySnapshot(dataMap map[string]string) (*EquitySnapshot, error) {
	openPriceCents, err := getCentsValueForKeyOrNil(csvKeyOpen, dataMap)
	if err != nil {
		return nil, err
	}
	closePriceCents, err := getCentsValueForKeyOrNil(csvKeyClose, dataMap)
	if err != nil {
		return nil, err
	}
	highPriceCents, err := getCentsValueForKeyOrNil(csvKeyHigh, dataMap)
	if err != nil {
		return nil, err
	}
	lowPriceCents, err := getCentsValueForKeyOrNil(csvKeyLow, dataMap)
	if err != nil {
		return nil, err
	}
	return &EquitySnapshot{
		Time:            time.Now(),
		OpenPriceCents:  openPriceCents,
		ClosePriceCents: closePriceCents,
		HighPriceCents:  highPriceCents,
		LowPriceCents:   lowPriceCents,
	}, nil
}

func getCentsValueForKeyOrNil(key csvKey, dataMap map[string]string) (*int64, error) {
	value, ok := dataMap[key.String()]
	if !ok {
		fmt.Println(fmt.Sprintf("no value found for key %v", key.String()))
		return nil, nil
	}
	float, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, err
	}
	centsValue := int64(math.Round(float) * 100)
	return &centsValue, nil
}
