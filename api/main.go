package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ignj/go-invite/data"
	"github.com/ignj/go-invite/handlers"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome Home!")
}

func main() {
	// add env file
	// if err := gofotenv.Load(); err != nil {
	// 	log.Println("No .env file found")
	// }

	// create database instance
	// todo: inject conn string fetched from secrets .env
	db := data.NewEventsDB()

	// create handler
	eh := handlers.NewEvents(db)

	// create new server mux
	sm := mux.NewRouter().StrictSlash(true)

	// handlers for api
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/events", eh.GetEvents)
	getR.HandleFunc("/events/{id}", eh.GetEventById)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/events", eh.Create)

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/events/{id}/attendees", eh.AddAttendees)

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create a new server
	s := http.Server{
		Addr:    ":8080", // configure the bind address
		Handler: ch(sm),  // set the default handler
		// ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	defer data.DisconnectDB()
	s.Shutdown(ctx)
}
