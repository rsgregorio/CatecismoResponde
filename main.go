package main

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/rsgregorio/CatecismoResponde/controllers"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres password=root dbname=catecismo sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tmpl := template.Must(template.ParseFiles("views/template.html"))

	http.HandleFunc("/", controllers.Handler(db, tmpl))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8080", nil)
}
