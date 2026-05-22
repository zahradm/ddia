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
