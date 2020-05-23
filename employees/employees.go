package employees

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/EugZ/go-holidays/ask"
	_ "github.com/go-sql-driver/mysql"
)

func ExamineEmployees() {
	dbUser, dbPass, dbName := ask.AskCredsToConnect()

	dbAddr := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", dbUser, dbPass, dbName)
	db, err := sql.Open("mysql", dbAddr)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	fmt.Println("Successfully connected to db")

	type queryType struct {
		query    string
		question string
	}

	queries := []queryType{
		{
			query:    "SELECT title, first_name, last_name, salary FROM dept_manager INNER JOIN titles ON dept_manager.emp_no = titles.emp_no INNER JOIN salaries ON dept_manager.emp_no = salaries.emp_no INNER JOIN employees ON dept_manager.emp_no = employees.emp_no WHERE title LIKE '%Manager%';",
			question: "Find all current managers with their title, first/last name, current salary. Continue? y or n: ",
		},
		{
			query:    "SELECT DISTINCT dept_name, title, first_name, last_name, hire_date, EXTRACT(YEAR FROM dept_emp.to_date) - EXTRACT(YEAR FROM dept_emp.from_date) AS how_many_years FROM employees INNER JOIN dept_emp ON employees.emp_no = dept_emp.emp_no INNER JOIN departments ON dept_emp.dept_no = departments.dept_no INNER JOIN titles ON employees.emp_no = titles.emp_no INNER JOIN salaries ON employees.emp_no = salaries.emp_no;",
			question: "Find all employees (title, first/last name, hire date, experience in years). Continue? y or n: ",
		},
		{
			// NOTE: there isn't 2020 year in provided db with employyes;
			query:    "SELECT dept_name, COUNT(dept_emp.emp_no) AS employee_count, SUM(salary) AS salary_sum FROM dept_emp INNER JOIN departments ON dept_emp.dept_no = departments.dept_no INNER JOIN salaries ON dept_emp.emp_no = salaries.emp_no WHERE MONTH(salaries.to_date) = MONTH(CURRENT_DATE()) AND YEAR(salaries.to_date) = YEAR(CURRENT_DATE()) GROUP BY dept_name;",
			question: "Find all departments, their employee count, sum salary. Continue? y or n: ",
		},
	}

	for _, q := range queries {
		var answerToCont string

		fmt.Print(q.question)
		fmt.Fscan(os.Stdin, &answerToCont)

		if answerToCont == "y" {
			requestToDB(db, q.query)
		} else {
			continue
		}
	}

}

func requestToDB(db *sql.DB, query string) {

	rows, err := db.Query(query)

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
