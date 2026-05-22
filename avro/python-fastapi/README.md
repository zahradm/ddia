# Avro with FastAPI

This example demonstrates how to use Avro with FastAPI.

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
  "fields": [
    { "name": "id", "type": "int" },
    { "name": "name", "type": "string" },
    { "name": "email", "type": "string" }
  ]
}
```

## 2. Create the FastAPI Application

Create a file named `main.py`:

```python
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

```

## 3. Create `requirements.txt`

Create a `requirements.txt` file:

```
fastapi
uvicorn
avro
```

## 4. Run the Application

Install the dependencies:

```bash
pip install -r requirements.txt
```

Run the FastAPI server:

```bash
uvicorn main:app --reload
```

## 5. Test the API

You can use a client like `curl` or a Python script to test the API.

### Create a user:

```bash
# First, create a serialized user object. You can use a python script for this.
# create_user_payload.py
import io
import json
from avro.schema import parse
from avro.io import DatumWriter, BinaryEncoder

with open("user.avsc", "r") as f:
    schema = parse(f.read())

writer = DatumWriter(schema)
bytes_writer = io.BytesIO()
encoder = BinaryEncoder(bytes_writer)
writer.write({"id": 1, "name": "test", "email": "test@example.com"}, encoder)
with open("user.bin", "wb") as f:
    f.write(bytes_writer.getvalue())


# Now send it with curl
curl -X POST http://127.0.0.1:8000/users/ -H "Content-Type: application/avro" --data-binary "@user.bin"
```

### Get a user:

```bash
curl http://127.0.0.1:8000/users/1
```
The response will be a binary Avro message. You would need to deserialize it to read the content.
