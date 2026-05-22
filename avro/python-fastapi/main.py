import io
import json
from fastapi import FastAPI, HTTPException
from fastapi.responses import Response
from avro.schema import parse
from avro.io import DatumWriter, DatumReader, BinaryEncoder, BinaryDecoder

app = FastAPI()

users_db = {}

# Load Avro schema
with open("user.avsc", "r") as f:
    schema_str = f.read()
schema = parse(schema_str)

@app.post("/users/", response_class=Response)
async def create_user(user_avro: bytes):
    bytes_reader = io.BytesIO(user_avro)
    decoder = BinaryDecoder(bytes_reader)
    reader = DatumReader(schema)
    try:
        user_data = reader.read(decoder)
    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Invalid Avro message: {e}")

    if user_data['id'] in users_db:
        raise HTTPException(status_code=400, detail="User with this ID already exists")

    users_db[user_data['id']] = user_data

    # Serialize response
    bytes_writer = io.BytesIO()
    encoder = BinaryEncoder(bytes_writer)
    writer = DatumWriter(schema)
    writer.write(user_data, encoder)
    return Response(content=bytes_writer.getvalue(), media_type="application/avro")

@app.get("/users/{user_id}", response_class=Response)
async def get_user(user_id: int):
    if user_id not in users_db:
        raise HTTPException(status_code=404, detail="User not found")
    
    user = users_db[user_id]
    
    # Serialize response
    bytes_writer = io.BytesIO()
    encoder = BinaryEncoder(bytes_writer)
    writer = DatumWriter(schema)
    writer.write(user, encoder)
    return Response(content=bytes_writer.getvalue(), media_type="application/avro")
