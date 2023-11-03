package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"os"

	"github.com/joho/godotenv" // Importe a biblioteca godotenv
	_ "github.com/lib/pq"
	"github.com/rsgregorio/CatecismoResponde/controllers"
)

func init() {
	// Carregue as variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		panic("Erro ao carregar o arquivo .env: " + err.Error())
	}
}

func main() {

	connStr := "user=" + os.Getenv("DB_USER") + " " +
		"password=" + os.Getenv("DB_PASSWORD") + " " +
		"dbname=" + os.Getenv("DB_NAME") + " " +
		"host=" + os.Getenv("DB_HOST") + " " +
		"port=" + os.Getenv("DB_PORT") + " " +
		"sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tmpl := template.Must(template.ParseFiles("views/template.html", "views/menu.html", "views/content.html", "views/footer.html", "views/sobre.html"))

	http.HandleFunc("/", controllers.Handler(db, tmpl))
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/sobre", controllers.Sobre(tmpl))

	http.ListenAndServe(":8080", nil)
}
