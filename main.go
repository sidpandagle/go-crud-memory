package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var users []User
var nextId int = 1

func main() {

	http.HandleFunc("/", home)
	http.HandleFunc("/user/update", updateUser)
	http.HandleFunc("/user/delete", deleteUser)
	http.HandleFunc("/user", getUser)
	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/user/create", createUser)
	err := http.ListenAndServe(":8000", nil)
	log.Fatal(err)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	log.Println("Hello World!")
	fmt.Fprintln(w, "Hello World!")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Ällow", "POST")
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	var newUser User
	newUser.Id = nextId
	nextId += 1

	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users = append(users, newUser)

	log.Println("Creating user")
	fmt.Fprintln(w, json.NewEncoder(w).Encode(newUser))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}
	for i := 0; i < len(users); i++ {
		if users[i].Id == id {
			json.NewEncoder(w).Encode(users[i])
			return
		}
	}
	http.NotFound(w, r)

}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}
	for i := 0; i < len(users); i++ {
		if users[i].Id == id {
			users = append(users[:i], users[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}
	var updateUser User
	err = json.NewDecoder(r.Body).Decode(&updateUser)
	for i := 0; i < len(users); i++ {
		if users[i].Id == id {
			users[i] = updateUser
			return
		}
	}
	http.NotFound(w, r)

}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// // // Create User
// curl -X POST http://localhost:8000/user/create ^
// -H "Content-Type: application/json" ^
// -d "{\"name\": \"John Doe\"}"

// // // Get Users
// curl -X GET http://localhost:8000/users

// // // Get User by ID
// curl -X GET http://localhost:8000/users/{id}

// // // Update User
// curl -X PUT http://localhost:8000/user/update?id=1 ^
// -H "Content-Type: application/json" ^
// -d "{\"name\": \"Jane Doe\"}"

// // // Delete User
// curl -X DELETE http://localhost:8000/users/{id}
