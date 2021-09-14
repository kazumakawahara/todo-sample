package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kazumakawahara/todo-sample/apperrors"
	"github.com/kazumakawahara/todo-sample/interfaces/presenter"
	"github.com/kazumakawahara/todo-sample/usecase"
	"github.com/kazumakawahara/todo-sample/usecase/input"
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
