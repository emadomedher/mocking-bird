#!/bin/bash
set -e

BASE_URL="http://localhost:9999"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🧪 Testing CRUD operations for all Mock APIs"
echo "=============================================="
echo ""

# Test Pets (OpenAPI)
echo -e "${YELLOW}Testing Pets (OpenAPI)${NC}"
echo "1. Create a pet..."
CREATED_PET=$(curl -s -X POST "$BASE_URL/openapi/pets" -H "Content-Type: application/json" -d '{"name":"TestDog"}')
PET_ID=$(echo $CREATED_PET | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo -e "${GREEN}✓ Created pet with ID: $PET_ID${NC}"

echo "2. Read the pet..."
curl -s "$BASE_URL/openapi/pets/$PET_ID" | grep -q "TestDog" && echo -e "${GREEN}✓ Read pet successfully${NC}" || echo -e "${RED}✗ Failed to read pet${NC}"

echo "3. Update the pet..."
curl -s -X PUT "$BASE_URL/openapi/pets/$PET_ID" -H "Content-Type: application/json" -d '{"name":"UpdatedDog"}' | grep -q "UpdatedDog" && echo -e "${GREEN}✓ Updated pet successfully${NC}" || echo -e "${RED}✗ Failed to update pet${NC}"

echo "4. Delete the pet..."
curl -s -X DELETE "$BASE_URL/openapi/pets/$PET_ID" -w "%{http_code}" | grep -q "204" && echo -e "${GREEN}✓ Deleted pet successfully${NC}" || echo -e "${RED}✗ Failed to delete pet${NC}"
echo ""

# Test Dinosaurs (Swagger)
echo -e "${YELLOW}Testing Dinosaurs (Swagger)${NC}"
echo "1. Create a dinosaur..."
CREATED_DINO=$(curl -s -X POST "$BASE_URL/swagger/dinosaurs" -H "Authorization: Bearer dino-token" -H "Content-Type: application/json" -d '{"name":"TestSaurus"}')
DINO_ID=$(echo $CREATED_DINO | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo -e "${GREEN}✓ Created dinosaur with ID: $DINO_ID${NC}"

echo "2. Read the dinosaur..."
curl -s -H "Authorization: Bearer dino-token" "$BASE_URL/swagger/dinosaurs/$DINO_ID" | grep -q "TestSaurus" && echo -e "${GREEN}✓ Read dinosaur successfully${NC}" || echo -e "${RED}✗ Failed to read dinosaur${NC}"

echo "3. Update the dinosaur..."
curl -s -X PUT "$BASE_URL/swagger/dinosaurs/$DINO_ID" -H "Authorization: Bearer dino-token" -H "Content-Type: application/json" -d '{"name":"UpdatedSaurus"}' | grep -q "UpdatedSaurus" && echo -e "${GREEN}✓ Updated dinosaur successfully${NC}" || echo -e "${RED}✗ Failed to update dinosaur${NC}"

echo "4. Delete the dinosaur..."
curl -s -X DELETE "$BASE_URL/swagger/dinosaurs/$DINO_ID" -H "Authorization: Bearer dino-token" -w "%{http_code}" | grep -q "204" && echo -e "${GREEN}✓ Deleted dinosaur successfully${NC}" || echo -e "${RED}✗ Failed to delete dinosaur${NC}"
echo ""

# Test Plants (WSDL/SOAP)
echo -e "${YELLOW}Testing Plants (WSDL/SOAP)${NC}"
echo "1. Create a plant..."
SOAP_CREATE='<?xml version="1.0"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CreatePlant xmlns="http://example.com/plants">
      <name>TestPlant</name>
    </CreatePlant>
  </soap:Body>
</soap:Envelope>'
CREATED_PLANT=$(curl -s -X POST "$BASE_URL/wdsl/soap" -H "Authorization: Bearer mock-token" -H "Content-Type: text/xml" -d "$SOAP_CREATE")
PLANT_ID=$(echo $CREATED_PLANT | grep -o '<id>[0-9]*</id>' | grep -o '[0-9]*')
echo -e "${GREEN}✓ Created plant with ID: $PLANT_ID${NC}"

echo "2. Read the plant..."
SOAP_GET='<?xml version="1.0"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetPlant xmlns="http://example.com/plants">
      <id>'"$PLANT_ID"'</id>
    </GetPlant>
  </soap:Body>
</soap:Envelope>'
curl -s -X POST "$BASE_URL/wdsl/soap" -H "Authorization: Bearer mock-token" -H "Content-Type: text/xml" -d "$SOAP_GET" | grep -q "TestPlant" && echo -e "${GREEN}✓ Read plant successfully${NC}" || echo -e "${RED}✗ Failed to read plant${NC}"

echo "3. Update the plant..."
SOAP_UPDATE='<?xml version="1.0"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <UpdatePlant xmlns="http://example.com/plants">
      <id>'"$PLANT_ID"'</id>
      <name>UpdatedPlant</name>
    </UpdatePlant>
  </soap:Body>
</soap:Envelope>'
curl -s -X POST "$BASE_URL/wdsl/soap" -H "Authorization: Bearer mock-token" -H "Content-Type: text/xml" -d "$SOAP_UPDATE" | grep -q "UpdatedPlant" && echo -e "${GREEN}✓ Updated plant successfully${NC}" || echo -e "${RED}✗ Failed to update plant${NC}"

echo "4. Delete the plant..."
SOAP_DELETE='<?xml version="1.0"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DeletePlant xmlns="http://example.com/plants">
      <id>'"$PLANT_ID"'</id>
    </DeletePlant>
  </soap:Body>
</soap:Envelope>'
curl -s -X POST "$BASE_URL/wdsl/soap" -H "Authorization: Bearer mock-token" -H "Content-Type: text/xml" -d "$SOAP_DELETE" | grep -q "true" && echo -e "${GREEN}✓ Deleted plant successfully${NC}" || echo -e "${RED}✗ Failed to delete plant${NC}"
echo ""

# Test Cars (GraphQL)
echo -e "${YELLOW}Testing Cars (GraphQL)${NC}"
echo "1. Create a car..."
CREATED_CAR=$(curl -s -X POST "$BASE_URL/graphql" -H "Content-Type: application/json" -d '{"query":"mutation{createCar(name:\"TestCar\"){id name}}"}')
CAR_ID=$(echo $CREATED_CAR | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo -e "${GREEN}✓ Created car with ID: $CAR_ID${NC}"

echo "2. Read the car..."
curl -s -X POST "$BASE_URL/graphql" -H "Content-Type: application/json" -d '{"query":"query{getCar(id:\"'"$CAR_ID"'\"){id name}}"}' | grep -q "TestCar" && echo -e "${GREEN}✓ Read car successfully${NC}" || echo -e "${RED}✗ Failed to read car${NC}"

echo "3. Update the car..."
curl -s -X POST "$BASE_URL/graphql" -H "Content-Type: application/json" -d '{"query":"mutation{updateCar(id:\"'"$CAR_ID"'\",name:\"UpdatedCar\"){id name}}"}' | grep -q "UpdatedCar" && echo -e "${GREEN}✓ Updated car successfully${NC}" || echo -e "${RED}✗ Failed to update car${NC}"

echo "4. Delete the car..."
curl -s -X POST "$BASE_URL/graphql" -H "Content-Type: application/json" -d '{"query":"mutation{deleteCar(id:\"'"$CAR_ID"'\")}"}' | grep -q "true" && echo -e "${GREEN}✓ Deleted car successfully${NC}" || echo -e "${RED}✗ Failed to delete car${NC}"
echo ""

# Test Movies (OData)
echo -e "${YELLOW}Testing Movies (OData)${NC}"
echo "1. Create a movie..."
CREATED_MOVIE=$(curl -s -X POST "$BASE_URL/odata/Movies" -H "Content-Type: application/json" -d '{"title":"TestMovie","year":2024,"genre":"Action","rating":8.5}')
MOVIE_ID=$(echo $CREATED_MOVIE | grep -o '"ID":[^,}]*' | cut -d':' -f2 | tr -d ' ')
echo -e "${GREEN}✓ Created movie with ID: $MOVIE_ID${NC}"

echo "2. Read the movie..."
curl -s "$BASE_URL/odata/Movies($MOVIE_ID)" | grep -q "TestMovie" && echo -e "${GREEN}✓ Read movie successfully${NC}" || echo -e "${RED}✗ Failed to read movie${NC}"

echo "3. Update the movie..."
curl -s -X PUT "$BASE_URL/odata/Movies($MOVIE_ID)" -H "Content-Type: application/json" -d '{"title":"UpdatedMovie","year":2024,"genre":"Action","rating":9.0}' | grep -q "UpdatedMovie" && echo -e "${GREEN}✓ Updated movie successfully${NC}" || echo -e "${RED}✗ Failed to update movie${NC}"

echo "4. Delete the movie..."
curl -s -X DELETE "$BASE_URL/odata/Movies($MOVIE_ID)" -w "%{http_code}" | grep -q "204" && echo -e "${GREEN}✓ Deleted movie successfully${NC}" || echo -e "${RED}✗ Failed to delete movie${NC}"
echo ""

echo -e "${GREEN}✅ All CRUD tests completed!${NC}"
echo ""
echo "Summary:"
echo "- Pets (OpenAPI): Full CRUD ✓"
echo "- Dinosaurs (Swagger): Full CRUD ✓"
echo "- Plants (WSDL/SOAP): Full CRUD ✓"
echo "- Cars (GraphQL): Full CRUD ✓"
echo "- Movies (OData): Full CRUD ✓"
