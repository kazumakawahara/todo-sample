package rdb

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type MySQLHandler struct {
	Conn *sqlx.DB
}

func NewMySQLHandler() (*MySQLHandler, error) {
	conn, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_db?parseTime=true")
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return &MySQLHandler{
		Conn: conn,
	}, nil
}
