package types

type TaskStatus string

const (
	ToDo       TaskStatus = "ToDo"
	InProgress TaskStatus = "InProgress"
	Done       TaskStatus = "Done"
)

type TaskPriority string

const (
	High   TaskPriority = "High"
	Medium TaskPriority = "Medium"
	Low    TaskPriority = "Low"
)

type Task struct {
	TaskID      int          `json:"id" db:"task_id"`
	Title       string       `json:"title,omitempty" db:"title"`
	Priority    TaskPriority `json:"priority,omitempty" db:"priority"`
	Status      TaskStatus   `json:"status,omitempty" db:"status"`
	AssignedTo  int          `json:"assigned_to,omitempty" db:"assigned_to"`
	AssignedToName string	`json:"assigned_to_name,omitempty" db:"assigned_to_name"`
	Description string       `json:"description,omitempty" db:"description"`
	CreatedBy   int          `json:"created_by" db:"created_by"`
	CreatedAt   string       `json:"created_at" db:"created_at"`
	UpdatedAt   string       `json:"updated_at" db:"updated_at"`
}

type ChangeType string

const (
	Assignee ChangeType = "Assignee"
	Status   ChangeType = "Status"
	Priority ChangeType = "Priority"
)

type TaskUpdate struct {
	UpdateID    int           `json:"id" db:"update_id"`
	TaskID      int           `json:"task_id" db:"task_id"`
	UserID      int           `json:"user_id" db:"user_id"`
	ChangeType  ChangeType    `json:"change_type" db:"change_type"`
	OldAssignee *int          `json:"old_assignee,omitempty" db:"old_assignee"`
	NewAssignee *int          `json:"new_assignee,omitempty" db:"new_assignee"`
	OldStatus   *TaskStatus   `json:"old_status,omitempty" db:"old_status"`
	NewStatus   *TaskStatus   `json:"new_status,omitempty" db:"new_status"`
	OldPriority *TaskPriority `json:"old_priority,omitempty" db:"old_priority"`
	NewPriority *TaskPriority `json:"new_priority,omitempty" db:"new_priority"`
	UpdatedAt   string        `json:"updated_at" db:"updated_at"`
}

// AI Task Suggestion types
type TaskSuggestion struct {
	SuggestionID   int    `json:"id" db:"suggestion_id"`
	TaskID         int    `json:"task_id" db:"task_id"`
	UserID         int    `json:"user_id" db:"user_id"`
	SuggestionText string `json:"suggestion_text" db:"suggestion_text"`
	SubTasks       []Task `json:"sub_tasks,omitempty"`
	Accepted       bool   `json:"accepted" db:"accepted"`
	CreatedAt      string `json:"created_at" db:"created_at"`
}

type AITaskBreakdownRequest struct {
	TaskID      int    `json:"task_id"`
	Description string `json:"description"`
	Context     string `json:"context,omitempty"`
}

type AITaskBreakdownResponse struct {
	TaskID      int              `json:"task_id"`
	Suggestions []TaskSuggestion `json:"suggestions"`
	Analysis    string           `json:"analysis"`
}
