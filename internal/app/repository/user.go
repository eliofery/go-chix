package repository

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/eliofery/go-chix/internal/app/model"
	"github.com/eliofery/go-chix/pkg/log"
	"log/slog"
)

// UserQuery интерфейс для запросов связанных с пользователями
type UserQuery interface {
	GetUsers() (*[]model.User, error) // GetUsers Получение всех пользователей
}

type userQuery struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

// GetUsers Получение всех пользователей
func (q *userQuery) GetUsers() (*[]model.User, error) {
	qb := q.builder.Select("id", "name", "age").
		From(model.UserTableName)

	var users []model.User
	rows, err := qb.Query()
	for rows.Next() {
		var user model.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			log.Error("Не удалось получить пользователей", slog.String("err", err.Error()))
			return nil, fmt.Errorf("не удалось получить пользователей")
		}
		users = append(users, user)
	}

	return &users, nil
}