package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aachex/service/internal/enricher"
	"github.com/aachex/service/internal/model"
)

type usersRepository interface {
	GetFiltered(ctx context.Context, offset, limit int, filter map[string][]any) ([]model.User, error)
	Create(ctx context.Context, name, surname, patronymic string, age int, gender, nationality string) (int64, error)
	Delete(ctx context.Context, uid int64) error
}

type UsersController struct {
	users usersRepository
}

func NewUsersController(ur usersRepository) *UsersController {
	return &UsersController{
		users: ur,
	}
}

func (c *UsersController) RegisterHandlers(mux *http.ServeMux) {
	prefix := "/api/v1"

	mux.HandleFunc(
		"POST "+prefix+"/users/new",
		c.CreateUser)

	mux.HandleFunc(
		"POST "+prefix+"/users/get",
		c.GetUsers)

	mux.HandleFunc(
		"DELETE "+prefix+"/users/delete",
		c.DeleteUser)
}

// GetUsers возвращает список пользователей. В качестве тела запроса принимает параметры фильтрации.
//
// Пример:
//
//	{
//		"name": ["Artem", "Dmitry"],
//		"surname": ["Filin", "Okunev"]
//	}
//
// При таком теле запроса метод вернёт всех пользователей с именами Artem или Dmitry и фамилиями Filin или Okunev.
//
// HTTP POST /users/get
func (c *UsersController) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Пагинация
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter, err := readBody[map[string][]any](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, err := c.users.GetFiltered(r.Context(), offset, limit, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeReponse(users, w)
}

func (c *UsersController) CreateUser(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Patronymic string `json:"patronymic"`
	}

	body, err := readBody[reqBody](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if body.Name == "" || body.Surname == "" {
		http.Error(w, "name or surname cannot be empty", http.StatusBadRequest)
		return
	}

	user := model.User{
		Name:       body.Name,
		Surname:    body.Surname,
		Patronymic: body.Patronymic,
	}

	enricher.EnrichUser(&user)

	id, err := c.users.Create(r.Context(), user.Name, user.Surname, user.Patronymic, user.Age, user.Gender, user.Nationality)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Id = id

	w.WriteHeader(http.StatusCreated)
	writeReponse(user, w)
}

func (c *UsersController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.users.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
