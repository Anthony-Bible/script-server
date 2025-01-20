package main

import (
	"log"
	"net/http"
)

func displayInBrowserHandler(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the Content-Type to text/plain to ensure it's displayed in the browser.
		w.Header().Set("Content-Type", "text/plain")
		//make sure the content is displayed inline in the browser
		w.Header().Set("Content-Disposition", "inline")
		h.ServeHTTP(w, r)
	})
}
func tattle(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello, friend!"))
	// log out all the queery components
	for k, v := range r.URL.Query() {
		log.Printf("Key: %s, Value: %s", k, v)

	}

}
func main() {
	fs := http.FileServer(http.Dir("./scripts"))
	http.Handle("/", displayInBrowserHandler(fs)) // Wrap the FileServer with our custom handler
	http.Handle("/tattle", http.HandlerFunc(tattle))
	port := "8080"
	log.Printf("Serving on port %s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
