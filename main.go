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
	// testStruct := Holiday{
	// 	Date:        "2020-05-21",
	// 	LocalName:   "TEST THURSDAY",
	// 	Name:        "TEST THURSDAY",
	// 	CountryCode: "UA",
	// 	Fixed:       false,
	// 	Global:      false,
	// 	Counties:    nil,
	// 	LaunchYear:  nil,
	// 	Type:        "TEST TYPE",
	// }
	// testStruct := Holiday{
	// 	Date:        "2020-05-19",
	// 	LocalName:   "TEST ORDINARY",
	// 	Name:        "TEST ORDINARY",
	// 	CountryCode: "UA",
	// 	Fixed:       false,
	// 	Global:      false,
	// 	Counties:    nil,
	// 	LaunchYear:  nil,
	// 	Type:        "TEST TYPE",
	// }
	// testStruct := Holiday{
	// 	Date:        "2020-05-23",
	// 	LocalName:   "TEST SATURDAY",
	// 	Name:        "TEST SATURDAY",
	// 	CountryCode: "UA",
	// 	Fixed:       false,
	// 	Global:      false,
	// 	Counties:    nil,
	// 	LaunchYear:  nil,
	// 	Type:        "TEST TYPE",
	// }
	// testSlice := []Holiday{testStruct}
	// for _, holiday := range holidaysDecoded {
	// 	testSlice = append(testSlice, holiday)
	// }
	// holidaysDecoded = testSlice
	// ======================== DELETE
	return holidaysDecoded
}

func checkHolidays(holidays []Holiday, currentDate time.Time) {
	const layoutISO = "2006-01-02"
	var isToday bool
	var closestHoliday Holiday
	var holidayLength float64
	var dateRange string

	for _, holiday := range holidays {
		holidayDate, err := time.Parse(layoutISO, holiday.Date)

		if err != nil {
			fmt.Println(err)
		}
		if currentDate.Equal(holidayDate) {
			isToday = true
			// TODO: try to DRY
			closestHoliday = holiday
			holidayLength, dateRange = calcHolidayLength(holidayDate)
			break
		}
		if currentDate.Before(holidayDate) {
			// TODO: try to DRY
			closestHoliday = holiday
			holidayLength, dateRange = calcHolidayLength(holidayDate)
			break
		}
	}
	conclusion(closestHoliday, holidayLength, dateRange, isToday)
}

func calcHolidayLength(holidayDate time.Time) (float64, string) {
	const hoursInADay = 24
	var holidayStart time.Time
	var holidayEnd time.Time
	var holidayDays float64
	var dateRange string

	holidayWeekday := holidayDate.Weekday()
	// TODO: DRY IT !!!
	switch holidayWeekday.String() {
	case "Monday":
		holidayStart = holidayDate.AddDate(0, 0, -2)
		holidayEnd = holidayDate.AddDate(0, 0, 1)
		// TODO: to separate func
		dateRange = fmt.Sprintf("From %v %v to %v %v", holidayStart.Day(), holidayStart.Month(), holidayEnd.Day(), holidayEnd.Month())
		holidayDays = holidayEnd.Sub(holidayStart).Hours() / hoursInADay
	case "Friday":
		holidayEnd = holidayDate.AddDate(0, 0, 3)
		// TODO: to separate func
		dateRange = fmt.Sprintf("From %v %v to %v %v", holidayStart.Day(), holidayStart.Month(), holidayEnd.Day(), holidayEnd.Month())
		holidayDays = holidayEnd.Sub(holidayDate).Hours() / hoursInADay
	case "Saturday":
		holidayStart = holidayDate
		holidayEnd = holidayDate.AddDate(0, 0, 3)
		// TODO: to separate func
		dateRange = fmt.Sprintf("From %v %v to %v %v", holidayStart.Day(), holidayStart.Month(), holidayEnd.Day(), holidayEnd.Month())
		holidayDays = holidayEnd.Sub(holidayStart).Hours() / hoursInADay
	case "Sunday":
		holidayStart = holidayDate.AddDate(0, 0, -1)
		holidayEnd = holidayDate.AddDate(0, 0, 2)
		// TODO: to separate func
		dateRange = fmt.Sprintf("From %v %v to %v %v", holidayStart.Day(), holidayStart.Month(), holidayEnd.Day(), holidayEnd.Month())
		holidayDays = holidayEnd.Sub(holidayStart).Hours() / hoursInADay
	default:
		holidayEnd = holidayDate.AddDate(0, 0, 1)
		dateRange = fmt.Sprintf("on %v of %v", holidayDate.Day(), holidayDate.Month())
		holidayDays = holidayEnd.Sub(holidayDate).Hours() / hoursInADay
	}
	return holidayDays, dateRange
}

func conclusion(holiday Holiday, lengthInDays float64, dateRange string, isToday bool) {
	todayOrNot := "The closest holiday"
	dayOrDays := "days"

	if isToday {
		todayOrNot = "Today"
	}
	if lengthInDays == 1 {
		dayOrDays = "day"
	}

	output := fmt.Sprintf("%v is a %v on %v. It will last %v %v: %v.", todayOrNot, holiday.Name, holiday.Date, lengthInDays, dayOrDays, dateRange)

	fmt.Println(output)
}
