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

	return holidaysDecoded
}

func checkHolidays(holidays []Holiday, currentDate time.Time) {
	const layoutISO = "2006-01-02"
	for _, holiday := range holidays {
		holidayDate, err := time.Parse(layoutISO, holiday.Date)
		if err != nil {
			fmt.Println(err)
		}
		if currentDate.Equal(holidayDate) {
			fmt.Println("TODAY IS A HOLIDAY", holiday)
			break
		}
		if currentDate.Before(holidayDate) {
			fmt.Println("TODAY IS:", currentDate)
			fmt.Println("NEXT HOLIDAY IS:", holiday)
			break
		}

	}
}

// TODO: check if holiday date is on Fri or Mon
// TODO: calc the whole length of holidays
