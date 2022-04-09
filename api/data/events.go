package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Event struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Date        time.Time          `json:"date,omitempty" bson:"date,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Attendees   []Attendee         `json:"attendees" bson:"attendees"`
}

type Attendee struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FullName    string             `json:"fullName,omitempty" bson:"fullName,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	PhoneNumber string             `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	AttendStatus ConfirmationStatus `json:"attendStatus" bson:"attendStatus"`
}

type ConfirmationStatus int

const (
	Pending ConfirmationStatus = iota
	Yes
	No
)

var ctx context.Context
var client *mongo.Client

type EventsDB struct {
	eventCollection *mongo.Collection
}

func NewEventsDB() *EventsDB {
	// uri := os.Getenv("MONGODB_URI")
	uri := "mongodb://admin:admin@db:27017/?authSource=admin"

	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://docs.mongodb.com/drivers/go/current/usage-examples/#environment-variable")
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return &EventsDB{
		eventCollection: client.Database("test").Collection("events"),
	}
}

func DisconnectDB() {
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func (db *EventsDB) AddEvent(e Event) *mongo.InsertOneResult {
	result, err := db.eventCollection.InsertOne(ctx, &Event{
		Date:        time.Now(),
		Title:       e.Title,
		Description: e.Description,
		Attendees:   []Attendee{},
	})

	if err != nil {
		log.Println("Error:", err)
	}

	return result
}

func (db *EventsDB) GetEvents() []Event {
	var events []Event

	cursor, err := db.eventCollection.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var event Event
		cursor.Decode(&event)
		events = append(events, event)
	}
	if err := cursor.Err(); err != nil {
		panic(err)
	}

	return events
}

func (db *EventsDB) GetEventById(id primitive.ObjectID) Event {
	var event Event
	err := db.eventCollection.FindOne(ctx, Event{ID: id}).Decode(&event)
	if err != nil {
		panic(err)
	}
	return event
}

func (db *EventsDB) AddAttendees(eventId primitive.ObjectID, attendees []Attendee) *mongo.UpdateResult {
	event := db.GetEventById(eventId)
	event.Attendees = append(event.Attendees, initializeAttendees(attendees)...)

	result, err := db.eventCollection.UpdateByID(ctx, event.ID, bson.M{
		"$set": event,
	})
	if err != nil {
		log.Println("Error:", err)
		panic(err)
	}

	return result
}

func initializeAttendees(attendees []Attendee) []Attendee{
	var result []Attendee
	for _, element := range attendees{
		element.AttendStatus = Pending
		element.ID = primitive.NewObjectID()
		result = append(result,element)
	}
	log.Println(result)
	return result
}

func (db *EventsDB) UpdateAttendeeStatus(eventId primitive.ObjectID, attendeeId primitive.ObjectID, status ConfirmationStatus) *mongo.UpdateResult{
	log.Println(eventId, attendeeId, status)
	filter := bson.D{
		primitive.E{Key: "_id", Value: eventId},
		primitive.E{Key: "attendees._id", Value: attendeeId},
	}
	update := bson.M{
		"$set": bson.M{
			"attendees.$.attendStatus": status,
		},
	}
	result, _ := db.eventCollection.UpdateOne(ctx, filter, update)
	return result
}

func (db *EventsDB) RemoveAttendee(eventId primitive.ObjectID, attendeeId primitive.ObjectID) *mongo.UpdateResult{
	filter := bson.D{
		primitive.E{Key: "_id", Value: eventId},
	}
	update := bson.M{
		"$pull": bson.M{
			"attendees": bson.D{primitive.E{Key: "_id", Value: attendeeId}},
		},
	}
	result, _ := db.eventCollection.UpdateOne(ctx, filter, update)
	return result
}