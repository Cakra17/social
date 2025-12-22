package handlers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cakra17/social/internal/handlers"
	"github.com/cakra17/social/internal/models"
	"github.com/cakra17/social/internal/store"
	"github.com/cakra17/social/internal/utils"
	"github.com/cakra17/social/pkg/jwt"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

var (
	db = store.ConnectDB(store.DBConfig{
		DB_USERNAME: "admin",
		DB_PASSWORD: "adminsecret",
		DB_HOST: "localhost",
		DB_PORT: "5432",
		DB_NAME: "social",
		DB_MaxOpenConn: 30,
		DB_MaxIdleConn: 30,
		DB_MaxConnLifetime: 15 * time.Minute,
		DB_MaxConnIdletime: 15 * time.Minute,
	})

	jwtAuthenticator = jwt.NewJWTAuthenticator("mysecret", 5 * time.Hour)
	logger = utils.NewLogger()

	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	userRepo = store.NewUserRepo(db, logger)

	userHandler = handlers.NewUserHandler(handlers.UserHandlerConfig{
		UserRepo: userRepo,
		JWTAuthenticator: jwtAuthenticator,
		Redis: rdb,
	})
)

func TestCreateUserBadPayload(t *testing.T) {
	req := httptest.NewRequest("POST", "/users", nil)
	rr := httptest.NewRecorder()

	userHandler.CreateUser(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler return wromg status code: got %v want %v", status, http.StatusBadRequest)
	}

	var payload models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &payload)
	if err != nil {
		log.Fatalf("Failed to Unmarshaling json: %v", err)
	}

	if payload.Message != "Payload Malformed" {
		t.Errorf("handler return wromg message: got %s want %v", payload.Message, "Payload Malformed")
	}
}

func TestCreateInvalidPayload(t *testing.T) {
	userPayload := models.RegisterPayload{
		Username: "dwda",
		Email: "dasadsawdaw",
		Password: "wdwad",
	}
	jsonBytes,_ := json.Marshal(&userPayload)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBytes))
	rr := httptest.NewRecorder()

	userHandler.CreateUser(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler return wromg status code: got %v want %v", status, http.StatusBadRequest)
	}
	
	var payload models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &payload)
	if err != nil {
		log.Fatalf("Failed to Unmarshaling json: %v", err)
	}

	if payload.Message != "Invalid Payload" {
		t.Errorf("handler return wromg message: got %s want %v", payload.Message,"Invalid Payload")
	}
}

func TestCreateDuplicateUser(t *testing.T) {
	userPayload := models.RegisterPayload{
		Username: "nightfall",
		Email: "nightfall@gmail.com",
		Password: "nightfallgantenk",
	}
	jsonBytes,_ := json.Marshal(&userPayload)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBytes))
	rr := httptest.NewRecorder()

	userHandler.CreateUser(rr, req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler return wromg status code: got %v want %v", status, http.StatusBadRequest)
	}
	
	var payload models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &payload)
	if err != nil {
		log.Fatalf("Failed to Unmarshaling json: %v", err)
	}

	if payload.Message != "Credentials already used" {
		t.Errorf("handler return wromg message: got %s want %v", payload.Message,"Credentials already used")
	}
}

func TestCreateValidUser(t *testing.T) {
	userPayload := models.RegisterPayload{
		Username: "nightfall",
		Email: "nightfall123@gmail.com",
		Password: "nightfallgantenk",
	}

	jsonBytes,_ := json.Marshal(&userPayload)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBytes))
	rr := httptest.NewRecorder()

	userHandler.CreateUser(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler return wromg status code: got %v want %v", status, http.StatusCreated)
	}

	var payload models.Response
	err := json.Unmarshal(rr.Body.Bytes(), &payload)
	if err != nil {
		log.Fatalf("Failed to Unmarshaling json: %v", err)
	}

	if !payload.Success {
		t.Errorf("Failed to Create user")
	}

	if payload.Data == nil {
		t.Errorf("handler didn't return correct data: got %v expect %v", payload.Data, userPayload)
	}
}

func TestLoginWithInvalidPayload(t *testing.T) {
	req := httptest.NewRequest("POST", "/login", nil)
	rr := httptest.NewRecorder()

	userHandler.Authenticate(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler return wromg status code: got %v want %v", status, http.StatusBadRequest)
	}

	var payload models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &payload)
	if err != nil {
		log.Fatalf("Failed to Unmarshaling json: %v", err)
	}

	if payload.Message != "Payload Malformed" {
		t.Errorf("handler return wromg message: got %s want %v", payload.Message, "Payload Malformed")
	}
}

func TestLoginWithBadPayload(t *testing.T) {
	userPayload := models.LoginPayload{
		Email: "dasadsawdaw",
		Password: "wdwad",
	}
	jsonBytes,_ := json.Marshal(&userPayload)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBytes))
	rr := httptest.NewRecorder()

	userHandler.Authenticate(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler return wromg status code: got %v want %v", status, http.StatusBadRequest)
	}
	
	var payload models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &payload)
	if err != nil {
		log.Fatalf("Failed to Unmarshaling json: %v", err)
	}

	if payload.Message != "Invalid Payload" {
		t.Errorf("handler return wromg message: got %s want %v", payload.Message,"Invalid Payload")
	}
}

func TestLoginWithFalseCredentials(t *testing.T) {
	userPayload := models.LoginPayload{
		Email: "nightfalloff@gmail.com",
		Password: "nightfallgantenk",
	}
	jsonBytes,_ := json.Marshal(&userPayload)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBytes))
	rr := httptest.NewRecorder()

	userHandler.Authenticate(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler return wromg status code: got %v want %v", status, http.StatusNotFound)
	}

	var payload models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &payload)
	if err != nil {
		log.Fatalf("Failed to Unmarshaling json: %v", err)
	}

	if payload.Message != "User not found" {
		t.Errorf("handler return wromg message: got %s want %v", payload.Message, "User not found")
	}
}

func TestLoginWithWrongPassword(t *testing.T) {
	userPayload := models.LoginPayload{
		Email: "nightfall@gmail.com",
		Password: "nightfallgantenk123",
	}
	jsonBytes,_ := json.Marshal(&userPayload)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBytes))
	rr := httptest.NewRecorder()

	userHandler.Authenticate(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler return wromg status code: got %v want %v", status, http.StatusBadRequest)
	}
	
	var payload models.ErrorResponse
	err := json.Unmarshal(rr.Body.Bytes(), &payload)
	if err != nil {
		log.Fatalf("Failed to Unmarshaling json: %v", err)
	}

	if payload.Message != "Wrong password" {
		t.Errorf("handler return wromg message: got %s want %v", payload.Message, "Wrong password")
	}
}

func TestLoginWithValidCredentials(t *testing.T) {
	userPayload := models.LoginPayload{
		Email: "nightfall@gmail.com",
		Password: "nightfallgantenk",
	}
	jsonBytes,_ := json.Marshal(&userPayload)
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBytes))
	rr := httptest.NewRecorder()

	userHandler.Authenticate(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler return wromg status code: got %v want %v", status, http.StatusOK)
	}
	
	var payload models.Response
	err := json.Unmarshal(rr.Body.Bytes(), &payload)
	if err != nil {
		log.Fatalf("Failed to Unmarshaling json: %v", err)
	}

	if !payload.Success {
		t.Errorf("Failed to login")
	}

	if payload.Message != "success to login" {
		t.Errorf("handler return wromg message: got %s want %v", payload.Message, "success to login")
	}
}