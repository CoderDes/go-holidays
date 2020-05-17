package main

import (
	"bufio"
	"fmt"
	"net/http"
)

func main() {
	resp, err := http.Get("https://date.nager.at/api/v2/publicholidays/2020/UA")

	if err != nil {
		fmt.Println("ERROR:", err)
	}

	defer resp.Body.Close()

	fmt.Println("RESPONSE STATUS:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)

	var result string
	for scanner.Scan() {
		result += scanner.Text()
	}
	fmt.Println(result)

}
