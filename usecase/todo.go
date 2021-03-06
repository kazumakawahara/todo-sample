package usecase

import (
	"github.com/kazumakawahara/todo-sample/apperrors"
	"github.com/kazumakawahara/todo-sample/domain/tododomain"
	"github.com/kazumakawahara/todo-sample/usecase/input"
	"github.com/kazumakawahara/todo-sample/usecase/output"
)

type TodoUsecase interface {
	CreateTodo(in *input.Todo) (*output.Todo, error)
	FetchTodo(id int) (*output.Todo, error)
	FetchTodos() ([]*output.Todo, error)
	UpdateTodo(in *input.Todo) (*output.Todo, error)
	DeleteTodo(id int) error
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

func (u *todoUsecase) FetchTodo(id int) (*output.Todo, error) {
	idVo, err := tododomain.NewID(id)
	if err != nil {
		return nil, apperrors.InvalidParameter
	}

	todoDm, err := u.todoRepository.FetchTodoByID(idVo)
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

func (u *todoUsecase) FetchTodos() ([]*output.Todo, error) {
	todosDm, err := u.todoRepository.FetchTodos()
	if err != nil {
		return nil, err
	}

	todosDto := make([]*output.Todo, len(todosDm))
	for i, todoDm := range todosDm {
		todosDto[i] = &output.Todo{
			ID:                 todoDm.ID().Value(),
			Title:              todoDm.Title().Value(),
			ImplementationDate: todoDm.ImplementationDate().Value(),
			DueDate:            todoDm.DueDate().Value(),
			StatusID:           todoDm.Status().Value(),
			PriorityID:         todoDm.Priority().Value(),
			Memo:               todoDm.Memo().Value(),
		}
	}

	return todosDto, nil
}

func (u *todoUsecase) UpdateTodo(in *input.Todo) (*output.Todo, error) {
	idVo, err := tododomain.NewID(in.ID)
	if err != nil {
		return nil, apperrors.InvalidParameter
	}

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

	statusVo, err := tododomain.NewStatus(in.StatusID)
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

	todoDm := tododomain.NewTodo(
		idVo,
		titleVo,
		implementationDateVo,
		dueDateVo,
		statusVo,
		priorityVo,
		memoVo,
	)

	if _, err = u.todoRepository.UpdateTodo(todoDm); err != nil {
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

func (u *todoUsecase) DeleteTodo(id int) error {
	idVo, err := tododomain.NewID(id)
	if err != nil {
		return apperrors.InvalidParameter
	}

	if err = u.todoRepository.DeleteTodo(idVo); err != nil {
		return err
	}

	return nil
}
