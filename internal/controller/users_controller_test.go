package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/aachex/service/internal/model"
	"github.com/aachex/service/internal/pagination"
	"github.com/aachex/service/internal/repository/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres driver
)

func TestGetUsers(t *testing.T) {
	loadEnv(t)
	db := openDb(t)
	users := postgres.NewUsersRepository(db)

	r, err := http.NewRequest(http.MethodPost, "/api/v1/users/get?offset=0&limit=10", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	c := NewUsersController(users, nil)
	pagination.Middleware(c.GetUsers)(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Wanted status code 201, got %d", w.Result().StatusCode)
	}
}

func TestCreateAndDeleteUser(t *testing.T) {
	loadEnv(t)
	db := openDb(t)
	users := postgres.NewUsersRepository(db)

	r, err := http.NewRequest(http.MethodPost, "/api/v1/users/new", bytes.NewReader([]byte(`{"name": "test", "surname": "testsurname"}`)))
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	c := NewUsersController(users, nil)
	c.CreateUser(w, r)

	if w.Result().StatusCode != http.StatusCreated {
		t.Errorf("Wanted status code 201, got %d", w.Result().StatusCode)
	}

	// delete created user

	b, err := io.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
	}

	var createdUser model.User
	err = json.Unmarshal(b, &createdUser)
	if err != nil {
		t.Error(err)
	}

	r, err = http.NewRequest(http.MethodDelete, "/api/v1/users/delete/{id}", nil)
	if err != nil {
		t.Error(err)
	}
	r.SetPathValue("id", strconv.FormatInt(createdUser.Id, 10))

	w = httptest.NewRecorder()
	c.DeleteUser(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Wanted status code 200, got %d", w.Result().StatusCode)
	}
}

// helpers

func openDb(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DB_CONN"))
	if err != nil {
		t.Error(err)
	}
	return db
}

func loadEnv(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Error("Failed to load .env file")
	}
}
