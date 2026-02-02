package main

import (
	"Proj_2/taskstore"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type taskServer struct {
	store *taskstore.TaskStore
}

func NewTaskServer() *taskServer {
	store := taskstore.New()
	return &taskServer{store: store}
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling task create at %s\n", r.URL.Path)

	type RequestTask struct {
		Text string `json:"text"`
		Tags []string `json:"tags"`
		Due time.Time `json:"due"`
	}

	type ResponseId struct {
		Id int `json:"id"`
	}

	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "except application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var rt RequestTask
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags, rt.Due)
	renderJSON(w, ResponseId{Id: id})
}

func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get all tasks at %s\n", r.URL.Path)

	allTasks := ts.store.GetAllTasks()
	renderJSON(w, allTasks)
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling get task at %s\n", r.URL.Path)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	renderJSON(w, task)
}

func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling delete task at %s\n", r.URL.Path)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (ts *taskServer) deleteAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling delete all tasks at %s\n", r.URL.Path)
	ts.store.DeleteAllTasks()
}

func (ts *taskServer) tagHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling tasks by tag at %s\n", r.URL.Path)

	tag := r.PathValue(("tag"))

	tasks := ts.store.GetTaskByTag(tag)
	renderJSON(w, tasks)
}

func (ts *taskServer) dueHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling tasks by due at %s\n", r.URL.Path)

	badRequestError := func() {
		http.Error(w, fmt.Sprintf("except /due/<year>/<month>/<day>, got %v", r.URL.Path), http.StatusBadRequest)
	}

	year, errYear := strconv.Atoi(r.PathValue("year"))
	month, errMonth := strconv.Atoi(r.PathValue("month"))
	day, errDay := strconv.Atoi(r.PathValue("day"))
	if errYear != nil || errMonth != nil || errDay != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}

	tasks := ts.store.GetTaskByDueData(year, time.Month(month), day)
	renderJSON(w, tasks)
}

func coreHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	server := NewTaskServer()

	router.HandleFunc("/", coreHandler)

	router.HandleFunc("/task/", server.createTaskHandler).Methods("POST")
	router.HandleFunc("/task/", server.getAllTasksHandler).Methods("GET")
	router.HandleFunc("/task/", server.deleteAllTasksHandler).Methods("DELETE")
	router.HandleFunc("/task/{id}/", server.getTaskHandler).Methods("GET")
	router.HandleFunc("/task/{id}/", server.deleteTaskHandler).Methods("DELETE")
	router.HandleFunc("/tag/{tag}/", server.tagHandler).Methods("GET")
	router.HandleFunc("/due/{year}/{month}/{day}/", server.dueHandler).Methods("GET")
	
	port := "8080"
	log.Printf("Сервер запущен на http://localhost:%s", port)
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}