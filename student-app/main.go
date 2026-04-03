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
	// Get DB URL from Render
	connStr := os.Getenv("DATABASE_URL")

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create table if not exists
	createTable()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/view", viewHandler)

	fmt.Println("Server running...")
	http.ListenAndServe(":10000", nil)
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS students (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		course TEXT
	)`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintf(w, `
		<h2>Student Form</h2>
		<form action="/submit" method="post">
			Name: <input type="text" name="name"><br><br>
			Email: <input type="email" name="email"><br><br>
			Course: <input type="text" name="course"><br><br>
			<input type="submit" value="Submit">
		</form>
		<br>
		<a href="/view">View Students</a>
	`)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		course := r.FormValue("course")

		_, err := db.Exec("INSERT INTO students(name, email, course) VALUES($1, $2, $3)", name, email, course)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Fprintf(w, "Data saved successfully! <br><a href='/'>Go Back</a>")
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, email, course FROM students")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintf(w, "<h2>Students List</h2><table border='1'><tr><th>ID</th><th>Name</th><th>Email</th><th>Course</th></tr>")

	for rows.Next() {
		var id int
		var name, email, course string
		rows.Scan(&id, &name, &email, &course)

		fmt.Fprintf(w, "<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td></tr>", id, name, email, course)
	}

	fmt.Fprintf(w, "</table><br><a href='/'>Back</a>")
}
