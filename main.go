package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"regexp"

	_ "github.com/lib/pq"
)

type Data struct {
	Numero string
	Texto  template.HTML
}

func highlight(text, query string) template.HTML {
	re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(query))
	highlightedText := re.ReplaceAllString(text, `<mark>$0</mark>`)
	return template.HTML(highlightedText)
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template.html"))

	db, err := sql.Open("postgres", "user=postgres password=root dbname=catecismo sslmode=disable")
	if err != nil {
		panic(err)
	}

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	text := r.FormValue("text")
	data := []Data{}

	rows, err := db.Query("SELECT numero, texto FROM paragrafos WHERE texto ILIKE $1 ORDER BY numero", "%"+text+"%")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var numero string
		var texto template.HTML
		err = rows.Scan(&numero, &texto)
		if err != nil {
			panic(err)
		}
		texto = highlight(string(texto), text)
		data = append(data, Data{Numero: numero, Texto: texto})
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
