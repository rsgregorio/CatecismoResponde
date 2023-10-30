package main

import (
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/lib/pq"
)

type Result struct {
	Numero int
	Texto  string
}

var tmpl = template.Must(template.ParseFiles("template.html"))

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "user=postgres password=root dbname=catecismo sslmode=disable")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if r.Method == http.MethodPost {
		text := r.FormValue("text")
		rows, err := db.Query("SELECT numero, texto FROM paragrafos WHERE texto ILIKE $1 ORDER BY numero", "%"+text+"%")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []Result
		for rows.Next() {
			var result Result
			if err := rows.Scan(&result.Numero, &result.Texto); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			results = append(results, result)
		}
		tmpl.Execute(w, results)
	} else {
		tmpl.Execute(w, nil)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
