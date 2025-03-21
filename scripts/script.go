package main

import (
	"fmt"

	"github.com/adarsh-jaiss/zocket/db"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		panic(err)
	}
}

func main() {
	conn, err := db.Connect()
	if err != nil {
		fmt.Printf("error connecting database: %v", err)
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to DB")
	fmt.Println("Dropping tables...")

	// Drop tables in reverse order of dependencies
	tables := []string{"task_suggestions", "task_updates", "tasks", "users"}
	for _, table := range tables {
		fmt.Printf("dropping %v table\n", table)
		if table == "tasks" {
			// Drop enums first if dropping tasks table
			dropEnum := "DROP TYPE IF EXISTS task_status, priority_en CASCADE"
			_, err := conn.Exec(dropEnum)
			if err != nil {
				fmt.Printf("Error dropping enums: %v\n", err)
				return
			}
		}
		err := db.DropTable(conn, table)
		if err != nil {
			fmt.Printf("Error dropping table %s: %v\n", table, err)
			return
		}
	}

	fmt.Println("all tables dropped successfully")
	fmt.Println("creating schema...")

	err = db.CreateSchema(conn)
	if err != nil {
		fmt.Printf("Error creating schema: %v\n", err)
		panic(err)
	}

	fmt.Println("schema created...")
}