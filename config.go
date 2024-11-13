package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
)

type Config struct {
    DatabaseURL  string
    ServerPort   string
    LogLevel     string // New configuration option for logging
}

func NewConfig() *Config {
    return &Config{
        DatabaseURL:  getEnv("DATABASE_URL", "defaultDatabaseURL"),
        ServerPort:   getEnv("SERVER_PORT", "4000"),
        LogLevel:     getEnv("LOG_LEVEL", "INFO"), // INFO as default
    }
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}

func (c *Config) validate() error {
    if c.DatabaseURL == "" {
        return fmt.Errorf("missing database URL")
    }
    if c.ServerPort == "" {
        return fmt.Errorf("missing server port")
    }
    // Could add more validation on LogLevel here if needed
    return nil
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.RequestURI)
        next.ServeHTTP(w, r)
    })
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome to GoForum!")
}

func main() {
    config := NewConfig()

    if err := config.validate(); err != nil {
        fmt.Printf("Configuration error: %s\n", err)
        os.Exit(1)
    }

    if config.LogLevel == "DEBUG" {
        log.SetFlags(log.LstdFlags | log.Lshortfile)
    }

    r := mux.NewRouter()
    // Apply our logging middleware to all routes
    r.Use(loggingMiddleware)
    r.HandleFunc("/", mainHandler)

    http.Handle("/", r)

    fmt.Println("Server started on port:", config.ServerPort)
    log.Fatal(http.ListenAndServe(":"+config.ServerPort, nil))
}