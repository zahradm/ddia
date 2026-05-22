# Protobuf with Go (REST)

This example demonstrates how to use Protocol Buffers (Protobuf) with a Go HTTP server for a RESTful interface. While this is possible, it's more common to use Protobuf with an RPC framework like gRPC. For a gRPC example, see the `go-grpc` directory.

## What is Protobuf?

Protocol Buffers is a free and open-source cross-platform data format used to serialize structured data. It is useful in developing programs to communicate with each other over a network or for storing data.

### Key Features:

*   **Language-agnostic:** You define your data structure once in a `.proto` file, and then you can generate source code for various languages (like Python, Go, Java, C++, etc.).
*   **Efficient:** Protobuf serialization is binary, which makes it smaller and faster to send over the network compared to text-based formats like JSON or XML.
*   **Strongly-typed:** The schema enforces data types, which helps prevent errors.
*   **Schema Evolution:** You can update your data structures without breaking existing services that are built on the old format. Protobuf is backward-compatible.

## 1. Define the Protobuf Schema


Create a file named `user.proto` with the following content:

```proto
syntax = "proto3";

package user;

option go_package = ".;user";

message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
}
```

## 2. Generate Go Code

First, install the `protoc` compiler and the Go plugins:

```bash
# Install protoc (if you don't have it)
# On macOS: brew install protobuf
# On Linux: sudo apt-get install -y protobuf-compiler

# Install the Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

Make sure your `$GOPATH/bin` is in your `PATH`.

Then, generate the Go code from the `.proto` file:

```bash
protoc --go_out=. --go_opt=paths=source_relative user.proto
```

This will create `user.pb.go`.

## 3. Create the Go Application

Create a file named `main.go`:

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"google.golang.org/protobuf/proto"
	"ddia/protobuf/go/user"
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
```

## 4. Create `go.mod` and `go.sum`

Run the following commands to create the module files and download dependencies:

```bash
go mod init ddia/protobuf/go
go mod tidy
```

## 5. Run the Application

```bash
go run main.go
```

## 6. Test the API

You can use a client like `curl` or a Go script to test the API. You will need a serialized protobuf message to send. You can create one with the python script from the previous example.

### Create a user:

```bash
curl -X POST http://localhost:8080/users -H "Content-Type: application/protobuf" --data-binary "@../python-fastapi/user.bin"
```

### Get a user:

```bash
curl http://localhost:8080/users/1
```
The response will be a binary Protobuf message.
