<p align="center">
  <img src="assets/mocking-bird-banner.svg" alt="Mocking Bird" width="800"/>
</p>

<p align="center">
  <strong>A fully-featured mock API server supporting 7 protocols out of the box</strong>
</p>

<p align="center">
  <a href="#-quick-start">Quick Start</a> •
  <a href="#-supported-protocols">Protocols</a> •
  <a href="#-crud-testing">CRUD Testing</a> •
  <a href="#-examples">Examples</a>
</p>

<br/>

---

## What is Mocking Bird?

**Mocking Bird** is a standalone mock API server that speaks every major API protocol. Use it for development, testing, or as a reference implementation for multi-protocol API tooling.

### ✨ Features

- 🎯 **7 Protocols Supported**: OpenAPI, Swagger, SOAP/WSDL, GraphQL, OData, gRPC, JSON-RPC
- 🔄 **Full CRUD Operations**: Create, Read, Update, Delete for all REST-like APIs
- 💾 **In-Memory Persistence**: Data persists during server runtime (no database needed)
- 🔐 **Multiple Auth Types**: Bearer tokens, Basic auth, API keys
- 📝 **Comprehensive Specs**: Each protocol includes spec endpoints for auto-discovery
- 🧪 **Test Suite Included**: Automated CRUD tests for all endpoints

---

## 📦 Supported Protocols

| Protocol | Resource | Description | Auth | Port/Path |
|----------|----------|-------------|------|-----------|
| **OpenAPI 3.0** | Pets | Pet store with CRUD operations | None | `:9999/openapi` |
| **Swagger 2.0** | Dinosaurs | Dinosaur catalog with CRUD | Bearer: `dino-token` | `:9999/swagger` |
| **WSDL/SOAP** | Plants | Plant database with CRUD | Bearer: `mock-token` | `:9999/wdsl/soap` |
| **GraphQL** | Cars | Car inventory with CRUD | Basic: `graphql-user:graphql-pass` | `:9999/graphql` |
| **OData v4** | Movies | Movie database with query support | None | `:9999/odata` |
| **gRPC** | Clothes | Clothing items (4 categories) | None | `:50051-50054` |
| **JSON-RPC** | Calculator | Math operations | None | `:9999/jsonrpc` |

---

## 🚀 Quick Start

### Install

```bash
git clone git@github.com:emadomedher/mocking-bird.git
cd mocking-bird
go build -o mockingbird .
```

### Run

```bash
./mockingbird
```

Server starts on `http://localhost:9999` (gRPC on ports `50051-50054`)

### Test

```bash
./test-crud.sh
```

Runs automated CRUD tests for all APIs and verifies in-memory persistence.

---

## 📚 API Documentation

### OpenAPI (Pets)

**Spec**: `http://localhost:9999/openapi/openapi.json`

**Endpoints**:
```bash
GET    /openapi/pets           # List pets
POST   /openapi/pets           # Create pet
GET    /openapi/pets/{id}      # Get pet by ID
PUT    /openapi/pets/{id}      # Update pet
DELETE /openapi/pets/{id}      # Delete pet
```

**Example**:
```bash
curl -X POST http://localhost:9999/openapi/pets \
  -H "Content-Type: application/json" \
  -d '{"name":"Fluffy"}'
```

---

### Swagger 2.0 (Dinosaurs)

**Spec**: `http://localhost:9999/swagger/swagger.json`

**Endpoints**:
```bash
GET    /swagger/dinosaurs           # List dinosaurs
POST   /swagger/dinosaurs           # Create dinosaur
GET    /swagger/dinosaurs/{id}      # Get dinosaur by ID
PUT    /swagger/dinosaurs/{id}      # Update dinosaur
DELETE /swagger/dinosaurs/{id}      # Delete dinosaur
```

**Authentication**: `Authorization: Bearer dino-token`

**Example**:
```bash
curl http://localhost:9999/swagger/dinosaurs \
  -H "Authorization: Bearer dino-token"
```

---

### WSDL/SOAP (Plants)

**WSDL Spec**: `http://localhost:9999/wdsl/wsdl`

**Operations**: `ListPlants`, `GetPlant`, `CreatePlant`, `UpdatePlant`, `DeletePlant`

**Authentication**: `Authorization: Bearer mock-token`

**Example**:
```bash
curl -X POST http://localhost:9999/wdsl/soap \
  -H "Authorization: Bearer mock-token" \
  -H "Content-Type: text/xml" \
  -d '<?xml version="1.0"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ListPlants xmlns="http://example.com/plants"/>
  </soap:Body>
</soap:Envelope>'
```

---

### GraphQL (Cars)

**Schema (SDL)**: `http://localhost:9999/graphql/schema`

**Endpoint**: `POST http://localhost:9999/graphql`

**Authentication**: Basic auth (`graphql-user:graphql-pass`)

**Queries**:
```graphql
query {
  listCars(limit: 10) {
    id
    name
  }
}

query {
  getCar(id: "1") {
    id
    name
  }
}
```

**Mutations**:
```graphql
mutation {
  createCar(name: "Tesla") {
    id
    name
  }
}

mutation {
  updateCar(id: "1", name: "Tesla Model S") {
    id
    name
  }
}

mutation {
  deleteCar(id: "1")
}
```

**Example**:
```bash
curl -X POST http://localhost:9999/graphql \
  -u graphql-user:graphql-pass \
  -H "Content-Type: application/json" \
  -d '{"query":"query{listCars(limit:5){id name}}"}'
```

---

### OData v4 (Movies)

**Metadata**: `http://localhost:9999/odata/$metadata`

**Endpoints**:
```bash
GET    /odata/Movies              # List movies with query options
POST   /odata/Movies              # Create movie
GET    /odata/Movies(1)           # Get movie by ID
PUT    /odata/Movies(1)           # Update movie
DELETE /odata/Movies(1)           # Delete movie
```

**Query Options**: `$filter`, `$orderby`, `$top`, `$skip`, `$select`, `$expand`

**Example**:
```bash
# List movies with filtering
curl "http://localhost:9999/odata/Movies?\$filter=Year gt 2000&\$orderby=Rating desc"

# Create a movie
curl -X POST http://localhost:9999/odata/Movies \
  -H "Content-Type: application/json" \
  -d '{"title":"Inception","year":2010,"genre":"Sci-Fi","rating":8.8}'
```

---

### gRPC (Clothes)

**Proto File**: `clothes.proto`

**Service**: `ClothesService`

**Categories & Ports**:
- **Hats** - `localhost:50051`
- **Shoes** - `localhost:50052`
- **Pants** - `localhost:50053`
- **Shirts** - `localhost:50054`

**Method**: `ListClothes(ListClothesRequest) returns (ListClothesResponse)`

**Example** (using `grpcurl`):
```bash
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 ClothesService/ListClothes
```

---

### JSON-RPC (Calculator)

**OpenRPC Spec**: `http://localhost:9999/jsonrpc/openrpc.json`

**Methods**: `add`, `subtract`, `multiply`, `divide`, `rpc.discover`

**Example**:
```bash
curl -X POST http://localhost:9999/jsonrpc \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"add","params":{"a":5,"b":3}}'

# Response: {"jsonrpc":"2.0","id":1,"result":8}
```

---

## 🧪 CRUD Testing

Run the comprehensive test suite to verify all CRUD operations:

```bash
./test-crud.sh
```

**What it tests**:
- ✅ Create operations for all APIs
- ✅ Read operations (list and get by ID)
- ✅ Update operations with data validation
- ✅ Delete operations with proper cleanup
- ✅ In-memory persistence across requests

**Sample output**:
```
🧪 Testing CRUD operations for all Mock APIs
==============================================

Testing Pets (OpenAPI)
✓ Created pet with ID: 4
✓ Read pet successfully
✓ Updated pet successfully
✓ Deleted pet successfully

Testing Cars (GraphQL)
✓ Created car with ID: 3
✓ Read car successfully
✓ Updated car successfully
✓ Deleted car successfully

✅ All CRUD tests completed!
```

---

## 🎯 Use Cases

### Local Development
Replace external API dependencies with Mocking Bird during development. No network calls, instant responses.

### Integration Testing
Test your API clients against all supported protocols without setting up real services.

### CI/CD Pipelines
Run Mocking Bird in Docker for automated testing in your build pipeline.

### Protocol Learning
Use as a reference implementation to learn how different API protocols work.

### Multi-Protocol Tools
Use as a test fixture when building tools that need to support multiple API formats (like [Skyline MCP](https://github.com/emadomedher/skyline-mcp)).

---

## 🐳 Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o mockingbird .

FROM alpine:latest
COPY --from=builder /app/mockingbird /usr/local/bin/
COPY --from=builder /app/clothes.proto /app/
EXPOSE 9999 50051 50052 50053 50054
CMD ["mockingbird"]
```

Build and run:
```bash
docker build -t mockingbird .
docker run -p 9999:9999 -p 50051-50054:50051-50054 mockingbird
```

---

## 🛠️ Configuration

**Environment Variables**:
```bash
# Dinosaur auth token
export DINOSAURS_SWAGGER2_TOKEN="custom-dino-token"

# GraphQL credentials
export GRAPHQL_USERNAME="admin"
export GRAPHQL_PASSWORD="secret"

# gRPC base port (default: 50051)
export GRPC_BASE_PORT=60051
```

---

## 📖 Seed Data

Each API comes pre-seeded with sample data:

| API | Seed Data |
|-----|-----------|
| Pets | 3 pets (pet-1, pet-2, pet-3) |
| Dinosaurs | 2 dinosaurs (t-rex, triceratops) |
| Plants | 2 plants (fern, cactus) |
| Cars | 2 cars (sedan, truck) |
| Movies | 8 movies (Matrix, Inception, etc.) |
| Clothes | 3 items per category (12 total) |

Data persists in memory for the lifetime of the server process.

---

## 🤝 Contributing

Contributions welcome! Ideas:
- Add more protocols (e.g., Thrift, Avro)
- Add persistence options (Redis, PostgreSQL)
- Add WebSocket support
- Add rate limiting examples
- Improve test coverage

---

## 📄 License

MIT License - see [LICENSE](LICENSE) for details

---

## 🔗 Related Projects

- **[Skyline MCP](https://github.com/emadomedher/skyline-mcp)** - Turn any API into MCP tools for AI agents
- Uses Mocking Bird as its test fixture for multi-protocol support

---

<p align="center">
  <img src="assets/mocking-bird-logo.svg" alt="Mocking Bird" width="200"/>
</p>

<p align="center">
  <sub>Built with Go. Mock everything.</sub>
</p>
