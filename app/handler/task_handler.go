package handler

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

type Handler struct {
	Client *mongo.Client
}

func NewHandler(client *mongo.Client) Handler {
	return Handler{Client: client}
}

// Task структура для зберігання завдань
type Task struct {
	UUID     string    `json:"uuid" bson:"uuid"`
	Count    int       `json:"count" bson:"count"`
	Type     string    `json:"type" bson:"type"`
	Status   string    `json:"status" bson:"status"`
	Subtasks []SubTask `json:"subtasks" bson:"subtasks"`
}

// SubTask структура для зберігання підзавдань
type SubTask struct {
	UUID   string `json:"uuid" bson:"uuid"`
	Type   string `json:"type" bson:"type"`
	Status string `json:"status" bson:"status"`
}

// createTask обробник запиту POST /createTask
func (h Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)

	task.Status = "created"
	task.UUID = primitive.NewObjectID().Hex()
	task.Subtasks = []SubTask{}

	collection := h.Client.Database("micro-tasks").Collection("tasks")
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
func (h Handler) AddSubTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var subTask SubTask
	_ = json.NewDecoder(r.Body).Decode(&subTask)

	uuid := r.URL.Query().Get("uuid")
	log.Println(uuid)
	collection := h.Client.Database("micro-tasks").Collection("tasks")
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

func (h Handler) CheckStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uuid := r.FormValue("uuid")

	collection := h.Client.Database("micro-tasks").Collection("tasks")
	filter := bson.M{"uuid": uuid}

	var task Task
	err := collection.FindOne(context.Background(), filter).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to fetch task")
		return
	}

	response := map[string]interface{}{
		"uuid":   task.UUID,
		"status": task.Status,
	}

	json.NewEncoder(w).Encode(response)
}
