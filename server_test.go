package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test OpenAPI Pets endpoints
func TestOpenAPI_ListPets(t *testing.T) {
	srv, _ := setupTestServer(t)
	defer srv.store.db.Close()

	req := httptest.NewRequest("GET", "/openapi/pets?limit=5", nil)
	w := httptest.NewRecorder()
	srv.handlePets(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var pets []EnhancedPet
	if err := json.Unmarshal(w.Body.Bytes(), &pets); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(pets) == 0 {
		t.Fatal("expected at least one pet")
	}

	// Verify nested structure
	if pets[0].Owner.Name == "" {
		t.Error("expected owner name to be populated")
	}
	if len(pets[0].Medical) == 0 {
		t.Error("expected medical records")
	}
	if len(pets[0].Tags) == 0 {
		t.Error("expected tags")
	}
}

func TestOpenAPI_CreatePet(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	body := `{"name":"Buddy (Beagle)"}`
	req := httptest.NewRequest("POST", "/openapi/pets", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.handlePets(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var pet EnhancedPet
	if err := json.Unmarshal(w.Body.Bytes(), &pet); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if pet.ID == "" {
		t.Error("expected pet ID")
	}
	if pet.Name != "Buddy" {
		t.Errorf("expected name 'Buddy', got %q", pet.Name)
	}
	if pet.Breed != "Beagle" {
		t.Errorf("expected breed 'Beagle', got %q", pet.Breed)
	}
}

func TestOpenAPI_GetPet(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	req := httptest.NewRequest("GET", "/openapi/pets/1", nil)
	w := httptest.NewRecorder()
	srv.handlePet(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var pet EnhancedPet
	if err := json.Unmarshal(w.Body.Bytes(), &pet); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if pet.ID != "1" {
		t.Errorf("expected ID '1', got %q", pet.ID)
	}
}

func TestOpenAPI_UpdatePet(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	body := `{"name":"Max Updated (Golden Retriever)"}`
	req := httptest.NewRequest("PUT", "/openapi/pets/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.handlePet(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var pet EnhancedPet
	if err := json.Unmarshal(w.Body.Bytes(), &pet); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if !strings.Contains(pet.Name, "Updated") {
		t.Errorf("expected updated name, got %q", pet.Name)
	}
}

func TestOpenAPI_DeletePet(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	// Create a pet first
	body := `{"name":"ToDelete (Test)"}`
	req := httptest.NewRequest("POST", "/openapi/pets", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	srv.handlePets(w, req)

	var created EnhancedPet
	json.Unmarshal(w.Body.Bytes(), &created)

	// Delete it
	req = httptest.NewRequest("DELETE", "/openapi/pets/"+created.ID, nil)
	w = httptest.NewRecorder()
	srv.handlePet(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	// Verify it's gone
	req = httptest.NewRequest("GET", "/openapi/pets/"+created.ID, nil)
	w = httptest.NewRecorder()
	srv.handlePet(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", w.Code)
	}
}

// Test Swagger Dinosaurs endpoints (with auth)
func TestSwagger_ListDinosaurs(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	req := httptest.NewRequest("GET", "/swagger/dinosaurs?limit=5", nil)
	req.Header.Set("Authorization", "Bearer dino-token")
	w := httptest.NewRecorder()
	srv.handleDinosaurs(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var dinos []EnhancedDinosaur
	if err := json.Unmarshal(w.Body.Bytes(), &dinos); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(dinos) == 0 {
		t.Fatal("expected at least one dinosaur")
	}

	// Verify nested structure
	if dinos[0].Species == "" {
		t.Error("expected species")
	}
	if dinos[0].Discovered.Year == 0 {
		t.Error("expected discovery year")
	}
	if len(dinos[0].Features) == 0 {
		t.Error("expected features")
	}
}

func TestSwagger_Unauthorized(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	req := httptest.NewRequest("GET", "/swagger/dinosaurs", nil)
	// No Authorization header
	w := httptest.NewRecorder()
	srv.handleDinosaurs(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestSwagger_CreateDinosaur(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	body := `{"name":"Spinosaurus"}`
	req := httptest.NewRequest("POST", "/swagger/dinosaurs", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer dino-token")
	w := httptest.NewRecorder()
	srv.handleDinosaurs(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var dino EnhancedDinosaur
	if err := json.Unmarshal(w.Body.Bytes(), &dino); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if dino.Name != "Spinosaurus" {
		t.Errorf("expected name 'Spinosaurus', got %q", dino.Name)
	}
}

// Test GraphQL Cars endpoints
func TestGraphQL_ListCars(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	query := `{"query":"query{listCars(limit:5){id name}}"}`
	req := httptest.NewRequest("POST", "/graphql", strings.NewReader(query))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("graphql-user", "graphql-pass")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			ListCars []map[string]interface{} `json:"listCars"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(resp.Data.ListCars) == 0 {
		t.Fatal("expected at least one car")
	}
}

func TestGraphQL_CreateCar(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	query := `{"query":"mutation{createCar(name:\"Porsche 911\"){id name}}"}`
	req := httptest.NewRequest("POST", "/graphql", strings.NewReader(query))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("graphql-user", "graphql-pass")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Data struct {
			CreateCar struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"createCar"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Data.CreateCar.Name != "Porsche 911" {
		t.Errorf("expected name 'Porsche 911', got %q", resp.Data.CreateCar.Name)
	}
}

// Test OData Movies endpoints
func TestOData_ListMovies(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	req := httptest.NewRequest("GET", "/odata/Movies", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Value []Movie `json:"value"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(resp.Value) == 0 {
		t.Fatal("expected at least one movie")
	}
}

func TestOData_FilterMovies(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	// Filter movies with Year > 2010
	req := httptest.NewRequest("GET", "/odata/Movies?$filter=Year%20gt%202010", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Value []Movie `json:"value"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// All filtered movies should have Year > 2010
	for _, m := range resp.Value {
		if m.Year <= 2010 {
			t.Errorf("movie %q has year %d, expected > 2010", m.Title, m.Year)
		}
	}
}

func TestOData_OrderByMovies(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	req := httptest.NewRequest("GET", "/odata/Movies?$orderby=Rating%20desc&$top=3", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp struct {
		Value []Movie `json:"value"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(resp.Value) == 0 {
		t.Fatal("expected at least one movie")
	}

	// Verify descending order
	for i := 1; i < len(resp.Value); i++ {
		if resp.Value[i].Rating > resp.Value[i-1].Rating {
			t.Errorf("movies not in descending order: %f > %f", resp.Value[i].Rating, resp.Value[i-1].Rating)
		}
	}
}

func TestOData_CreateMovie(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	body := `{"Title":"Test Movie","Year":2024,"Genre":"Action","Rating":7.5,"Director":"Test Director"}`
	req := httptest.NewRequest("POST", "/odata/Movies", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var movie Movie
	if err := json.Unmarshal(w.Body.Bytes(), &movie); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if movie.Title != "Test Movie" {
		t.Errorf("expected title 'Test Movie', got %q", movie.Title)
	}
}

// Test SOAP Plants endpoints
func TestSOAP_ListPlants(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	soapEnv := `<?xml version="1.0"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ListPlants xmlns="http://example.com/plants">
      <Limit>5</Limit>
    </ListPlants>
  </soap:Body>
</soap:Envelope>`

	req := httptest.NewRequest("POST", "/wdsl", strings.NewReader(soapEnv))
	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("Authorization", "Bearer mock-token")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Basic XML validation
	var envelope struct {
		XMLName xml.Name `xml:"Envelope"`
	}
	if err := xml.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("failed to parse SOAP response: %v", err)
	}
}

func TestSOAP_CreatePlant(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	soapEnv := `<?xml version="1.0"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CreatePlant xmlns="http://example.com/plants">
      <name>Bamboo</name>
    </CreatePlant>
  </soap:Body>
</soap:Envelope>`

	req := httptest.NewRequest("POST", "/wdsl", strings.NewReader(soapEnv))
	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("Authorization", "Bearer mock-token")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "Bamboo") {
		t.Error("expected plant name in response")
	}
}

// Test JSON-RPC Calculator endpoints
func TestJSONRPC_Add(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	body := `{"jsonrpc":"2.0","method":"add","params":{"a":5,"b":3},"id":1}`
	req := httptest.NewRequest("POST", "/jsonrpc", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Jsonrpc string  `json:"jsonrpc"`
		Result  float64 `json:"result"`
		ID      int     `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Result != 8 {
		t.Errorf("expected result 8, got %f", resp.Result)
	}
}

func TestJSONRPC_Multiply(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	body := `{"jsonrpc":"2.0","method":"multiply","params":{"a":4,"b":7},"id":2}`
	req := httptest.NewRequest("POST", "/jsonrpc", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Jsonrpc string  `json:"jsonrpc"`
		Result  float64 `json:"result"`
		ID      int     `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Result != 28 {
		t.Errorf("expected result 28, got %f", resp.Result)
	}
}

func TestJSONRPC_Divide(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	body := `{"jsonrpc":"2.0","method":"divide","params":{"a":10,"b":2},"id":3}`
	req := httptest.NewRequest("POST", "/jsonrpc", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Jsonrpc string  `json:"jsonrpc"`
		Result  float64 `json:"result"`
		ID      int     `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Result != 5 {
		t.Errorf("expected result 5, got %f", resp.Result)
	}
}

func TestJSONRPC_DivideByZero(t *testing.T) {
	srv, mux := setupTestServer(t)
	_ = mux
	defer srv.store.db.Close()

	body := `{"jsonrpc":"2.0","method":"divide","params":{"a":10,"b":0},"id":4}`
	req := httptest.NewRequest("POST", "/jsonrpc", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Jsonrpc string `json:"jsonrpc"`
		Error   *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
		ID int `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error for division by zero")
	}
}

// Helper function to set up a test server
func setupTestServer(t *testing.T) (*Server, *http.ServeMux) {
	store, err := NewStore()
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	srv := &Server{store: store}
	
	// Set up routes like in main()
	mux := http.NewServeMux()
	mux.HandleFunc("/openapi/openapi.json", srv.handleOpenAPI)
	mux.HandleFunc("/swagger/swagger.json", srv.handleSwagger2)
	mux.HandleFunc("/wdsl/wsdl", srv.handleWSDL)
	mux.HandleFunc("/graphql", srv.handleGraphQL)
	mux.HandleFunc("/graphql/schema", srv.handleGraphQLSchema)
	mux.HandleFunc("/openapi/pets", srv.handlePets)
	mux.HandleFunc("/openapi/pets/", srv.handlePet)
	mux.HandleFunc("/swagger/dinosaurs", srv.handleDinosaurs)
	mux.HandleFunc("/swagger/dinosaurs/", srv.handleDinosaur)
	mux.HandleFunc("/wdsl/soap", srv.handleSOAP)
	mux.HandleFunc("/wdsl", srv.handleSOAP) // SOAP endpoint
	mux.HandleFunc("/odata/", srv.handleOData)
	mux.HandleFunc("/jsonrpc", srv.handleJSONRPC)
	mux.HandleFunc("/jsonrpc/openrpc.json", srv.handleOpenRPCSpec)
	
	return srv, mux
}
