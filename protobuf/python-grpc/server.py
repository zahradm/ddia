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
