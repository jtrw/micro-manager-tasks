package main

import (
"context"
//"encoding/json"
"os"
"os/signal"
"syscall"
"log"
"time"
"github.com/jessevdk/go-flags"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
server "micro-manager-tasks/m/v2/app/server"
)

type Options struct {
Listen string `short:"l" long:"listen" env:"LISTEN" default:":8080" description:"listen address"`
Secret string `short:"s" long:"secret" env:"EVENT_SECRET_KEY" default:"123"`
PinSize int `long:"pinszie" env:"PIN_SIZE" default:"5" description:"pin size"`
MaxExpire time.Duration `long:"expire" env:"MAX_EXPIRE" default:"24h" description:"max lifetime"`
MaxPinAttempts int `long:"pinattempts" env:"PIN_ATTEMPTS" default:"3" description:"max attempts to enter pin"`
WebRoot string `long:"web" env:"WEB" default:"/" description:"web ui location"`
}

// // Task структура для зберігання завдань
// type Task struct {
// 	UUID     string      `json:"uuid" bson:"uuid"`
// 	Count    int         `json:"count" bson:"count"`
// 	Type     string      `json:"type" bson:"type"`
// 	Status   string      `json:"status" bson:"status"`
// 	Subtasks []SubTask   `json:"subtasks" bson:"subtasks"`
// }
//
// // SubTask структура для зберігання підзавдань
// type SubTask struct {
// 	UUID   string `json:"uuid" bson:"uuid"`
// 	Type   string `json:"type" bson:"type"`
// 	Status string `json:"status" bson:"status"`
// }

//var client *mongo.Client
var revision string

// // createTask обробник запиту POST /createTask
// func createTask(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var task Task
// 	_ = json.NewDecoder(r.Body).Decode(&task)
//
// 	task.Status = "created"
// 	task.UUID = primitive.NewObjectID().Hex()
// 	task.Subtasks = []SubTask{}
//
// 	collection := client.Database("micro-tasks").Collection("tasks")
// 	_, err := collection.InsertOne(context.Background(), task)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode("Failed to create task")
// 		return
// 	}
//
// 	response := map[string]string{"uuid": task.UUID}
// 	json.NewEncoder(w).Encode(response)
// }
//
// // addSubTask обробник запиту POST /addSubTask
// func addSubTask(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var subTask SubTask
// 	_ = json.NewDecoder(r.Body).Decode(&subTask)
//
// 	uuid := r.URL.Query().Get("uuid")
//     log.Println(uuid)
// 	collection := client.Database("micro-tasks").Collection("tasks")
// 	filter := bson.M{"uuid": uuid}
// 	update := bson.M{
// 		"$push": bson.M{"subtasks": subTask},
// 	}
//
// 	res, err := collection.UpdateOne(context.Background(), filter, update)
// 	if err != nil {
// 	    log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode("Failed to add subtask")
// 		return
// 	}
// 	log.Println(res)
//
// 	// Перевіряємо кількість підзавдань та змінюємо статус, якщо потрібно
// 	var task Task
// 	err = collection.FindOne(context.Background(), filter).Decode(&task)
// 	if err != nil {
// 	    log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode("Failed to update task status!")
// 		return
// 	}
//
// 	if len(task.Subtasks) == task.Count {
// 		update := bson.M{
// 			"$set": bson.M{"status": "done"},
// 		}
//
// 		_, err := collection.UpdateOne(context.Background(), filter, update)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			json.NewEncoder(w).Encode("Failed to update task status")
// 			return
// 		}
// 	}
//
// 	response := map[string]string{"uuid": subTask.UUID}
// 	json.NewEncoder(w).Encode(response)
// }
//
// func checkStatus(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	uuid := r.FormValue("uuid")
//
// 	collection := client.Database("your_database_name").Collection("tasks")
// 	filter := bson.M{"uuid": uuid}
//
// 	var task Task
// 	err := collection.FindOne(context.Background(), filter).Decode(&task)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode("Failed to fetch task")
// 		return
// 	}
//
// 	response := map[string]interface{}{
// 		"uuid":   task.UUID,
// 		"status": task.Status,
// 	}
//
// 	json.NewEncoder(w).Encode(response)
// }

func main() {
log.Printf("Micro Manager tasks %s\n", revision)

var opts Options
parser := flags.NewParser(&opts, flags.Default)
_, err := parser.Parse()
if err != nil {
log.Fatal(err)
}

ctx, cancel := context.WithCancel(context.Background())
go func() {
if x := recover(); x != nil {
log.Printf("[WARN] run time panic:\n%v", x)
panic(x)
}

// catch signal and invoke graceful termination
stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
<-stop
log.Printf("[WARN] interrupt signal")
cancel()
}()

clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
client, err := mongo.Connect(ctx, clientOptions)
if err != nil {
log.Fatal(err)
}

srv := server.Server{
Listen:         opts.Listen,
PinSize:        opts.PinSize,
MaxExpire:      opts.MaxExpire,
MaxPinAttempts: opts.MaxPinAttempts,
WebRoot:        opts.WebRoot,
Secret:         opts.Secret,
Version:        revision,
Client:         client,
}
if err := srv.Run(ctx); err != nil {
log.Printf("[ERROR] failed, %+v", err)
}


// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
//
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
// 	var err error
// 	client, err = mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	router := mux.NewRouter()
//
// 	router.HandleFunc("/createTask", createTask).Methods("POST")
// 	router.HandleFunc("/addSubTask", addSubTask).Methods("POST")
// 	router.HandleFunc("/checkStatus", checkStatus).Methods("GET")
//
// 	log.Fatal(http.ListenAndServe(":8080", router))
}
