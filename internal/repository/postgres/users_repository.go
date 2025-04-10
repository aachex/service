package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/aachex/service/internal/model"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

// GetFiltered возвращает список пользователей, поля которых равны указанным значениям.
// Параметр filter является мапой, в которой ключи - имена свойств, а значения - желаемые значения для свойств.
func (r *UsersRepository) GetFiltered(ctx context.Context, offset, limit int, filter map[string]any) ([]model.User, error) {
	// запрос по умолчанию, который вернёт выборку пользователей
	query := `
		SELECT *
		FROM (SELECT id, name, surname, patronymic FROM users OFFSET $1 LIMIT $2) 
		WHERE true
	`

	params := []any{offset, limit}

	// Параметры 1 и 2 - offset и limit
	pholder := 3
	for k, v := range filter {
		if k != "" {
			// k - имя поля в базе данных
			// v - желаемое значение для k
			// pholder - номер плейсхолдера ($1, $2 и т. д.)
			query += fmt.Sprintf(" AND %s = $%d", k, pholder)
			params = append(params, v)
		}
		pholder++
	}

	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	users := make([]model.User, 0)
	var u model.User
	for rows.Next() {
		err = rows.Scan(&u.Id, &u.Name, &u.Surname, &u.Patronymic)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UsersRepository) Create(ctx context.Context, name, surname, patronymic string) (*model.User, error) {
	if name == "" || surname == "" {
		return nil, errors.New("name or surname cannot be empty")
	}

	row := r.db.QueryRowContext(ctx, "INSERT INTO users(name, surname, patronymic) VALUES($1, $2, $3) RETURNING id", name, surname, patronymic)

	var uid int64
	if err := row.Scan(&uid); err != nil {
		return nil, err
	}

	created := model.User{Id: uid, Name: name, Surname: surname, Patronymic: patronymic}
	return &created, nil
}

func (r *UsersRepository) Delete(ctx context.Context, uid int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", uid)
	return err
}

func (r *UsersRepository) Exists(ctx context.Context, id int64) bool {
	row := r.db.QueryRowContext(ctx, "SELECT id FROM users WHERE id = $1", id)
	return row.Scan() != sql.ErrNoRows
}
