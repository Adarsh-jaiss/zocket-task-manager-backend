package tasks

import (
	"database/sql"
	"fmt"

	"github.com/adarsh-jaiss/zocket/internal/ai"
	users "github.com/adarsh-jaiss/zocket/internal/user"
	"github.com/adarsh-jaiss/zocket/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func CreateTask(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var task types.Task
		if err := c.BodyParser(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Get user ID from JWT token
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		userID := int(claims["user_id"].(float64))
		task.CreatedBy = userID

		// Set default values if not provided
		if task.Status == "" {
			task.Status = types.ToDo
		}
		if task.Priority == "" {
			task.Priority = types.Medium
		}

		taskID, err := CreateTaskInStore(db, task)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create task",
			})
		}

		task.TaskID = taskID
		return c.Status(fiber.StatusCreated).JSON(task)
	}
}

func GetTask(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		taskID, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid task ID",
			})
		}

		task, err := GetTaskFromStore(db, taskID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Task not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve task",
			})
		}

		return c.Status(fiber.StatusOK).JSON(task)
	}
}

func UpdateTask(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		taskID, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid task ID",
			})
		}

		var task types.Task
		if err := c.BodyParser(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Get existing task to check ownership
		existingTask, err := GetTaskFromStore(db, taskID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Task not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve task",
			})
		}

		// Get user ID from JWT token
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		userID := int(claims["user_id"].(float64))

		// Only allow task creator or assignee to update
		if existingTask.CreatedBy != userID && existingTask.AssignedTo != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Not authorized to update this task",
			})
		}

		task.TaskID = taskID
		task.CreatedBy = existingTask.CreatedBy
		task.CreatedAt = existingTask.CreatedAt

		if err := UpdateTaskInStore(db, task); err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update task",
			})
		}

		return c.Status(fiber.StatusOK).JSON(task)
	}
}

func DeleteTask(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		taskID, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid task ID",
			})
		}

		// Get existing task to check ownership
		existingTask, err := GetTaskFromStore(db, taskID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Task not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve task",
			})
		}

		// Get user ID from JWT token
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		userID := int(claims["user_id"].(float64))

		// Only allow task creator to delete
		if existingTask.CreatedBy != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Not authorized to delete this task",
			})
		}

		if err := DeleteTaskFromStore(db, taskID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete task",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Task deleted successfully",
		})
	}
}

func ListTasks(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tasks, err := ListTasksFromStore(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve tasks",
			})
		}

		users, err := users.GetAllUsers(db)
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to map users to tasks",
			})
		}

		// Create a map of user IDs to full names
		userMap := make(map[int]string)
		for _, user := range users {
			userMap[user.ID] = user.FirstName + " " + user.LastName
		}

		// Replace AssignedTo IDs with full names in tasks
		for i := range tasks {
			if tasks[i].AssignedTo != 0 {
				if name, ok := userMap[tasks[i].AssignedTo]; ok {
					tasks[i].AssignedToName = name
				}
			}
		}

		return c.Status(fiber.StatusOK).JSON(tasks)
	}
}

func AnalyzeTask(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		taskID, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid task ID",
			})
		}

		var req types.AITaskBreakdownRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		req.TaskID = taskID

		// Get existing task
		task, err := GetTaskFromStore(db, req.TaskID)
		if err != nil {
			fmt.Println(err)
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Task not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve task",
			})
		}

		// Get user ID from JWT token
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		userID := int(claims["user_id"].(float64))

		// Only allow task creator or assignee to request analysis
		if task.CreatedBy != userID && task.AssignedTo != userID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Not authorized to analyze this task",
			})
		}

		// If additional description is provided, append it to task description
		if req.Description != "" {
			task.Description += "\n\nAdditional Context:\n" + req.Description
		}

		// Initialize Gemini client
		gemini, err := ai.NewGeminiClient()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to initialize AI service",
			})
		}
		defer gemini.Close()

		// Get AI analysis
		analysis, err := gemini.AnalyzeTask(task)
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to analyze task",
			})
		}

		// Store suggestions in database
		for _, suggestion := range analysis.Suggestions {
			suggestion.TaskID = task.TaskID
			suggestion.UserID = userID
			if err := StoreSuggestion(db, suggestion); err != nil {
				// Log error but continue
				c.App().Config().ErrorHandler(c, err)
			}
		}

		return c.Status(fiber.StatusOK).JSON(analysis)
	}
}
