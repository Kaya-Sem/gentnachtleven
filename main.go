package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Kaya-Sem/gentnachtleven/cmd"
)

// PageData holds the data structure to be passed to the HTML template
type PageData struct {
	Events []cmd.Event
}

// FormatDate formats a Date struct as a string
func FormatDate(date cmd.Date) string {
	return fmt.Sprintf("%02d/%02d/%d", date.Day, date.Month, date.Year)
}

func main() {
	// Define HTTP handler for the root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Create a list of EventProvider interfaces
		providers := []cmd.EventProvider{
			cmd.NewPalmariumProvider(),
			cmd.NewClubWinterCircusProvider(),
		}

		// Collect all events
		var allEvents []cmd.Event
		for _, provider := range providers {
			events, err := provider.ScrapeEvents()
			if err != nil {
				log.Printf("Error scraping events from provider: %v", err)
				continue
			}
			allEvents = append(allEvents, events...)
		}

		// Create a template with functions
		funcMap := template.FuncMap{
			"formatDate": FormatDate,
		}

		// Parse the template file
		tmpl, err := template.New("events.html").Funcs(funcMap).ParseFiles("templates/events.html")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
			return
		}

		// Execute the template with the collected events
		err = tmpl.Execute(w, PageData{Events: allEvents})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
			return
		}
	})

	// Create a server with proper timeouts
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server
	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}

