package controllers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/rsgregorio/CatecismoResponde/models"
)

func Handler(db *sql.DB, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		text := r.FormValue("text")
		data, err := models.FetchData(db, text)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data)
	}
}
