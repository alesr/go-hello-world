package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// a very simple HTML template
const tmpl = `
<h1>Go &mdash; Hello World!</h1>

{{if .Error}}
    <p>Problem connecting to MySQL:</p>
    <code style="background:#eee;border:1px solid #ccc;padding:10px;">{{.Error.Error}}</code>
{{else}}
    <p>{{.Message}}</p>
{{end}}
`

// possible messages to be rendered
var (
	messageNoService = "There are no binded services, try binding one in the console."
	messageSuccess   = "Successfully connected to MySQL!"
)

// TemplateData contains values passed to the template
type TemplateData struct {
	Message string
	Error   error
}

// one handler to rule them all :)
func helloHandler(w http.ResponseWriter, r *http.Request) {
	message, err := tryToConnectToDb()

	t := template.Must(template.New("tmpl").Parse(tmpl))
	t.Execute(w, TemplateData{
		Message: message,
		Error:   err,
	})
}

// try to connect to the database
//
// If MYSQL_URI environment variable exists, a connection towards MySQL
// will be established.
func tryToConnectToDb() (string, error) {
	uri := os.Getenv("MYSQL_URI")
	if uri == "" {
		return messageNoService, nil
	}

	dsn, err := dbDsn(uri)
	if err != nil {
		return "", err
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", err
	}
	defer db.Close()

	return messageSuccess, nil
}

// dbDsn returns MySQL DSN to use for connection
func dbDsn(uri string) (string, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
			"%s@tcp(%s)%s",
			url.User.String(),
			url.Host,
			url.Path,
		),
		nil
}

// where the magic starts :)
func main() {
	http.HandleFunc("/", helloHandler)
	panic(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
