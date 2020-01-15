package tools

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func lastTimeContainDateAndHour(st string) (string, error) {
	var dt string
	if st == "yyyy-mm-dd-hh" {
		dt = time.Now().Add(-time.Minute * 60).Format("2006-01-02-15")
		return dt, nil
	}
	if st == "yyyy-mm-dd_hh" {
		dt = time.Now().Add(-time.Minute * 60).Format("2006-01-02_15")
		return dt, nil
	}
	if st == "yyyymmdd-hh" {
		dt = time.Now().Add(-time.Minute * 60).Format("20060102-15")
		return dt, nil
	}
	if st == "yyyymmdd_hh" {
		dt = time.Now().Add(-time.Minute * 60).Format("20060102_15")
		return dt, nil
	}

	errMsg := fmt.Sprintf("source time format: %s invalid", st)
	return "", errors.New(errMsg)
}

func today(st string) (string, error) {
	var (
		dt     string
		errMsg = fmt.Sprintf("source time format: %s invalid", st)
	)

	switch st {
	case "yyyy-mm-dd":
		dt = time.Now().AddDate(0, 0, 0).Format("2006-01-02")
		return dt, nil
	case "yyyymmdd":
		dt = time.Now().AddDate(0, 0, 0).Format("20060102")
		return dt, nil
	case "yyyy_mm_dd":
		dt = time.Now().AddDate(0, 0, 0).Format("2006_01_02")
		return dt, nil
	case "yyyy-mm_dd":
		dt = time.Now().AddDate(0, 0, 0).Format("2006-01_02")
		return dt, nil
	case "yyyy_mm-dd":
		dt = time.Now().AddDate(0, 0, 0).Format("2006_01-02")
		return dt, nil
	default:
		return dt, errors.New(errMsg)
	}

}

func tomorrow(st string) (string, error) {

	var (
		dt     string
		errMsg = fmt.Sprintf("source time format: %s invalid", st)
	)

	switch st {
	case "yyyy-mm-dd":
		dt = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		return dt, nil
	case "yyyymmdd":
		dt = time.Now().AddDate(0, 0, 1).Format("20060102")
		return dt, nil
	case "yyyy_mm_dd":
		dt = time.Now().AddDate(0, 0, 1).Format("2006_01_02")
		return dt, nil
	case "yyyy-mm_dd":
		dt = time.Now().AddDate(0, 0, 1).Format("2006-01_02")
		return dt, nil
	case "yyyy_mm-dd":
		dt = time.Now().AddDate(0, 0, 1).Format("2006_01-02")
		return dt, nil
	default:
		return dt, errors.New(errMsg)
	}
}

func yesterday(st string) (string, error) {
	var (
		dt     string
		errMsg = fmt.Sprintf("source time format: %s invalid", st)
	)

	switch st {
	case "yyyy-mm-dd":
		dt = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		return dt, nil
	case "yyyymmdd":
		dt = time.Now().AddDate(0, 0, -1).Format("20060102")
		return dt, nil
	case "yyyy_mm_dd":
		dt = time.Now().AddDate(0, 0, -1).Format("2006_01_02")
		return dt, nil
	case "yyyy-mm_dd":
		dt = time.Now().AddDate(0, 0, -1).Format("2006-01_02")
		return dt, nil
	case "yyyy_mm-dd":
		dt = time.Now().AddDate(0, 0, -1).Format("2006_01-02")
		return dt, nil
	default:
		return dt, errors.New(errMsg)
	}
}

func lastHour(st string) string {
	h := time.Now().Hour() - 1
	if h == -1 {
		return fmt.Sprintf("%d", 23)
	}

	if h < 10 {
		return fmt.Sprintf("0%d", h)
	}

	return fmt.Sprintf("%d", h)
}

func nextHour(st string) string {
	h := time.Now().Hour() + 1
	if h < 10 {
		return fmt.Sprintf("0%d", h)
	}
	return fmt.Sprintf("%d", h)
}

func GetLastTime(st, split string) (string, error) {
	// 获取上个小时的时间，格式类型：2019091213
	if strings.Contains(st, "yyyy") && strings.Contains(st, "mm") && strings.Contains(st, "dd") && strings.Contains(st, "hh") {
		return lastTimeContainDateAndHour(st)
	}

	// 获取上个小时的日期，格式类型：20190912
	if strings.Contains(st, "yyyy") && strings.Contains(st, "mm") && strings.Contains(st, "dd") && !strings.Contains(st, "hh") {
		if split == "HOURLY" && time.Now().Hour() != 0 {
			return today(st)
		}
		return yesterday(st)
	}

	// 获取上个小时
	if !strings.Contains(st, "yyyy") && !strings.Contains(st, "mm") && !strings.Contains(st, "dd") && strings.Contains(st, "hh") {
		return lastHour(st), nil
	}

	return yesterday(st)

}
