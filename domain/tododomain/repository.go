package tododomain

type Repository interface {
	CreateTodo(todo *Todo) (ID, error)
	FetchTodoByID(id ID) (*Todo, error)
	FetchTodos() ([]*Todo, error)
	UpdateTodo(todo *Todo) (ID, error)
	DeleteTodo(id ID) error
}
