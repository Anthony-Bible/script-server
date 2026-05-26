package main

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

//go:embed templates/404.html
var notFoundHTML string

// notFoundTmpl renders the custom "WANTED" 404 page.
var notFoundTmpl = template.Must(template.New("404").Parse(notFoundHTML))

type connector struct {
	db *sql.DB
}

type row struct {
	Id    int
	Value string
}

func createTable(db *sql.DB, tableName string) {
	//_, err := db.Exec("CREATE TABLE IF NOT EXISTS tattle (id INTEGER PRIMARY KEY, value TEXT)")
	_, err := db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY, value TEXT)", tableName))
	if err != nil {
		fmt.Println(err)
	}
}

func newConnector() (*connector, error) {
	db, err := sql.Open("sqlite3", "file:tattle.db:?cache=shared")

	if err != nil {
		return nil, err
	}
	return &connector{db: db}, nil
}
func displayInBrowserHandler(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Content-Type to text/plain to ensure it's displayed in the browser.
		w.Header().Set("Content-Type", "text/plain")
		//make sure the content is displayed inline in the browser
		w.Header().Set("Content-Disposition", "inline")
		h.ServeHTTP(w, r)
	})
}

// notFoundResponder wraps a ResponseWriter so that a 404 from the inner
// handler is replaced with the custom WANTED-poster 404 page.
type notFoundResponder struct {
	http.ResponseWriter
	r           *http.Request
	intercepted bool
}

// WriteHeader intercepts 404 status codes and renders the custom page
// instead of the inner handler's default body.
func (n *notFoundResponder) WriteHeader(status int) {
	if status == http.StatusNotFound {
		n.intercepted = true
		renderNotFound(n.ResponseWriter, n.r)
		return
	}
	n.ResponseWriter.WriteHeader(status)
}

// Write swallows the inner handler's body bytes once a 404 has been
// intercepted, so the default "404 page not found" never reaches the client.
func (n *notFoundResponder) Write(b []byte) (int, error) {
	if n.intercepted {
		return len(b), nil
	}
	return n.ResponseWriter.Write(b)
}

// withCustom404 wraps a handler and serves a custom 404 page when the
// inner handler would otherwise return http.StatusNotFound.
func withCustom404(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nr := &notFoundResponder{ResponseWriter: w, r: r}
		h.ServeHTTP(nr, r)
	})
}

// renderNotFound writes the custom WANTED 404 page for the given request.
func renderNotFound(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Del("Content-Disposition")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotFound)
	if err := notFoundTmpl.Execute(w, map[string]string{"Path": r.URL.Path}); err != nil {
		slog.Error("failed to render 404 template", "error", err)
	}
}

// Save - save the values to the database
func (c connector) save(values map[string][]string) {
	// save the values to the databased
	db := c.db
	// convert the values to a json string
	// insert the values into the database
	valuesJson := fmt.Sprintf("%v", values)
	_, err := db.Exec("INSERT INTO tattle (value) VALUES (?)", valuesJson)
	if err != nil {
		slog.Error(fmt.Sprintf("there was an error executing the query: %w", err))
	}

}

// SaveError - save the error to the database in the error table
func (c connector) saveError(values map[string][]string) {
	// save the values to the databased
	db := c.db
	// convert the values to a json string
	// insert the values into the database
	valuesJson := fmt.Sprintf("%v", values)
	_, err := db.Exec("INSERT INTO error (value) VALUES (?)", valuesJson)
	if err != nil {
		slog.Error(fmt.Sprintf("there was an error executing the query: %w", err))
	}
}

// Load - load the values from the database
func (c connector) load() []row {
	// load the values from the database
	db := c.db
	// select the values from the database
	// create the select statement
	rows, err := db.Query("SELECT * FROM tattle")
	if err != nil {
		slog.Error(fmt.Sprintf("there was an error executing the query: %w", err))
	}
	// iterate over the rows and print them out
	var allValues []row
	for rows.Next() {
		r := row{}
		err := rows.Scan(&r.Id, &r.Value)
		if err != nil {
			slog.Error(fmt.Sprintf("there was an error scanning the row: %w", err))
		}
		allValues = append(allValues, r)

	}
	return allValues
}
func (c connector) tattle(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Oops, you've been tattled on!"))
	// log out all the queery components
	slog.Info(fmt.Sprintf("Values passed: %v\n", r.URL.Query()))
	// get the values from the query and insert them into the database
	values := r.URL.Query()
	c.save(values)

	// put the values in a sqlite database

}

// Tell - tell the values from the database
func (c connector) tell(w http.ResponseWriter, r *http.Request) {
	// get the values from the database
	values := c.load()
	// print the values to the screen
	jsonBytes, err := json.Marshal(values)
	if err != nil {
		slog.Error(fmt.Sprintf("there was an error marshalling the values: %w", err))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)

}

// Report an error takes a request and logs and stores the error
func (c connector) reportError(w http.ResponseWriter, r *http.Request) {
	// log the error
	slog.Error(fmt.Sprintf("There was an error: %v", r.URL.Query()))
	// save the error
	values := r.URL.Query()
	c.saveError(values)
	// send a response
	_, _ = w.Write([]byte("Error reported"))
}
func main() {
	cn, err := newConnector()
	if err != nil {
		log.Fatal(err)
	}
	createTable(cn.db, "tattle")
	createTable(cn.db, "error")
	http.Handle("POST /tattle", http.HandlerFunc(cn.tattle))
	http.Handle("GET /tattle", http.HandlerFunc(cn.tell))
	http.Handle("POST /error", http.HandlerFunc(cn.reportError))
	fs := http.FileServer(http.Dir("./scripts"))
	http.Handle("/", withCustom404(displayInBrowserHandler(fs))) // Wrap the FileServer with our custom handler
	port := "8080"
	log.Printf("Serving on port %s...", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
