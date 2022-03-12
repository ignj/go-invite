package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ignj/go-invite/data"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Events struct {
	eventsDB *data.EventsDB
}

type Invitation struct {
	InvitationHash string `json:"invitationHash"`
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

func (e *Events) GetInvitationLink(rw http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	eventId, _ := primitive.ObjectIDFromHex(params["id"])
	attendeeId, _ := primitive.ObjectIDFromHex(params["attendeeId"])

	// TODO: check if entities exist / are valid in db

	eventUser := eventId.String() + "_" + attendeeId.String()

	json.NewEncoder(rw).Encode(&Invitation{
		base64.StdEncoding.EncodeToString([]byte(eventUser)),
	})
}

func (e *Events) AcceptInvitation(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hashedInvitation, _ := base64.StdEncoding.DecodeString(params["hash"])

	ids := strings.Split(string(hashedInvitation), "_")

	eventId, _ := primitive.ObjectIDFromHex(ids[0][10:34])
	attendeeId, _ := primitive.ObjectIDFromHex(ids[1][10:34])

	json.NewEncoder(rw).Encode(e.eventsDB.UpdateAttendeeStatus(eventId, attendeeId, data.Yes)) 
}

func (e *Events) DeclineInvitation(rw http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	hashedInvitation, _ := base64.StdEncoding.DecodeString(params["hash"])

	ids := strings.Split(string(hashedInvitation), "_")

	eventId, _ := primitive.ObjectIDFromHex(ids[0][10:34])
	attendeeId, _ := primitive.ObjectIDFromHex(ids[1][10:34])

	json.NewEncoder(rw).Encode(e.eventsDB.UpdateAttendeeStatus(eventId, attendeeId, data.No))
}

func (e *Events) RemoveAttendee(rw http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	eventId, _ := primitive.ObjectIDFromHex(params["id"])
	attendeeId, _ := primitive.ObjectIDFromHex(params["attendeeId"])

	//todo: if attendee confirmed attendance, return 409 and prevent action

	json.NewEncoder(rw).Encode(e.eventsDB.RemoveAttendee(eventId, attendeeId))
}