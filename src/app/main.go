package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const HTMLTemplate = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Главная страница</title>
	</head>
	<body>
		<h3>Notes</h3>
		<form method="POST" action="add">
			<label>New note:</label><br />
			<input type="text" name="text"><br />
			<input type="submit">
		</form>
		<ul>
			{{range .Notes}}
			<li>{{.Id}}: {{.Text}}
				<form method="POST" action="delete">
				<button type="submit" name="id" value="{{.Id}}">Delete</button>
				</form>
			</li>
			{{end}}
		</ul>
	</body>
</html>
`

func createDBConnection() (*sql.DB, error) {

	dbUser, exists := os.LookupEnv("DB_USER")
	if !exists {
		dbUser = "root"
	}

	dbPassword, exists := os.LookupEnv("DB_PASSWORD")
	if !exists {
		dbPassword = "password"
	}

	dbHost, exists := os.LookupEnv("DB_HOST")
	if !exists {
		dbHost = "localhost"
	}

	dbPort, exists := os.LookupEnv("DB_PORT")
	if !exists {
		dbPort = "3306"
	}

	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		dbName = "db"
	}

	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	log.Println("opening sql connection")

	db, err := createDBConnection()
	if err != nil {
		log.Fatalf("failed to open sql connection: %v", err)
	}

	log.Println("successfully connected to mysql")

	tmpl, err := template.New("test").Parse(HTMLTemplate)
	if err != nil {
		log.Fatalf("failed to execute template: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		notes, err := GetAll(db)
		if err != nil {
			log.Fatalf("failed to get notes: %v", err)
		}
		tmpl.Execute(w, struct{ Notes []Note }{notes})
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		err := Add(db, r)
		if err != nil {
			log.Fatalf("failed to add note: %v", err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		err := Delete(db, r)
		if err != nil {
			log.Fatalf("failed to delete note: %v", err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Fatal(http.ListenAndServe(":7777", nil))
}
