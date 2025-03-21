package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/adarsh-jaiss/zocket/internal/websocket"
	"github.com/adarsh-jaiss/zocket/types"
)

func CreateTaskInStore(db *sql.DB, task types.Task) (int, error) {
	query := `
		INSERT INTO tasks (title, priority, status, assigned_to, description, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING task_id
	`
	var taskID int
	err := db.QueryRow(
		query,
		task.Title,
		task.Priority,
		task.Status,
		task.AssignedTo,
		task.Description,
		task.CreatedBy,
	).Scan(&taskID)

	if err != nil {
		return 0, err
	}

	// Broadcast task creation
	task.TaskID = taskID
	taskJSON, _ := json.Marshal(map[string]interface{}{
		"type": "task_created",
		"data": task,
	})
	websocket.GetManager().BroadcastToAll(taskJSON)

	return taskID, nil
}

func GetTaskFromStore(db *sql.DB, taskID int) (types.Task, error) {
	var task types.Task
	query := `
		SELECT task_id, title, priority, status, assigned_to, description, created_by, created_at, updated_at
		FROM tasks WHERE task_id = $1
	`
	err := db.QueryRow(query, taskID).Scan(
		&task.TaskID,
		&task.Title,
		&task.Priority,
		&task.Status,
		&task.AssignedTo,
		&task.Description,
		&task.CreatedBy,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return types.Task{}, err
	}
	return task, nil
}

func UpdateTaskInStore(db *sql.DB, task types.Task) error {
	query := `
		UPDATE tasks 
		SET 
			title = COALESCE(NULLIF($1, ''), title), 
			priority = CASE WHEN $2 = '' THEN priority ELSE $2::priority_en END,
			status = CASE WHEN $3 = '' THEN status ELSE $3::task_status END,
			assigned_to = COALESCE(NULLIF($4, 0), assigned_to), 
			description = COALESCE(NULLIF($5, ''), description), 
			updated_at = NOW()
		WHERE task_id = $6
	`
	result, err := db.Exec(
		query,
		task.Title,
		task.Priority,
		task.Status,
		task.AssignedTo,
		task.Description,
		task.TaskID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// Broadcast task update
	taskJSON, _ := json.Marshal(map[string]interface{}{
		"type": "task_updated",
		"data": task,
	})
	websocket.GetManager().BroadcastToAll(taskJSON)

	return nil
}

func DeleteTaskFromStore(db *sql.DB, taskID int) error {
	query := `DELETE FROM tasks WHERE task_id = $1`
	result, err := db.Exec(query, taskID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return err
	}
	if rowsAffected == 0 {
		
		return sql.ErrNoRows
	}

	// Broadcast task deletion
	taskJSON, _ := json.Marshal(map[string]interface{}{
		"type": "task_deleted",
		"data": taskID,
	})
	websocket.GetManager().BroadcastToAll(taskJSON)

	return nil
}

func ListTasksFromStore(db *sql.DB) ([]types.Task, error) {
	query := `
		SELECT task_id, title, priority, status, assigned_to, description, created_by, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []types.Task
	for rows.Next() {
		var task types.Task
		err := rows.Scan(
			&task.TaskID,
			&task.Title,
			&task.Priority,
			&task.Status,
			&task.AssignedTo,
			&task.Description,
			&task.CreatedBy,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func StoreSuggestion(db *sql.DB, suggestion types.TaskSuggestion) error {
	// Convert subtasks to JSON for storage
	subTasksJSON, err := json.Marshal(suggestion.SubTasks)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO task_suggestions (task_id, user_id, suggestion_text, sub_tasks, accepted, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING suggestion_id
	`

	var suggestionID int
	err = db.QueryRow(
		query,
		suggestion.TaskID,
		suggestion.UserID,
		suggestion.SuggestionText,
		subTasksJSON,
		suggestion.Accepted,
	).Scan(&suggestionID)

	if err != nil {
		return err
	}

	// Broadcast suggestion creation
	suggestion.SuggestionID = suggestionID
	suggestionJSON, _ := json.Marshal(map[string]interface{}{
		"type": "suggestion_created",
		"data": suggestion,
	})
	websocket.GetManager().BroadcastToAll(suggestionJSON)

	return nil
}

func GetTaskSuggestions(db *sql.DB, taskID int) ([]types.TaskSuggestion, error) {
	query := `
		SELECT suggestion_id, task_id, user_id, suggestion_text, sub_tasks, accepted, created_at
		FROM task_suggestions
		WHERE task_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suggestions []types.TaskSuggestion
	for rows.Next() {
		var suggestion types.TaskSuggestion
		var subTasksJSON []byte

		err := rows.Scan(
			&suggestion.SuggestionID,
			&suggestion.TaskID,
			&suggestion.UserID,
			&suggestion.SuggestionText,
			&subTasksJSON,
			&suggestion.Accepted,
			&suggestion.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse subtasks JSON
		if subTasksJSON != nil {
			if err := json.Unmarshal(subTasksJSON, &suggestion.SubTasks); err != nil {
				return nil, err
			}
		}

		suggestions = append(suggestions, suggestion)
	}

	return suggestions, nil
}
