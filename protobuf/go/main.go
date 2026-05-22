package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"ddia/protobuf/go/user"
	"github.com/gorilla/mux"
	"google.golang.org/protobuf/proto"
)

var usersDb = make(map[int32]*user.User)

func createUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	newUser := &user.User{}
	if err := proto.Unmarshal(body, newUser); err != nil {
		http.Error(w, "Error unmarshaling protobuf", http.StatusBadRequest)
		return
	}

	if _, exists := usersDb[newUser.Id]; exists {
		http.Error(w, "User with this ID already exists", http.StatusBadRequest)
		return
	}

	usersDb[newUser.Id] = newUser
	fmt.Printf("Created user: %v\n", newUser)

	w.WriteHeader(http.StatusCreated)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	foundUser, exists := usersDb[int32(id)]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	data, err := proto.Marshal(foundUser)
	if err != nil {
		http.Error(w, "Error marshaling protobuf", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/protobuf")
	w.Write(data)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}
