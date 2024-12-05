package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
)

type ForumThread struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
}

var forumThreads []ForumThread

func createForumThread(w http.ResponseWriter, r *http.Request) {
    var newForumThread ForumThread
    if err := json.NewDecoder(r.Body).Decode(&newForumThread); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if newForumThread.Title == "" || newForumThread.Description == "" {
        http.Error(w, "Missing title or description", http.StatusBadRequest)
        return
    }
    forumThreads = append(forumThreads, newForumThread)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newForumThread)
}

func retrieveForumThreads(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(forumThreads)
}

func findForumThreadByID(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, thread := range forumThreads {
        if thread.ID == params["id"] {
            json.NewEncoder(w).Encode(thread)
            return
        }
    }
    http.Error(w, "Thread not found", http.StatusNotFound)
}

func modifyForumThread(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for index, thread := range forumThreads {
        if thread.ID == params["id"] {
            var updatedThread ForumThread
            if err := json.NewDecoder(r.Body).Decode(&updatedThread); err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            if updatedThread.Title != "" {
                thread.Title = updatedThread.Title
            }
            if updatedThread.Description != "" {
                thread.Description = updatedThread.Description
            }
            forumThreads[index] = thread
            json.NewEncoder(w).Encode(thread)
            return
        }
    }
    http.Error(w, "Thread not found", http.StatusNotFound)
}

func removeForumThread(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for index, thread := range forumThreads {
        if thread.ID == params["id"] {
            forumThreads = append(forumThreads[:index], forumThreads[index+1:]...)
            break
        }
    }
    w.WriteHeader(http.StatusNoContent)
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/threads", createForumThread).Methods("POST")
    router.HandleFunc("/threads", retrieveForumThreads).Methods("GET")
    router.HandleFunc("/threads/{id}", findForumThreadByID).Methods("GET")
    router.HandleFunc("/threads/{id}", modifyForumThread).Methods("PUT")
    router.HandleFunc("/threads/{id}", removeForumThread).Methods("DELETE")

    httpAddr := os.Getenv("HTTP_ADDR")
    if httpAddr == "" {
        httpAddr = ":8000"
    }

    log.Fatal(http.ListenAndServe(httpAddr, router))
}