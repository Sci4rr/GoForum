package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type User struct {
    gorm.Model
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

func initDB() {
    db, err = gorm.Open("sqlite3", "forum.db")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    db.AutoMigrate(&User{})
}

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func generateJWT(email string) (string, error) {
    expirationTime := time.Now().Add(30 * time.Minute)
    claims := &Claims{
        Email: email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    hashedPassword, err := hashPassword(user.Password)
    if err != nil {
        http.Error(w, "Error while hashing password", http.StatusInternalServerError)
        return
    }
    user.Password = hashedPassword

    var exists User
    db.Where("username = ?", user.Username).Or("email = ?", user.Email).First(&exists)
    if exists.ID != 0 {
        http.Error(w, "Username or email already exists", http.StatusBadRequest)
        return
    }

    result := db.Create(&user)
    if result.Error != nil {
        http.Error(w, "Error creating user", http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    var foundUser User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    db.Where("email = ?", user.Email).First(&foundUser)
    if foundUser.ID == 0 || !checkPasswordHash(user.Password, foundUser.Password) {
        http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
        return
    }

    tokenString, err := generateJWT(user.Email)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    w.Write([]byte(fmt.Sprintf("Bearer %s", tokenString)))
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
    email := r.URL.Query().Get("email")
    var user User
    db.Where("email = ?", email).First(&user)
    if user.ID == 0 {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func main() {
    initDB()

    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/login", loginHandler)
    http.HandleFunc("/profile", profileHandler)

    log.Fatal(http.ListenAndServe(":8080", nil))
}