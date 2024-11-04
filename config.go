package main

import (
    "fmt"
    "os"
)

type Config struct {
    DatabaseURL string
    ServerPort  string
}

func NewConfig() *Config {
    return &Config{
        DatabaseURL: getEnv("DATABASE_URL", "defaultDatabaseURL"),
        ServerPort:  getEnv("SERVER_PORT", "4000"),
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

    return nil
}

func main() {
    config := NewConfig()

    if err := config.validate(); err != nil {
        fmt.Printf("Configuration error: %s\n", err)
        os.Exit(1)
    }

    fmt.Println("Server started on port:", config.ServerPort)
}