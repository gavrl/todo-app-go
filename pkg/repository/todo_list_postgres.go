package repository

import (
	"fmt"
	"github.com/gavrl/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, nil
	}

	createUserListsQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", userListsTable)
	_, err = tx.Exec(createUserListsQuery, userId, id)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, nil
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id where ul.user_id = $1",
		todoListTable, userListsTable)
	err := r.db.Select(&lists, query, userId)

	return lists, err
}

func (r *TodoListPostgres) GetById(userId int, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul "+
		"on tl.id = ul.list_id where ul.user_id = $1 AND ul.list_id = $2",
		todoListTable, userListsTable)
	err := r.db.Get(&list, query, userId, listId)

	return list, err
}

func (r *TodoListPostgres) Delete(userId int, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul where tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2",
		todoListTable, userListsTable)

	_, err := r.db.Exec(query, userId, listId)

	return err
}

func (r *TodoListPostgres) Update(userId int, listId int, input todo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id = $%d AND ul.user_id = $%d",
		todoListTable, setQuery, userListsTable, argId, argId+1)

	args = append(args, listId, userId)

	logrus.Debugf("updated query: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
