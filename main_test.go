package main

import (
	"fmt"
	"testing"
	"time"
)

func TestDefineHolidayLength(t *testing.T) {
	const layoutISO = "2006-01-02"

	type testHoliday struct {
		Date        string
		LocalName   string
		Name        string
		CountryCode string
		Fixed       bool
		Global      bool
		Counties    interface{}
		LaunchYear  interface{}
		Type        string
	}

	testsHolidays := []testHoliday{
		testHoliday{
			Date:        "2020-06-17",
			LocalName:   "Test holiday that between Monday and Friday",
			Name:        "Tuesday, Wednesday or Thursday holiday",
			CountryCode: "UA",
			Fixed:       false,
			Global:      false,
			Counties:    nil,
			LaunchYear:  nil,
			Type:        "weekday",
		},
		testHoliday{
			Date:        "2020-06-22",
			LocalName:   "Holiday on Friday",
			Name:        "Holiday on Weekend",
			CountryCode: "UA",
			Fixed:       false,
			Global:      false,
			Counties:    nil,
			LaunchYear:  nil,
			Type:        "fri",
		},
		testHoliday{
			Date:        "2020-06-20",
			LocalName:   "Holiday on Saturday",
			Name:        "Holiday on Weekend",
			CountryCode: "UA",
			Fixed:       false,
			Global:      false,
			Counties:    nil,
			LaunchYear:  nil,
			Type:        "weekend",
		},
		testHoliday{
			Date:        "2020-06-21",
			LocalName:   "Holiday on Sunday",
			Name:        "Holiday on Weekend",
			CountryCode: "UA",
			Fixed:       false,
			Global:      false,
			Counties:    nil,
			LaunchYear:  nil,
			Type:        "weekend",
		},
		testHoliday{
			Date:        "2020-06-22",
			LocalName:   "Holiday on Monday",
			Name:        "Holiday on Weekend",
			CountryCode: "UA",
			Fixed:       false,
			Global:      false,
			Counties:    nil,
			LaunchYear:  nil,
			Type:        "mon",
		},
	}

	for _, value := range testsHolidays {
		holidayDate, err := time.Parse(layoutISO, value.Date)

		if err != nil {
			fmt.Println(err)
		}

		holidayDays, _ := defineHolidayLength(holidayDate)

		if (value.Type == "mon" || value.Type == "weekend" || value.Type == "fri") && holidayDays != 3 {
			t.Error("Expected:", 3, "Got:", holidayDays)
		}
		if value.Type == "weekday" && holidayDays != 1 {
			t.Error("Expected", 1, "Got", holidayDays)
		}
	}
}
