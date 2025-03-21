package main

import (
	"fmt"
	"log"

	"github.com/adarsh-jaiss/zocket/db"
	"github.com/adarsh-jaiss/zocket/internal/middleware"
	tasks "github.com/adarsh-jaiss/zocket/internal/tasks"
	users "github.com/adarsh-jaiss/zocket/internal/user"
	wsmanager "github.com/adarsh-jaiss/zocket/internal/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
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

	// Initialize WebSocket manager
	wsmanager.InitManager()

	app := fiber.New()
	app.Use(logger.New()) // Add logging middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
	}))

	// WebSocket middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	api := app.Group("/api")

	// public routes (auth)
	auth := api.Group("/auth")
	auth.Post("/signup", users.Signup(conn))
	auth.Post("/signin", users.SignIn(conn))

	// v1 (protected routes)
	v1 := api.Group("/v1", middleware.JWTProtected())

	// WebSocket endpoint
	v1.Get("/ws", wsmanager.WebsocketHandler())

	// user routes
	user := v1.Group("/user")
	user.Get("/:id", users.GetUser(conn))
	user.Get("",users.FetchAllUsers(conn))

	// task routes
	tasksGroup := v1.Group("/tasks")
	tasksGroup.Post("/", tasks.CreateTask(conn))
	tasksGroup.Get("/", tasks.ListTasks(conn))
	tasksGroup.Get("/:id", tasks.GetTask(conn))
	tasksGroup.Put("/:id", tasks.UpdateTask(conn))
	tasksGroup.Delete("/:id", tasks.DeleteTask(conn))
	tasksGroup.Post("/:id/analyze", tasks.AnalyzeTask(conn))

	log.Fatal(app.Listen(":8000"))
}
