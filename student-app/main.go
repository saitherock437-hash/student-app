package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {

	connStr := os.Getenv("DATABASE_URL")

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", formHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/view", viewHandler)

	fmt.Println("Server running...")
	http.ListenAndServe(":10000", nil)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<h2>Student Form</h2>
	<form action="/submit" method="post">
	Name: <input type="text" name="name"><br><br>
	Email: <input type="email" name="email"><br><br>
	Course: <input type="text" name="course"><br><br>
	<input type="submit" value="Submit">
	</form>
	<br>
	<a href="/view">View Students</a>
	`
	fmt.Fprint(w, html)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	course := r.FormValue("course")

	_, err := db.Exec("INSERT INTO students(name,email,course) VALUES($1,$2,$3)", name, email, course)
	if err != nil {
		fmt.Fprint(w, "Error saving data")
		return
	}

	fmt.Fprint(w, "Saved ✅ <br><a href='/'>Go Back</a>")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, email, course FROM students")
	if err != nil {
		fmt.Fprint(w, "Error fetching data")
		return
	}
	defer rows.Close()

	html := "<h2>Student List</h2><table border='1'><tr><th>ID</th><th>Name</th><th>Email</th><th>Course</th></tr>"

	for rows.Next() {
		var id int
		var name, email, course string
		rows.Scan(&id, &name, &email, &course)

		html += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td></tr>", id, name, email, course)
	}

	html += "</table><br><a href='/'>Go Back</a>"

	fmt.Fprint(w, html)
}