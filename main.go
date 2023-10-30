package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

type Data struct {
	Numero string
	Texto  template.HTML
}

func highlight(text, query string) template.HTML {
	words := strings.Fields(query) // Divide o texto em palavras
	for _, word := range words {
		text = strings.ReplaceAll(text, word, "<mark>"+word+"</mark>")
	}
	return template.HTML(text)
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template.html"))

	db, err := sql.Open("postgres", "user=postgres password=root dbname=catecismo sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	text := r.FormValue("text")
	data := []Data{}

	words := strings.Fields(text)         // Divide o texto em palavras
	tsQuery := strings.Join(words, " & ") // Junta as palavras com o operador &

	rows, err := db.Query("SELECT numero, texto FROM paragrafos WHERE to_tsvector('portuguese', texto) @@ to_tsquery('portuguese', $1) ORDER BY numero", tsQuery)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var numero string
		var texto string
		err = rows.Scan(&numero, &texto)
		if err != nil {
			panic(err)
		}
		textoHTML := highlight(texto, text)
		data = append(data, Data{Numero: numero, Texto: textoHTML})
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
