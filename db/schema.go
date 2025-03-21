package db

import (
	"database/sql"
	"fmt"
)

func CreateSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		logged_in_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TYPE task_status AS ENUM ('ToDo', 'InProgress', 'Done');

	CREATE TYPE priority_en AS ENUM ('High', 'Medium', 'Low');

	CREATE TABLE IF NOT EXISTS tasks (
		task_id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		priority priority_en NOT NULL DEFAULT 'Medium',
		status task_status NOT NULL DEFAULT 'ToDo',
		assigned_to INTEGER REFERENCES users(user_id),
		description TEXT,
		created_by INTEGER NOT NULL REFERENCES users(user_id),
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS task_suggestions (
		suggestion_id SERIAL PRIMARY KEY,
		task_id INTEGER REFERENCES tasks(task_id),
		user_id INTEGER NOT NULL REFERENCES users(user_id),
		suggestion_text TEXT NOT NULL,
		sub_tasks TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		accepted BOOLEAN DEFAULT FALSE
	);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("error creating schema: %w", err)
	}

	return nil
}

func DropTable(db *sql.DB, tableName string) error {
	drop := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName)

	if tableName == "tasks" {
		dropEnum := "DROP TYPE IF EXISTS task_status, priority_en"
		_, err := db.Exec(dropEnum)
		if err != nil {
			return fmt.Errorf("error dropping enum task_status: %w", err)
		}
	}

	_, err := db.Exec(drop)
	if err != nil {
		return fmt.Errorf("error dropping table %s: %w", tableName, err)
	}

	return nil
}
