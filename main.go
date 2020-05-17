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

	currentTime := getCurrentDate()
	fmt.Println(currentTime)

	fmt.Printf("======= %T", decodedHolidays)
}

// TODO: define current date
func getCurrentDate() time.Time {
	dt := time.Now()
	return dt
}

// TODO: loop over the each holiday and define date
func decodeHolidays(holidays string) interface{} {
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
	bs := []byte(holidays)

	var holidaysDecoded []Holiday

	err := json.Unmarshal(bs, &holidaysDecoded)

	if err != nil {
		fmt.Println(err)
	}

	return holidaysDecoded
}

// TODO: compare two dates
// TODO: check if holiday date is on Fri or Mon
// TODO: calc the whole length of holidays
