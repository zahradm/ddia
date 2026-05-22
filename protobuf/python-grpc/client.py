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
