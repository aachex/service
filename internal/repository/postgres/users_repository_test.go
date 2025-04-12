package postgres

import (
	"database/sql"
	"os"
	"testing"

	"slices"

	"github.com/aachex/service/internal/model"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres driver
)

// mock struct
type m struct {
	name        string
	surname     string
	patronymic  string
	age         int
	gender      string
	nationality string
}

var mock = m{
	name:        "Artem",
	surname:     "Dmitriev",
	patronymic:  "Evgenievich",
	age:         17,
	gender:      "male",
	nationality: "RU",
}

func TestCreate(t *testing.T) {
	loadEnv(t)
	db := openDb(t)

	repo := NewUsersRepository(db)

	id, err := repo.Create(t.Context(), mock.name, mock.surname, mock.patronymic, mock.age, mock.gender, mock.nationality)
	if err != nil {
		t.Error(err)
	}

	// clear db
	defer func() {
		err = repo.Delete(t.Context(), id)
		if err != nil {
			t.Error(err)
		}
	}()

	if !repo.Exists(t.Context(), id) {
		t.Errorf("user %d wasn't created", id)
	}
}

func TestGetFiltered(t *testing.T) {
	loadEnv(t)
	db := openDb(t)

	filter := make(map[string][]any)
	filter["name"] = []any{mock.name, "name"}
	filter["surname"] = []any{mock.surname}

	repo := NewUsersRepository(db)
	users := []model.User{
		{Name: mock.name, Surname: mock.surname, Patronymic: ""},
		{Name: "name", Surname: "surname", Patronymic: ""},
		{Name: "super name", Surname: "noname", Patronymic: "ha-ha-ha"},
		{Name: mock.name, Surname: mock.surname, Patronymic: "дайте пожалуйста оффер!!!"},
	}

	for i := range users {
		id, err := repo.Create(t.Context(), users[i].Name, users[i].Surname, users[i].Patronymic, 23, "male or female", "RU")
		if err != nil {
			t.Error(err)
		}
		users[i].Id = id
	}

	// clear db
	defer func() {
		for _, u := range users {
			err := repo.Delete(t.Context(), u.Id)
			if err != nil {
				t.Error(err)
			}
		}
	}()

	filtered, err := repo.GetFiltered(t.Context(), filter, 0, 100)
	if err != nil {
		t.Error(err)
	}

	for _, u := range filtered {
		if !contains(filter["name"], u.Name) || !contains(filter["surname"], u.Surname) {
			t.Error("incorrect work of GetFiltered: name or surname mismatch")
		}
	}
}

func TestUpdate(t *testing.T) {
	loadEnv(t)
	db := openDb(t)

	repo := NewUsersRepository(db)
	id, err := repo.Create(t.Context(), mock.name, mock.surname, mock.patronymic, mock.age, mock.gender, mock.nationality)
	if err != nil {
		t.Error(err)
	}

	updates := map[string]any{
		"name":        "Dima",
		"age":         18,
		"nationality": "JP",
	}

	err = repo.Update(t.Context(), id, updates)
	if err != nil {
		t.Error(err)
	}

	user, err := repo.GetById(t.Context(), id, 0, 100000)
	if err != nil {
		t.Error(err)
	}

	if user.Name != updates["name"] || user.Age != updates["age"] || user.Nationality != updates["nationality"] {
		t.Error("failed to update user")
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
	err := godotenv.Load("../../../.env")
	if err != nil {
		t.Error("Failed to load .env file")
	}
}

func contains(s []any, e any) bool {
	return slices.Contains(s, e)
}
