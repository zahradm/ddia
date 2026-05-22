# Protobuf with FastAPI (REST)

This example demonstrates how to use Protocol Buffers (Protobuf) with FastAPI over a RESTful interface. While this is possible, it's more common to use Protobuf with an RPC framework like gRPC. For a gRPC example, see the `python-grpc` directory.

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

message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
}
```

## 2. Generate Python Code

First, install the required tools:

```bash
pip install grpcio-tools
```

Then, generate the Python code from the `.proto` file:

```bash
python -m grpc_tools.protoc -I=. --python_out=. --pyi_out=. user.proto
```

This will create `user_pb2.py` and `user_pb2.pyi`.

## 3. Create the FastAPI Application

Create a file named `main.py`:

```python
from fastapi import FastAPI, HTTPException
from fastapi.responses import Response
from user_pb2 import User

app = FastAPI()

users_db = {}

@app.post("/users/", response_class=Response)
async def create_user(user_proto: bytes):
    user = User()
    try:
        user.ParseFromString(user_proto)
    except Exception as e:
        raise HTTPException(status_code=400, detail="Invalid Protobuf message")

    if user.id in users_db:
        raise HTTPException(status_code=400, detail="User with this ID already exists")

    users_db[user.id] = user
    return Response(content=user.SerializeToString(), media_type="application/protobuf")

@app.get("/users/{user_id}", response_class=Response)
async def get_user(user_id: int):
    if user_id not in users_db:
        raise HTTPException(status_code=404, detail="User not found")
    
    user = users_db[user_id]
    return Response(content=user.SerializeToString(), media_type="application/protobuf")

```

## 4. Create `requirements.txt`

Create a `requirements.txt` file:

```
fastapi
uvicorn
protobuf
grpcio-tools
```

## 5. Run the Application

Install the dependencies:

```bash
pip install -r requirements.txt
```

Run the FastAPI server:

```bash
uvicorn main:app --reload
```

## 6. Test the API

You can use a client like `curl` or a Python script to test the API.

### Create a user:

```bash
# First, create a serialized user object. You can use a python script for this.
# create_user_payload.py
from user_pb2 import User
user = User(id=1, name="test", email="test@example.com")
with open("user.bin", "wb") as f:
    f.write(user.SerializeToString())

# Now send it with curl
curl -X POST http://127.0.0.1:8000/users/ -H "Content-Type: application/protobuf" --data-binary "@user.bin"
```

### Get a user:

```bash
curl http://127.0.0.1:8000/users/1
```
The response will be a binary Protobuf message. You would need to deserialize it to read the content.
