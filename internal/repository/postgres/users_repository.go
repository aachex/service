package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/aachex/service/internal/model"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

// createFilteringQuery генерирует SQL-запрос, который фильтрует и возвращает данные в соответствии с фильтром filter.
func createFilteringQuery(offset, limit int, filter map[string][]any) (query string, params []any) {
	// запрос по умолчанию, который вернёт выборку пользователей
	query = `
		SELECT *
		FROM (SELECT id, name, surname, patronymic, age, gender, nationality FROM users OFFSET $1 LIMIT $2) 
		WHERE true
	`

	params = []any{offset, limit}

	// Начинаем с третьего параметра, потому что параметры 1 и 2 - offset и limit
	pholder := 3
	for field, targets := range filter {
		if field == "" || len(targets) == 0 {
			continue
		}

		// field - имя поля в базе данных
		// targets - желаемое значение для k
		// pholder - номер плейсхолдера ($1, $2 и т. д.)
		query += " AND ("
		for _, t := range targets {
			query += fmt.Sprintf(" %s = $%d OR", field, pholder)
			params = append(params, t)
			pholder++
		}

		query = strings.TrimSuffix(query, " OR") // убираем последний OR
		query += ")"
	}

	query += " ORDER BY id"

	return query, params
}

// GetFiltered возвращает список пользователей, поля которых равны указанным значениям.
// Параметр filter является мапой, в которой ключи - имена свойств, а значения - желаемые значения для свойств.
func (r *UsersRepository) GetFiltered(ctx context.Context, filter map[string][]any, offset, limit int) ([]model.User, error) {
	query, params := createFilteringQuery(offset, limit, filter)
	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0)
	var u model.User
	for rows.Next() {
		err = rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Patronymic, &u.Age, &u.Gender, &u.Nationality)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UsersRepository) GetById(ctx context.Context, id int64, offset, limit int) (user model.User, err error) {
	filter := map[string][]any{
		"id": {id},
	}
	match, err := r.GetFiltered(ctx, filter, offset, limit)
	if err != nil {
		return user, err
	}

	if len(match) == 0 {
		return user, nil
	}

	return match[0], nil
}

// Create создаёт нового пользователя в базе данных.
func (r *UsersRepository) Create(ctx context.Context, name, surname, patronymic string, age int, gender, nationality string) (int64, error) {
	row := r.db.QueryRowContext(
		ctx,
		`INSERT INTO users(name, surname, patronymic, age, gender, nationality) 
		VALUES($1, $2, $3, $4, $5, $6) RETURNING id`, name, surname, patronymic, age, gender, nationality)

	var uid int64
	if err := row.Scan(&uid); err != nil {
		return -1, err
	}

	return uid, nil
}

func (r *UsersRepository) Update(ctx context.Context, id int64, updates map[string]any) error {
	if len(updates) == 0 {
		return errors.New("no updates")
	}
	if _, ok := updates["id"]; ok {
		return errors.New("field id is not updatable")
	}

	// строим SQL-запрос, который обновит все поля, указанные в updates
	params := []any{}
	updQuery := "UPDATE USERS SET"
	pholder := 1
	for field, val := range updates {
		updQuery += fmt.Sprintf(" %s = $%d", field, pholder)
		if pholder < len(updates) {
			updQuery += ","
		}

		params = append(params, val)
		pholder++
	}

	updQuery += fmt.Sprintf(" WHERE id = $%d", pholder)
	params = append(params, id)

	_, err := r.db.ExecContext(ctx, updQuery, params...)
	if err != nil {
		return err
	}

	return nil
}

// Delete удаляет пользователя из базы данных по id.
func (r *UsersRepository) Delete(ctx context.Context, uid int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", uid)
	return err
}

// Exists воззвращает true, если пользователь с указанным id существует, иначе false.
func (r *UsersRepository) Exists(ctx context.Context, id int64) bool {
	row := r.db.QueryRowContext(ctx, "SELECT id FROM users WHERE id = $1", id)
	return row.Scan() != sql.ErrNoRows
}
