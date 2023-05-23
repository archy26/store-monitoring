package utils

import (
	"fmt"
	"strconv"
	"time"
)

func StringToInteger(value string) int64 {
	res, _ := strconv.ParseInt(value, 10, 64)

	return res
}

func StringToTime(timeStr string) time.Time {
	layout := "2006-01-02 15:04:05.000000 MST" // Format layout based on the provided string
	fmt.Println(timeStr)
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		panic(err)
	}

	return parsedTime
}
