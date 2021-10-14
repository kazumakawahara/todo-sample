package tododomain

type Todo struct {
	id                 ID
	title              Title
	implementationDate ImplementationDate
	dueDate            DueDate
	status             Status
	priority           Priority
	memo               Memo
}

func NewTodoWhenUnCreated(
	title Title,
	implementationDate ImplementationDate,
	dueDate DueDate,
	priority Priority,
	memo Memo,
) *Todo {
	return &Todo{
		title:              title,
		implementationDate: implementationDate,
		dueDate:            dueDate,
		status:             TODO, // todo作成時はstatusは作業前
		priority:           priority,
		memo:               memo,
	}
}

func NewTodo(
	id ID,
	title Title,
	implementationDate ImplementationDate,
	dueDate DueDate,
	status Status,
	priority Priority,
	memo Memo) *Todo {
	return &Todo{
		id:                 id,
		title:              title,
		implementationDate: implementationDate,
		dueDate:            dueDate,
		status:             status,
		priority:           priority,
		memo:               memo,
	}
}

func (t *Todo) ID() ID {
	return t.id
}

func (t *Todo) Title() Title {
	return t.title
}

func (t *Todo) ImplementationDate() ImplementationDate {
	return t.implementationDate
}

func (t *Todo) DueDate() DueDate {
	return t.dueDate
}

func (t *Todo) Status() Status {
	return t.status
}

func (t *Todo) Priority() Priority {
	return t.priority
}

func (t *Todo) Memo() Memo {
	return t.memo
}
