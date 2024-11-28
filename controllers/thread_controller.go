package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Thread struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

var threads []Thread

func createThread(w http.ResponseWriter, r *http.Request) {
	var newThread Thread
	if err := json.NewDecoder(r.Body).Decode(&newThread); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if newThread.Title == "" || newThread.Description == "" {
		http.Error(w, "Missing title or description", http.StatusBadRequest)
		return
	}
	threads = append(threads, newThread)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newThread)
}

func getThreads(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(threads)
}

func getThreadByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range threads {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Thread not found", http.StatusNotFound)
}

func updateThread(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range threads {
		if item.ID == params["id"] {
			var updatedThread Thread
			if err := json.NewDecoder(r.Body).Decode(&updatedThread); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if updatedThread.Title != "" {
				item.Title = updatedThread.Title
			}
			if updatedThread.Description != "" {
				item.Description = updatedThread.Description
			}
			threads[index] = item
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Thread not found", http.StatusNotFound)
}

func deleteThread(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range threads {
		if item.ID == params["id"] {
			threads = append(threads[:index], threads[index+1:]...)
			break
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/threads", createThread).Methods("POST")
	router.HandleFunc("/threads", getThreads).Methods("GET")
	router.HandleFunc("/threads/{id}", getThreadByID).Methods("GET")
	router.HandleFunc("/threads/{id}", updateThread).Methods("PUT")
	router.HandleFunc("/threads/{id}", deleteThread).Methods("DELETE")

	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8000"
	}

	log.Fatal(http.ListenAndServe(httpAddr, router))
}