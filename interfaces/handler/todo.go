package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/kazumakawahara/todo-sample/apperrors"
	"github.com/kazumakawahara/todo-sample/interfaces/presenter"
	"github.com/kazumakawahara/todo-sample/usecase"
	"github.com/kazumakawahara/todo-sample/usecase/input"
	"github.com/kazumakawahara/todo-sample/usecase/output"
)

type todoHandler struct {
	todoUsecase usecase.TodoUsecase
}

func NewTodoHandler(todoUsecase usecase.TodoUsecase) *todoHandler {
	return &todoHandler{
		todoUsecase: todoUsecase,
	}
}

func (h *todoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var in input.Todo
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		presenter.ErrorJSON(w, apperrors.InvalidParameter)
		return
	}

	out, err := h.todoUsecase.CreateTodo(&in)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusCreated, out)
}

func (h *todoHandler) FetchTodo(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		presenter.ErrorJSON(w, apperrors.InvalidParameter)
		return
	}

	out, err := h.todoUsecase.FetchTodo(todoID)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusOK, out)
}

func (h *todoHandler) FetchTodos(w http.ResponseWriter, r *http.Request) {
	out, err := h.todoUsecase.FetchTodos()
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusOK, out)
}

func (h *todoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		presenter.ErrorJSON(w, apperrors.InvalidParameter)
		return
	}

	in := input.Todo{
		ID: todoID,
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		presenter.ErrorJSON(w, apperrors.InvalidParameter)
		return
	}

	out, err := h.todoUsecase.UpdateTodo(&in)
	if err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	presenter.JSON(w, http.StatusOK, out)
}

func (h *todoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		presenter.ErrorJSON(w, apperrors.InvalidParameter)
		return
	}

	if err := h.todoUsecase.DeleteTodo(todoID); err != nil {
		presenter.ErrorJSON(w, err)
		return
	}

	resp := output.DeleteMessage{Message: "削除しました。"}

	presenter.JSON(w, http.StatusOK, resp)
}
