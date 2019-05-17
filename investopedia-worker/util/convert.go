package util

import (
	"strconv"
	"strings"
)

func centsFromFloatString(v string) (*int64, error) {
	centsStr := strings.Replace(v, ".", "", -1)
	centsStrNoCommas := strings.Replace(centsStr, ",", "", -1)
	i64, err := strconv.ParseInt(centsStrNoCommas, 10, 64)
	if err != nil {
		return nil, err
	}
	return &i64, nil
}

func MustParseCentsFromFloatString(v string) int64 {
	cents, err := centsFromFloatString(v)
	if err != nil {
		panic(err)
	}
	return *cents
}
