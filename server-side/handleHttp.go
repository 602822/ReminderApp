package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	CurrentDate string             `bson:"currentDate"`
	EventDate   string             `bson:"eventDate"`
}

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
}

var client *mongo.Client
var DBError error
var userId primitive.ObjectID
var newUser bool = true

func createNewUser() {
	id := primitive.NewObjectID()
	document := User{
		ID: id,
	}
	collection := client.Database("MyMongoDB").Collection("Users")
	_, err := collection.InsertOne(context.TODO(), document)

	if err != nil {
		fmt.Println("Error inserting User")
	}
	fmt.Println("User with Id :", id.String(), "inserted successfully")
	newUser = false
	userId = id
}

func retrieveUserEvents() {
	fmt.Println("UserId: ", userId)
	var events []bson.M
	context := context.TODO()
	eventCollection := client.Database("MyMongoDB").Collection("Events")
	filter, err := eventCollection.Find(context, bson.M{"_id": userId})

	if err != nil {
		fmt.Println("Error finding document in collection: ", err)
	}

	if err = filter.All(context, &events); err != nil {
		fmt.Println("Error decoding document into result: ", err)
	}
	fmt.Println("Events: ", events)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("You entered the main page")
	if newUser == true {
		createNewUser()
	}

	retrieveUserEvents()

	filepath := filepath.Join("client-side", "html", "index.html")
	http.ServeFile(w, r, filepath)
}

func newEventPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("You entered the create new event page")
	if r.Method == http.MethodPost {
		handleFormSubmit(w, r)
	}

	filepath := filepath.Join("client-side", "html", "createNewEvent.html")
	http.ServeFile(w, r, filepath)

}

func serveCSS(w http.ResponseWriter, r *http.Request) {
	filepath := filepath.Join("server-side", "output.css")
	http.ServeFile(w, r, filepath)
}

func serveJS(w http.ResponseWriter, r *http.Request) {
	filepath := filepath.Join("client-side", "js", "script.js")
	http.ServeFile(w, r, filepath)
}

func handleFormSubmit(w http.ResponseWriter, request *http.Request) {

	error := request.ParseForm()

	if error != nil {
		fmt.Println("Error parsing form data: ", error)
	}

	collection := client.Database("MyMongoDB").Collection("Events")

	title := request.FormValue("title")
	currentDate := request.FormValue("currentDate")
	eventDate := request.FormValue("eventDate")

	document := Event{
		ID:          userId,
		Title:       title,
		CurrentDate: currentDate,
		EventDate:   eventDate,
	}

	_, err := collection.InsertOne(context.TODO(), document)

	if err != nil {
		fmt.Println("Error adding document", err)
	}
	fmt.Println("Inserted a Event with Id: ", document.ID)

}

func connectToDB() {
	connectionString := "mongodb+srv://vefjeld:UPJ3jQhFqhhslAFO@event-data-cluster.a18gcij.mongodb.net/"
	clientOptions := options.Client().ApplyURI(connectionString)

	client, DBError = mongo.Connect(context.TODO(), clientOptions)
	if DBError != nil {
		fmt.Println("Error connecting to the DB", DBError)
		return
	}
}

func main() {
	connectToDB()
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/createNewEvent", newEventPage)
	http.HandleFunc("/server-side/output.css", serveCSS)
	http.HandleFunc("/js/script.js", serveJS)
	http.ListenAndServe(":8080", nil)
}
