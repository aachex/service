package postgres

import (
	"database/sql"
	"os"
	"testing"

	"github.com/aachex/service/internal/model"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres driver
)

// mock struct
type m struct {
	name    string
	surname string
}

var mock = m{
	name:    "Artem",
	surname: "Dmitriev",
}

func TestCreate(t *testing.T) {
	loadEnv(t)
	db := openDb(t)

	repo := NewUsersRepository(db)

	user, err := repo.Create(t.Context(), mock.name, mock.surname, "")
	if err != nil {
		t.Error(err)
	}

	if !repo.Exists(t.Context(), user.Id) {
		t.Errorf("user %d wasn't created", user.Id)
	}
}

func TestGetFiltered(t *testing.T) {
	loadEnv(t)
	db := openDb(t)

	filter := make(map[string]string)
	filter["name"] = mock.name
	filter["surname"] = mock.surname

	repo := NewUsersRepository(db)
	users := []model.User{
		{Name: mock.name, Surname: mock.surname, Patronymic: ""},
		{Name: "name", Surname: "surname", Patronymic: ""},
		{Name: "super name", Surname: "noname", Patronymic: "ha-ha-ha"},
		{Name: mock.name, Surname: mock.surname, Patronymic: "дайте пожалуйста оффер!!!"},
	}

	for i := range users {
		created, err := repo.Create(t.Context(), users[i].Name, users[i].Surname, users[i].Patronymic)
		if err != nil {
			t.Error(err)
		}
		users[i].Id = created.Id
	}

	filtered, err := repo.GetFiltered(t.Context(), filter)
	if err != nil {
		t.Error(err)
	}

	for _, u := range filtered {
		if u.Name != mock.name || u.Surname != mock.surname {
			t.Error("incorrect work of GetFiltered: name or surname mismatch")
		}
	}

	// clear db
	for _, u := range users {
		err := repo.Delete(t.Context(), u.Id)
		if err != nil {
			t.Error(err)
		}
	}
}

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
