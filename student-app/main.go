package main

import (
	"fmt"
	"net/http"
)

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

func main() {
	http.HandleFunc("/", homeHandler)

	fmt.Println("Server running...")
	http.ListenAndServe(":10000", nil)
}
