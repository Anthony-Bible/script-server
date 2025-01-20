package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"log"
	"log/slog"
	"net/http"
)

type connector struct {
	db *sql.DB
}

type row struct {
	Id    int
	Value string
}

func createTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS tattle (id INTEGER PRIMARY KEY, value TEXT)")
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
func main() {
	cn, err := newConnector()
	if err != nil {
		log.Fatal(err)
	}
	createTable(cn.db)
	http.Handle("POST /tattle", http.HandlerFunc(cn.tattle))
	http.Handle("GET /tattle", http.HandlerFunc(cn.tell))
	fs := http.FileServer(http.Dir("./scripts"))
	http.Handle("/", displayInBrowserHandler(fs)) // Wrap the FileServer with our custom handler
	port := "8080"
	log.Printf("Serving on port %s...", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
