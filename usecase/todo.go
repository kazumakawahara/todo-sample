package usecase

import (
	"github.com/kazumakawahara/todo-sample/apperrors"
	"github.com/kazumakawahara/todo-sample/domain/tododomain"
	"github.com/kazumakawahara/todo-sample/usecase/input"
	"github.com/kazumakawahara/todo-sample/usecase/output"
)

type TodoUsecase interface {
	CreateTodo(in *input.Todo) (*output.Todo, error)
}

type todoUsecase struct {
	todoRepository tododomain.Repository
}

func NewTodoUsecase(todoRepository tododomain.Repository) *todoUsecase {
	return &todoUsecase{
		todoRepository: todoRepository,
	}
}

func (u *todoUsecase) CreateTodo(in *input.Todo) (*output.Todo, error) {
	titleVo, err := tododomain.NewTitle(in.Title)
	if err != nil {
		return nil, apperrors.InvalidParameter
	}

	implementationDateVo, err := tododomain.NewImplementationDate(in.ImplementationDate)
	if err != nil {
		return nil, apperrors.InvalidParameter
	}

	dueDateVo, err := tododomain.NewDueDate(in.DueDate)
	if err != nil {
		return nil, apperrors.InvalidParameter
	}

	priorityVo, err := tododomain.NewPriority(in.PriorityID)
	if err != nil {
		return nil, apperrors.InvalidParameter
	}

	memoVo, err := tododomain.NewMemo(in.Memo)
	if err != nil {
		return nil, apperrors.InvalidParameter
	}

	todoDm := tododomain.NewTodoWhenUnCreated(
		titleVo,
		implementationDateVo,
		dueDateVo,
		priorityVo,
		memoVo,
	)

	idVo, err := u.todoRepository.CreateTodo(todoDm)
	if err != nil {
		return nil, err
	}

	todoDm, err = u.todoRepository.FetchTodoByID(idVo)
	if err != nil {
		return nil, err
	}

	return &output.Todo{
		ID:                 todoDm.ID().Value(),
		Title:              todoDm.Title().Value(),
		ImplementationDate: todoDm.ImplementationDate().Value(),
		DueDate:            todoDm.DueDate().Value(),
		StatusID:           todoDm.Status().Value(),
		PriorityID:         todoDm.Priority().Value(),
		Memo:               todoDm.Memo().Value(),
	}, nil
}
