# Protobuf with gRPC in Python

This example demonstrates how to use Protocol Buffers with gRPC in Python. gRPC is a modern, high-performance RPC framework that is a natural fit for Protobuf.

## 1. Define the Protobuf Schema

We will define a `UserService` in our `.proto` file. Create a file named `user.proto` with the following content:

```proto
syntax = "proto3";

package user;

message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
}

message UserRequest {
  int32 id = 1;
}

message UserResponse {
  User user = 1;
}

message CreateUserRequest {
  User user = 1;
}

message CreateUserResponse {
  User user = 1;
}

service UserService {
  rpc GetUser(UserRequest) returns (UserResponse);
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
}
```

## 2. Generate Python Code

First, install the required tools:

```bash
pip install grpcio grpcio-tools
```

Then, generate the Python code from the `.proto` file:

```bash
python -m grpc_tools.protoc -I=. --python_out=. --grpc_python_out=. user.proto
```

This will create `user_pb2.py`, `user_pb2.pyi`, `user_pb2_grpc.py`, and `user_pb2_grpc.pyi`.

## 3. Create the gRPC Server

Create a file named `server.py`:

```python
from concurrent import futures
import grpc
import user_pb2
import user_pb2_grpc

users_db = {}

class UserService(user_pb2_grpc.UserServiceServicer):
    def GetUser(self, request, context):
        if request.id not in users_db:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details('User not found')
            return user_pb2.UserResponse()
        
        user = users_db[request.id]
        return user_pb2.UserResponse(user=user)

    def CreateUser(self, request, context):
        user = request.user
        if user.id in users_db:
            context.set_code(grpc.StatusCode.ALREADY_EXISTS)
            context.set_details('User with this ID already exists')
            return user_pb2.CreateUserResponse()
        
        users_db[user.id] = user
        print(f"Created user: {user}")
        return user_pb2.CreateUserResponse(user=user)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    user_pb2_grpc.add_UserServiceServicer_to_server(UserService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started on port 50051")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
```

## 4. Create the gRPC Client

Create a file named `client.py`:

```python
import grpc
import user_pb2
import user_pb2_grpc

def run():
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = user_pb2_grpc.UserServiceStub(channel)

        # Create a user
        new_user = user_pb2.User(id=1, name="gRPC User", email="grpc@example.com")
        create_response = stub.CreateUser(user_pb2.CreateUserRequest(user=new_user))
        print(f"Client received from create: {create_response.user}")

        # Get a user
        response = stub.GetUser(user_pb2.UserRequest(id=1))
        print(f"Client received from get: {response.user}")

if __name__ == '__main__':
    run()
```

## 5. Create `requirements.txt`

```
grpcio
grpcio-tools
```

## 6. Run the Example

First, install dependencies:
```bash
python3 -m venv env
source env/bin/activate
pip install -r requirements.txt
python -m grpc_tools.protoc -I=. --python_out=. --grpc_python_out=. user.proto
```

Start the server:
```bash
python server.py
```

In another terminal, run the client:
```bash
python client.py
```
