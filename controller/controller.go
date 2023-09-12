package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	model "github.com/Veeresh-R-G/mongoapi/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "netflix"
const colName = "movies"

var collection *mongo.Collection

// connect with mongoDB
func init() {

	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatal(env_err)
	}

	CONN := os.Getenv("MONGO_URL")
	clientOptions := options.Client().ApplyURI(CONN)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Mongo Connection Successfull")

	collection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection : ", collection)
}

//MongoDB helpers

// insert 1 document
func insertOneMovie(movie model.Netflix) {
	resp, err := collection.InsertOne(context.TODO(), movie)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)

}

func updateOneMovie(movieID string) {
	//string to primitive.ObjectID that mongoDB understands
	id, _ := primitive.ObjectIDFromHex(movieID)

	//filter
	filter := bson.M{"_id": id}

	update := bson.M{"$set": bson.M{"watched": true}}

	res, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Updated : ", res)
}

func deleteOneMovie(movieID string) {

	id, _ := primitive.ObjectIDFromHex(movieID)

	filter := bson.M{"_id": id}

	//return value is the delete count
	res, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Movie Deleted : ", res)

}

func deleteAllMovie() {

	res, err := collection.DeleteMany(context.TODO(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deleted All Movies: ", res.DeletedCount)

}

func findAllMovies() []primitive.M {

	curr, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	var movies []primitive.M

	defer curr.Close(context.Background())

	for curr.Next(context.TODO()) {
		var movie bson.M

		if err = curr.Decode(&movie); err != nil {
			log.Fatal(err)
		}

		movies = append(movies, movie)
	}

	return movies
}

//GO Actual controller

func GetAllMovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	allMovies := findAllMovies()

	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Origin", "*")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix

	_ = json.NewDecoder(r.Body).Decode(&movie)

	insertOneMovie(movie)

	json.NewEncoder(w).Encode(movie)
}

func MarkedAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Origin", "*")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	params := mux.Vars(r)

	updateOneMovie(params["id"])

	json.NewEncoder(w).Encode(params["id"])
}

func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Origin", "*")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)

	deleteOneMovie(params["id"])

	json.NewEncoder(w).Encode(params["id"])

}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Origin", "*")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	deleteAllMovie()

	json.NewEncoder(w).Encode("All Movies Deleted")

}
