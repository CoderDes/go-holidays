package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	resp, err := http.Get("https://date.nager.at/api/v2/publicholidays/2020/UA")

	if err != nil {
		fmt.Println("ERROR:", err)
	}

	defer resp.Body.Close()

	fmt.Println("RESPONSE STATUS:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)

	var jsonEncodedHolidays string
	for scanner.Scan() {
		jsonEncodedHolidays += scanner.Text()
	}

	decodedHolidays := decodeHolidays(jsonEncodedHolidays)

	currentDate := getCurrentDate()

	checkHolidays(decodedHolidays, currentDate)
}

type Holiday struct {
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

func getCurrentDate() time.Time {
	now := time.Now()
	currentDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.UTC().Location())
	return currentDate
}

func decodeHolidays(holidays string) []Holiday {
	bs := []byte(holidays)

	var holidaysDecoded []Holiday

	err := json.Unmarshal(bs, &holidaysDecoded)

	if err != nil {
		fmt.Println(err)
	}
	// ================ DELETE
	testStruct := Holiday{
		Date:        "2020-05-22",
		LocalName:   "TEST Friday",
		Name:        "TEST Friday",
		CountryCode: "UA",
		Fixed:       false,
		Global:      false,
		Counties:    nil,
		LaunchYear:  nil,
		Type:        "TEST TYPE",
	}
	testSlice := []Holiday{testStruct}
	for _, holiday := range holidaysDecoded {
		testSlice = append(testSlice, holiday)
	}
	holidaysDecoded = testSlice
	// ======================== DELETE
	return holidaysDecoded
}

func checkHolidays(holidays []Holiday, currentDate time.Time) {
	const layoutISO = "2006-01-02"
	var holidayLength float64

	for _, holiday := range holidays {
		holidayDate, err := time.Parse(layoutISO, holiday.Date)

		if err != nil {
			fmt.Println(err)
		}
		if currentDate.Equal(holidayDate) {
			fmt.Println("TODAY IS A HOLIDAY", holiday)
			holidayLength = calcHolidayLength(holidayDate)
			fmt.Printf("Holiday length is %v days", holidayLength)
			break
		}
		if currentDate.Before(holidayDate) {
			fmt.Println("TODAY IS:", currentDate)
			fmt.Println("NEXT HOLIDAY IS:", holiday)
			holidayLength = calcHolidayLength(holidayDate)
			fmt.Printf("Holiday length is %v days", holidayLength)
			break
		}
	}
}

func calcHolidayLength(holidayDate time.Time) float64 {
	const hoursInADay = 24
	var holidayStart time.Time
	var holidayEnd time.Time
	var holidayDays float64

	holidayWeekday := holidayDate.Weekday()

	switch holidayWeekday.String() {
	case "Monday":
		holidayStart = holidayDate.AddDate(0, 0, -2)
		holidayEnd = holidayDate.AddDate(0, 0, 1)
		holidayDays = holidayEnd.Sub(holidayStart).Hours() / hoursInADay
	case "Friday":
		holidayEnd = holidayDate.AddDate(0, 0, 3)
		holidayDays = holidayEnd.Sub(holidayDate).Hours() / hoursInADay
	case "Sunday":
		holidayStart = holidayDate.AddDate(0, 0, -1)
		holidayEnd = holidayDate.AddDate(0, 0, 2)
		holidayDays = holidayEnd.Sub(holidayStart).Hours() / hoursInADay
	}
	return holidayDays
}
