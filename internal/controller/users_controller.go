package controller

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aachex/service/internal/enricher"
	"github.com/aachex/service/internal/logging"
	"github.com/aachex/service/internal/model"
	"github.com/aachex/service/internal/pagination"
)

type usersRepository interface {
	GetFiltered(ctx context.Context, filter map[string][]any, offset, limit int) ([]model.User, error)
	Create(ctx context.Context, name, surname, patronymic string, age int, gender, nationality string) (int64, error)
	Update(ctx context.Context, id int64, updates map[string]any) error
	Delete(ctx context.Context, uid int64) error
}

type UsersController struct {
	users  usersRepository
	logger *slog.Logger
}

func NewUsersController(ur usersRepository, l *slog.Logger) *UsersController {
	return &UsersController{
		users:  ur,
		logger: l,
	}
}

func (c *UsersController) RegisterHandlers(mux *http.ServeMux) {
	prefix := "/api/v1"

	mux.HandleFunc(
		"POST "+prefix+"/users/new",
		logging.Middleware(c.logger, c.CreateUser))

	mux.HandleFunc(
		"POST "+prefix+"/users/get",
		logging.Middleware(c.logger, pagination.Middleware(c.GetUsers)))

	mux.HandleFunc(
		"PATCH "+prefix+"/users/upd/{id}",
		logging.Middleware(c.logger, c.UpdateUser))

	mux.HandleFunc(
		"DELETE "+prefix+"/users/delete/{id}",
		logging.Middleware(c.logger, c.DeleteUser))
}

//	@summary	Получение пользователей с возможностью фильтрации по полям.
//	@produce	json
//	@success	200
//	@param		offset	query	integer				true	"offset"
//	@param		limit	query	integer				true	"limit"
//	@param		request	body	map[string][]any	true	"filter"
//	@router		/users/get [post]
func (c *UsersController) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Пагинация
	pag := r.Context().Value(pagination.CtxKey("pagination")).(pagination.Pagination)

	filter, err := readBody[map[string][]any](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, err := c.users.GetFiltered(r.Context(), filter, pag.Offset, pag.Limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeReponse(users, w)
}

type reqBody struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

//	@summary	Создание нового пользователя в базе данных.
//	@accept		json
//	@produce	json
//	@param		request	body		reqBody	true	"Request"
//	@success	200		{object}	model.User
//	@router		/users/new [post]
func (c *UsersController) CreateUser(w http.ResponseWriter, r *http.Request) {
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

//	@summary	Обновляет указанные данные у пользователя по id.
//	@accept		json
//	@success	200
//	@param		id		path	integer		true	"User ID"
//	@param		request	body	model.User	true	"Request"
//	@router		/users/upd/{id} [patch]
func (c *UsersController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updates, err := readBody[map[string]any](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.users.Update(r.Context(), id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

//	@summary	Удаление пользователя по id.
//	@success	200
//	@param		id	path	integer	true	"User ID"
//	@router		/users/delete/{id} [delete]
func (c *UsersController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
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
