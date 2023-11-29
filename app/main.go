package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Task структура для зберігання завдань
type Task struct {
	UUID     string      `json:"uuid" bson:"uuid"`
	Count    int         `json:"count" bson:"count"`
	Type     string      `json:"type" bson:"type"`
	Status   string      `json:"status" bson:"status"`
	Subtasks []SubTask   `json:"subtasks" bson:"subtasks"`
}

// SubTask структура для зберігання підзавдань
type SubTask struct {
	UUID   string `json:"uuid" bson:"uuid"`
	Type   string `json:"type" bson:"type"`
	Status string `json:"status" bson:"status"`
}

var client *mongo.Client

// createTask обробник запиту POST /createTask
func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)

	task.Status = "created"
	task.UUID = primitive.NewObjectID().Hex()
	task.Subtasks = []SubTask{}

	collection := client.Database("micro-tasks").Collection("tasks")
	_, err := collection.InsertOne(context.Background(), task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to create task")
		return
	}

	response := map[string]string{"uuid": task.UUID}
	json.NewEncoder(w).Encode(response)
}

// addSubTask обробник запиту POST /addSubTask
func addSubTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var subTask SubTask
	_ = json.NewDecoder(r.Body).Decode(&subTask)

	uuid := r.URL.Query().Get("uuid")
    log.Println(uuid)
	collection := client.Database("micro-tasks").Collection("tasks")
	filter := bson.M{"uuid": uuid}
	update := bson.M{
		"$push": bson.M{"subtasks": subTask},
	}

	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
	    log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to add subtask")
		return
	}
	log.Println(res)

	// Перевіряємо кількість підзавдань та змінюємо статус, якщо потрібно
	var task Task
	err = collection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
	    log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to update task status!")
		return
	}

	if len(task.Subtasks) == task.Count {
		update := bson.M{
			"$set": bson.M{"status": "done"},
		}

		_, err := collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("Failed to update task status")
			return
		}
	}

	response := map[string]string{"uuid": subTask.UUID}
	json.NewEncoder(w).Encode(response)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/createTask", createTask).Methods("POST")
	router.HandleFunc("/addSubTask", addSubTask).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
