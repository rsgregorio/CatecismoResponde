package models

import (
	"database/sql"
	"html/template"
	"regexp"
	"strconv"
	"strings"
)

type Data struct {
	Numero string
	Texto  template.HTML
}

func FetchData(db *sql.DB, text string) ([]Data, error) {
	data := []Data{}

	// Dividindo o texto em palavras.
	words := strings.Fields(text)

	// Construindo a parte da consulta ILIKE.
	ilikeQuery := ""
	params := []interface{}{}
	for i, word := range words {
		ilikeQuery += "unaccent(texto) ILIKE unaccent($" + strconv.Itoa(i+1) + ") OR "
		params = append(params, "%"+word+"%")
	}
	ilikeQuery = strings.TrimSuffix(ilikeQuery, " OR ") // Removendo o último " OR "

	// Criando a string da consulta.
	queryString := "SELECT numero, texto FROM paragrafos WHERE " + ilikeQuery + " ORDER BY numero"

	// Executando a consulta.
	rows, err := db.Query(queryString, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterando sobre os resultados e populando a slice data.
	for rows.Next() {
		var numero string
		var texto string
		err = rows.Scan(&numero, &texto)
		if err != nil {
			return nil, err
		}
		textoHTML := highlight(texto, text)
		data = append(data, Data{Numero: numero, Texto: textoHTML})
	}

	return data, nil
}

func highlight(text, query string) template.HTML {
	// Dividindo a consulta em palavras.
	words := strings.Fields(query)

	// Criando uma expressão regular para cada palavra na consulta.
	for _, word := range words {
		re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(word))
		text = re.ReplaceAllString(text, "<mark>$0</mark>")
	}

	// Convertendo o texto resultante em HTML seguro.
	return template.HTML(text)
}
