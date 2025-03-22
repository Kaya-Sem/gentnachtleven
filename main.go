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

		// Define the HTML template
		tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Gent Nachtleven Events</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f4f4f4;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            padding: 20px;
        }
        h1 {
            color: #333;
            text-align: center;
        }
        .event-list {
            max-height: 70vh;
            overflow-y: auto;
            padding: 10px;
            border: 1px solid #eee;
            border-radius: 4px;
        }
        .event-item {
            padding: 15px;
            margin-bottom: 15px;
            border-left: 4px solid #5c6bc0;
            background-color: #f9f9f9;
            border-radius: 4px;
        }
        .event-title {
            font-weight: bold;
            font-size: 18px;
            margin-bottom: 5px;
            color: #333;
        }
        .event-date {
            color: #5c6bc0;
            font-weight: bold;
            margin-bottom: 5px;
        }
        .event-description {
            color: #666;
            line-height: 1.5;
        }
        .provider-tag {
            display: inline-block;
            font-size: 12px;
            background-color: #e0e0e0;
            padding: 3px 8px;
            border-radius: 12px;
            margin-top: 8px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Gent Nachtleven Events</h1>
        <div class="event-list">
            {{if .Events}}
                {{range .Events}}
                    <div class="event-item">
                        <div class="event-date">{{formatDate .Date}}</div>
                        <div class="event-title">{{.Title}}</div>
                        <div class="event-description">{{.Description}}</div>
                        <div class="event-description">{{.Location}}</div>
                    </div>
                {{end}}
            {{else}}
                <p>No events found.</p>
            {{end}}
        </div>
    </div>
</body>
</html>
`

		// Create a new template and register a function to format dates
		t := template.New("events")
		t.Funcs(template.FuncMap{
			"formatDate": FormatDate,
		})

		// Parse the template
		t, err := t.Parse(tmpl)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
			return
		}

		// Execute the template with the collected events
		err = t.Execute(w, PageData{Events: allEvents})
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
