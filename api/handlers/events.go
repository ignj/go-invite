package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ignj/go-invite/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Events struct {
	eventsDB *data.EventsDB
}

func NewEvents(edb *data.EventsDB) *Events {
	return &Events{
		eventsDB: edb,
	}
}

func (e *Events) Create(rw http.ResponseWriter, r *http.Request) {
	var event data.Event
	json.NewDecoder(r.Body).Decode(&event)

	result := e.eventsDB.AddEvent(event)

	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(result)
}

func (e *Events) GetEvents(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("content-type", "application/json")
	json.NewEncoder(rw).Encode(e.eventsDB.GetEvents())
}

func (e *Events) GetEventById(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	rw.Header().Set("content-type", "application/json")
	json.NewEncoder(rw).Encode(e.eventsDB.GetEventById(id))
}

func (e *Events) AddAttendees(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	eventId, _ := primitive.ObjectIDFromHex(params["id"])

	var attendees []data.Attendee
	json.NewDecoder(r.Body).Decode(&attendees)
	json.NewEncoder(rw).Encode(e.eventsDB.AddAttendees(eventId, attendees))
}
