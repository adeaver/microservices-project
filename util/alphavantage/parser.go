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

func parseResponseCSV(body string, function FunctionType) ([]*EquitySnapshot, error) {
	lines := strings.SplitN(body, "\n", -1)
	keys := strings.SplitN(lines[0], ",", -1)
	var out []*EquitySnapshot
	for _, data := range lines[1:] {
		parts := strings.SplitN(data, ",", -1)
		dataMap := make(map[string]string)
		for i, p := range parts {
			dataMap[keys[i]] = p
		}
		snapshot, err := csvDataMapToEquitySnapshot(dataMap, function)
		if err != nil {
			// TODO: logging
			fmt.Printf("error on line %s: %s", data, err.Error())
			continue
		}
		out = append(out, snapshot)
	}
	return out, nil
}

func csvDataMapToEquitySnapshot(dataMap map[string]string, function FunctionType) (*EquitySnapshot, error) {
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
	var volumeShares *int64
	volumeSharesStr, ok := dataMap[string(csvKeyVolume)]
	if ok {
		volumeSharesInt, err := strconv.ParseInt(volumeSharesStr, 10, 64)
		if err != nil {
			return nil, err
		}
		volumeShares = &volumeSharesInt
	}
	timestampStr, ok := dataMap[string(csvKeyTimestamp)]
	if !ok {
		return nil, fmt.Errorf("no timestamp found")
	}
	time, err := parseTimestamp(timestampStr, function)
	if err != nil {
		return nil, err
	}
	return &EquitySnapshot{
		Time:            *time,
		OpenPriceCents:  openPriceCents,
		ClosePriceCents: closePriceCents,
		HighPriceCents:  highPriceCents,
		LowPriceCents:   lowPriceCents,
		VolumeShares:    volumeShares,
	}, nil
}

func getCentsValueForKeyOrNil(key csvKey, dataMap map[string]string) (*int64, error) {
	value, ok := dataMap[key.String()]
	if !ok {
		fmt.Println(fmt.Sprintf("no value found for key %v", key.String()))
		fmt.Println(fmt.Sprintf("%+v", dataMap))
		return nil, nil
	}
	float, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, err
	}
	centsValue := int64(math.Round(float * 100))
	return &centsValue, nil
}

func parseTimestamp(timestampString string, function FunctionType) (*time.Time, error) {
	var layout string
	switch function {
	case FunctionTypeDaily:
		layout = "2006-01-02"
	default:
		return nil, fmt.Errorf("unsupported function type")
	}
	t, err := time.Parse(layout, timestampString)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
