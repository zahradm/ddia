# Avro with Go

This example demonstrates how to use Avro with a Go HTTP server.

## What is Avro?

Apache Avro is a remote procedure call and data serialization framework developed within the Apache Hadoop project. It uses JSON for defining data types and protocols, and serializes data in a compact binary format.

### Key Features:

*   **Rich Data Structures:** Avro supports a rich set of primitive types (`null`, `boolean`, `int`, `long`, `float`, `double`, `bytes`, and `string`) and complex types (`record`, `enum`, `array`, `map`, `union`, and `fixed`).
*   **Schema Evolution:** This is one of Avro's most powerful features. It allows you to update the schema (e.g., add or remove fields) in a way that is both forward and backward compatible. When reading Avro data, the schema used to write the data is always present. This allows the reader to handle data written with an older or newer schema.
*   **Dynamic Typing:** Serialization and deserialization can happen without code generation. The schema is used at runtime to process the data.
*   **Integration:** It's a primary data format for Apache Kafka and is widely used in the Hadoop ecosystem.

## 1. Define the Avro Schema


Create a file named `user.avsc` with the following content:

```json
{
  "type": "record",
  "name": "User",
  "namespace": "main",
  "fields": [
    { "name": "id", "type": "int" },
    { "name": "name", "type": "string" },
    { "name": "email", "type": "string" }
  ]
}
```

## 2. Generate Go Code

First, install the `avrogen` tool:

```bash
go get github.com/hamba/avro/v2/cmd/avrogen
```

Then, generate the Go code from the `.avsc` file:

```bash
avrogen -pkg main -o user.go user.avsc
```

This will create `user.go`.

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
	"github.com/hamba/avro/v2"
)

var usersDb = make(map[int]*User)

func createUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	newUser := &User{}
	if err := avro.Unmarshal(newUser.Schema(), body, newUser); err != nil {
		http.Error(w, "Error unmarshaling avro", http.StatusBadRequest)
		return
	}

	if _, exists := usersDb[newUser.Id]; exists {
		http.Error(w, "User with this ID already exists", http.StatusBadRequest)
		return
	}

	usersDb[newUser.Id] = newUser
	fmt.Printf("Created user: %+v\n", newUser)

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

	foundUser, exists := usersDb[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	data, err := avro.Marshal(foundUser.Schema(), foundUser)
	if err != nil {
		http.Error(w, "Error marshaling avro", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/avro")
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
go mod init ddia/avro/go
go mod tidy
```

## 5. Run the Application

```bash
go run .
```

## 6. Test the API

You can use a client like `curl` or a Go script to test the API. You will need a serialized avro message to send. You can create one with the python script from the previous example.

### Create a user:

```bash
curl -X POST http://localhost:8080/users -H "Content-Type: application/avro" --data-binary "@../python-fastapi/user.bin"
```

### Get a user:

```bash
curl http://localhost:8080/users/1
```
The response will be a binary Avro message.
