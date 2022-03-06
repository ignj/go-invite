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
	Date        time.Time `json:"date,omitempty" bson:"date,omitempty"`
	Title       string `json:"title,omitempty" bson:"title,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Invited     []Person `json:"invited,omitempty" bson:"invited,omitempty"`
	Accepted    []Person `json:"accepted,omitempty" bson:"accepted,omitempty"`
	Rejected    []Person `json:"rejected,omitempty" bson:"rejected,omitempty"`
}

type Person struct {
	ID          primitive.ObjectID
	FullName    string
	Email       string
	PhoneNumber string
}

var ctx context.Context
var client *mongo.Client

type EventsDB struct {
	eventCollection *mongo.Collection
}

func NewEventsDB() *EventsDB{
	// uri := os.Getenv("MONGODB_URI")
	uri := "mongodb://admin:admin@localhost:27017/?authSource=admin"
	
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

func DisconnectDB(){
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func (db *EventsDB) AddEvent(e Event) *mongo.InsertOneResult {
	result, err := db.eventCollection.InsertOne(ctx, &Event{
		Date: time.Now(),
		Title: e.Title,
		Description: e.Description,
	})

	if (err != nil){
		log.Println("Error:", err)
	}

	return result
}

func (db *EventsDB) GetEvents() []Event {
	var events []Event
	
	cursor, err := db.eventCollection.Find(ctx, bson.M{})
	if (err != nil){
		panic(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx){
		var event Event
		cursor.Decode(&event)
		events = append(events, event)
	}
	if err := cursor.Err(); err != nil {
		panic(err)
	}

	return events
}

func (db *EventsDB) GetEventById(id primitive.ObjectID) Event{
	var event Event
	err := db.eventCollection.FindOne(ctx, Event{ID: id}).Decode(&event)
	if err != nil {
		panic(err)
	}
	return event
}