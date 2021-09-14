package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"golang.org/x/xerrors"
)

// リファクタリング前
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/todos", CreateTodo).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: router,
	}

	errorCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			errorCh <- err
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errorCh:
		panic(err)
	case s := <-signalCh:
		log.Printf("SIGNAL %s received", s.String())
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			panic(err)
		}
	}
}

func NewMySQLConnection() (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_db?parseTime=true")
	if err != nil {
		return nil, err
	}

	return db, nil
}

type Todo struct {
	ID                 int       `json:"id"                 db:"id"`
	Title              string    `json:"title"              db:"title"`
	ImplementationDate time.Time `json:"implementationDate" db:"implementation_date"`
	DueDate            time.Time `json:"dueDate"            db:"due_date"`
	StatusID           int       `json:"statusID"           db:"status_id"`
	PriorityID         *int      `json:"priorityID"         db:"priority_id"`
	Memo               *string   `json:"memo"               db:"memo"`
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		panic(err)
	}

	db, err := NewMySQLConnection()
	if err != nil {
		panic(err)
	}

	createQuery := `
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

	result, err := db.Exec(
		createQuery,
		todo.Title,
		todo.ImplementationDate,
		todo.DueDate,
		todo.StatusID,
		todo.PriorityID,
		todo.Memo,
	)
	if err != nil {
		panic(err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

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

	var outTodo Todo
	if err = db.QueryRowx(fetchQuery, lastInsertID).StructScan(&outTodo); err != nil {
		if xerrors.Is(err, sql.ErrNoRows) {
			panic(err)
		}

		panic(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&outTodo); err != nil {
		panic(err)
	}
}
