package persistence

import (
	"database/sql"

	"golang.org/x/xerrors"

	"github.com/kazumakawahara/todo-sample/apperrors"
	"github.com/kazumakawahara/todo-sample/domain/tododomain"
	"github.com/kazumakawahara/todo-sample/infrastructure/datasource"
	"github.com/kazumakawahara/todo-sample/infrastructure/rdb"
)

type todoRepository struct {
	*rdb.MySQLHandler
}

func NewTodoRepository(mysqlHandler *rdb.MySQLHandler) *todoRepository {
	return &todoRepository{mysqlHandler}
}

func (r *todoRepository) CreateTodo(todo *tododomain.Todo) (tododomain.ID, error) {
	query := `
        INSERT INTO todos
        (
          title,
          implementation_date,
          due_date,
          status_id,
          priority_id,
          memo
        )
        VALUES
          (?, ?, ?, ?, ?, ?)`

	result, err := r.Conn.Exec(
		query,
		todo.Title().Value(),
		todo.ImplementationDate().Value(),
		todo.DueDate().Value(),
		todo.Status().Value(),
		todo.Priority().Value(),
		todo.Memo().Value(),
	)
	if err != nil {
		panic(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, apperrors.InternalServerError
	}

	idVo, err := tododomain.NewID(int(id))
	if err != nil {
		return 0, apperrors.InternalServerError
	}

	return idVo, nil
}

func (r *todoRepository) FetchTodoByID(id tododomain.ID) (*tododomain.Todo, error) {
	fetchQuery := `
        SELECT
          todos.id                  id,
          todos.title               title,
          todos.implementation_date implementation_date,
          todos.due_date            due_date,
          todos.status_id           status_id,
          todos.priority_id         priority_id,
          todos.memo                memo
        FROM
          todos
        INNER JOIN
          statuses
        ON
          statuses.id = todos.status_id
        INNER JOIN
          priorities
        ON
          priorities.id = todos.priority_id
        WHERE
          todos.id = ?`

	var todoDto datasource.Todo
	if err := r.Conn.QueryRowx(fetchQuery, id.Value()).StructScan(&todoDto); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.TodoNotFound
		}

		return nil, apperrors.InternalServerError
	}

	todoDm := tododomain.NewTodo(
		tododomain.ID(todoDto.ID),
		tododomain.Title(todoDto.Title),
		tododomain.ImplementationDate(todoDto.ImplementationDate),
		tododomain.DueDate(todoDto.DueDate),
		tododomain.Status(todoDto.StatusID),
		tododomain.Priority(todoDto.PriorityID),
		tododomain.Memo(todoDto.Memo),
	)

	return todoDm, nil
}

func (r *todoRepository) FetchTodos() ([]*tododomain.Todo, error) {
	fetchQuery := `
        SELECT
            todos.id                  id,
            todos.title               title,
            todos.implementation_date implementation_date,
            todos.due_date            due_date,
            todos.status_id           status_id,
            todos.priority_id         priority_id,
            todos.memo                memo
        FROM
            todos
        INNER JOIN
            statuses
        ON
            statuses.id = todos.status_id
        INNER JOIN
            priorities
        ON
            priorities.id = todos.priority_id`

	rows, err := r.Conn.Queryx(fetchQuery)
	if err != nil {
		return nil, apperrors.InternalServerError
	}

	var todosDto []datasource.Todo
	for rows.Next() {
		var todoDto datasource.Todo
		if err := rows.StructScan(&todoDto); err != nil {
			return nil, apperrors.InternalServerError
		}

		todosDto = append(todosDto, todoDto)
	}

	todoDms := make([]*tododomain.Todo, len(todosDto))
	for i, todoDto := range todosDto {
		todoDms[i] = tododomain.NewTodo(
			tododomain.ID(todoDto.ID),
			tododomain.Title(todoDto.Title),
			tododomain.ImplementationDate(todoDto.ImplementationDate),
			tododomain.DueDate(todoDto.DueDate),
			tododomain.Status(todoDto.StatusID),
			tododomain.Priority(todoDto.PriorityID),
			tododomain.Memo(todoDto.Memo),
		)
	}

	return todoDms, nil
}

func (r *todoRepository) UpdateTodo(todo *tododomain.Todo) (tododomain.ID, error) {
	updateQuery := `
        UPDATE
            todos
        SET 
            title = ?,
            implementation_date = ?,
            due_date = ?,
            status_id = ?,
            priority_id = ?,
            memo = ?
        WHERE
            id = ?`

	if _, err := r.Conn.Exec(
		updateQuery,
		todo.Title().Value(),
		todo.ImplementationDate().Value(),
		todo.DueDate().Value(),
		todo.Status().Value(),
		todo.Priority().Value(),
		todo.Memo().Value(),
		todo.ID().Value(),
	); err != nil {
		return 0, apperrors.InternalServerError
	}
	return 0, nil
}

func (r *todoRepository) DeleteTodo(id tododomain.ID) error {
	deleteQuery := `
        DELETE FROM
            todos
        WHERE
            id = ?`

	if _, err := r.Conn.Exec(
		deleteQuery,
		id.Value(),
	); err != nil {
		return apperrors.InternalServerError
	}

	return nil
}
