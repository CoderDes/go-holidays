package employees

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ExamineEmployees() {
	fmt.Println("=============")
	fmt.Println("Employees part")

	// TODO: move to separate function
	var dbUser string
	var dbPass string
	var dbName string
	fmt.Print("Enter database user: ")
	fmt.Fscan(os.Stdin, &dbUser)
	fmt.Print("Enter password: ")
	fmt.Fscan(os.Stdin, &dbPass)
	fmt.Print("Enter database name you want connect to: ")
	fmt.Fscan(os.Stdin, &dbName)

	// TODO: ask for user and DBname from connection
	dbAddr := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", dbUser, dbPass, dbName)
	db, err := sql.Open("mysql", dbAddr)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	fmt.Println("Successfully connected to db")

	rows, err := db.Query("SELECT * FROM employees")

	if err != nil {
		panic(err.Error())
	}

	columns, err := rows.Columns()

	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)

		if err != nil {
			panic(err.Error())
		}

		var value string

		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------")
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

}

func selectCurrentManagers() {

}

func findAllEmployees() {

}

func findAllDepartments() {

}
