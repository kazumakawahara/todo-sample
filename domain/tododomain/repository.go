package tododomain

type Repository interface {
	CreateTodo(todo *Todo) (ID, error)
	FetchTodoByID(id ID) (*Todo, error)
}
