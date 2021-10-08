package main

import (
	//Goのサーバーで入ってきた各リクエストは個別のgoroutineで処理される
	//1. リクエストスコープな値の伝播 （スコープ内で収束される値）
	//2. キャンセル （タイムアウト処理などで使用する）
	"context"
	"strconv"

	//DBへのアクセスには sql.DB を利用
	"database/sql"

	//エンコードとは、データの符号化、つまり、他形式へのコードの変換を行うこと
	//エンコーダとは、エンコードを行うソフトウェア、もしくは、そのソフトウェアが組み込まれたハードウェアのことである。
	//JSONとは「JavaScript Object Notation」の略で、「JavaScriptのオブジェクトの書き方を元にしたデータ定義方法」のこと
	//encoding パッケージは，データをバイト表現あるいはテキスト表現に変換するほかのパッケージが共有するインターフェースを定義します。
	//json パッケージは RFC 7159 で定義された JSON のエンコードとデコードを実装します。
	//json形式のデータでやりとりする
	"encoding/json"

	"fmt"

	//logは、シンプルなログ作成の為の機能がまとめられたパッケージ。
	//標準エラー出力に任意のログメッセージを出力できる。
	"log"

	//net/httpパッケージはHTTPクライアントとサーバーの実装を提供する。
	//HTTPサーバー実装の概要を掴むのに重要なのはListenAndServe、Handle、HandleFuncの３つの関数とHandlerインターフェース
	//	net パッケージは，TCP/IP, UDP，ドメイン名解決，Unix ドメインソケットを含む，ネットワーク I/O の可搬的インターフェースを提供します
	"net/http"

	//os パッケージは，オペレーティングシステム機能へのプラットフォーム非依存のインターフェースを提供。
	//ただし，エラー処理は Go のスタイルです。呼び出し失敗の場合，エラーナンバーではなく， error 型の値を返します。
	"os"

	//signal パッケージは，入力シグナルへのアクセスを実装します。
	"os/signal"

	//syscall パッケージは，低レベル OS 基本要素へのインターフェースを含んでいます。
	"syscall"

	//time パッケージは，時間の測定と表示機能を提供します。
	"time"

	_ "github.com/go-sql-driver/mysql"

	//gorilla/muxはルーティング機能を提供します。
	//gorilla/muxでは*Router型変数が持つメソッドを利用することがメインになります。
	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"

	//最近のエラーパッケージの主流
	"golang.org/x/xerrors"
)

// リファクタリング前
//サーバー起動などのコード
func main() {

	//routerをmux.NewRouter()と定義
	//はじめにNewRouter関数で*Router型の変数を作る必要があります。
	router := mux.NewRouter()

	//func (r *Router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route
	//HandleFuncというメソッド　第一引数パス　第2引数で関数
	//パスとhttp.MethodPosが一致した場合にCreateTodoが実行される
	router.HandleFunc("/todos", CreateTodo).Methods(http.MethodPost)

	//演習　特定のidを取得する（エンドポイントを作る）
	//router.HandleFuncの正式な書き方
	//{id:[0-9]+}　正規表現　0~9を使う全ての桁の数字を受け取れる
	// /todos/5　←例：id5のtodoを取得する
	router.HandleFunc("/todos/{id:[0-9]+}", FetchTodo).Methods(http.MethodGet)

	//演習　todoの一覧を取得する（エンドポイントを作る）
	router.HandleFunc("/todos", FetchTodos).Methods(http.MethodGet)

	//演習　特定のtodoを更新する（エンドポイントを作る）
	router.HandleFunc("/todos/{id:[0-9]+}", UpdateTodo).Methods(http.MethodPut)

	//演習　特定のtodoを削除する（エンドポイントを作る）
	router.HandleFunc("/todos/{id:[0-9]+}", DeleteTodo).Methods(http.MethodDelete)

	//srvをポインター&http.Serverと定義
	//メモリのアドレスを定義
	//サーバーを起動するための構造体
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: router,
	}

	//エラーチャネル　make型　戻り値が1つ　エラーを受け取る場合
	//グレースフルシャットダウン→今受けてる処理まで処理しきる
	errorCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			errorCh <- err
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	// ch は chan
	//v := <- ch
	//受信
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
	//↑ここまでは見なくていい
}

func NewMySQLConnection() (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_db?parseTime=true")
	if err != nil { //エラーがnilじゃなかったら
		return nil, err //エラー(nilじゃない）
	}

	return db, nil //エラーじゃない（nilだよ）
}

//構造体
//jsonで表す場合とデータベースで表す場合
//*は記入なくても良いという意味(nilが入る）
type Todo struct {
	ID                 int       `json:"id"                 db:"id"`
	Title              string    `json:"title"              db:"title"`
	ImplementationDate time.Time `json:"implementationDate" db:"implementation_date"`
	DueDate            time.Time `json:"dueDate"            db:"due_date"`
	StatusID           int       `json:"statusID"           db:"status_id"`
	PriorityID         *int      `json:"priorityID"         db:"priority_id"`
	Memo               *string   `json:"memo"               db:"memo"`
}

//CreateTodoには(w http.ResponseWriter, r *http.Request)を定義する←これは固定
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	//関数とメソッドを使う
	//関数→（型）に紐づかない処理の塊　値レシーバー（json.NewDecoder）
	//メソッド→（型）に紐づく処理の塊　ポインターレシーバー（Decode）
	//todoが型　→（type Todo struct）
	//メソッドチェーン
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		panic(err)
	}
	//↑返り値が一つの場合はこの書き方; err != nil

	//データーベースの接続処理
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

	//チャネルに値が入っていない場合、受信はブロックする。ブロックせずに処理を行いたい場合は select を使う。
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

//演習　特定のidを取得する（エンドポイントを作る）
func FetchTodo(w http.ResponseWriter, r *http.Request) {
	//strconv.Atoi ・・文字列から数字に変換する標準パッケージ
	todoID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	db, err := NewMySQLConnection()
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
	if err = db.QueryRowx(fetchQuery, todoID).StructScan(&outTodo); err != nil {
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

//演習　複数のidを取得する（エンドポイントを作る）
func FetchTodos(w http.ResponseWriter, r *http.Request) {

	db, err := NewMySQLConnection()
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
            priorities.id = todos.priority_id`

	//var outTodo Todo
	rows, err := db.Queryx(fetchQuery)
	if err != nil {
		panic(err)
	}

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.StructScan(&todo); err != nil {
			panic(err)
		}

		todos = append(todos, todo)
	}

	todosOut := Todos{Todos: todos}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&todosOut); err != nil {
		panic(err)
	}
}

type Todos struct {
	Todos []Todo `json:"todos"`
}

//演習　todoを更新する（エンドポイントを作る）
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	var todo Todo

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		panic(err)
	}

	db, err := NewMySQLConnection()
	if err != nil {
		panic(err)
	}

	//INSERT・・データを挿入するクエリー
	//SELECT・・データを取得するクエリー
	//UPDATE・・データを更新するクエリー
	//Delete・・データを削除するクエリー
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

	if _, err := db.Exec(
		updateQuery,
		todo.Title,
		todo.ImplementationDate,
		todo.DueDate,
		todo.StatusID,
		todo.PriorityID,
		todo.Memo,
		todoID,
	); err != nil {
		panic(err)
	}

	todo.ID = todoID

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(&todo); err != nil {
		panic(err)
	}
}

//演習　特定のtodoを削除する（エンドポイントを作る）
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	db, err := NewMySQLConnection()
	if err != nil {
		panic(err)
	}
	updateQuery := `
        DELETE FROM
            todos
        WHERE
            id = ?`

	if _, err := db.Exec(
		updateQuery,
		todoID,
	); err != nil {
		panic(err)
	}

	resp := RespMessage{Message: "削除しました。"}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(&resp); err != nil {
		panic(err)
	}
}

type RespMessage struct {
	Message string `json:"message"`
}
