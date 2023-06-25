package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/misternano/practice-todo/config"
	"github.com/misternano/practice-todo/internal/storage"
	"github.com/misternano/practice-todo/internal/todo"
	"github.com/misternano/practice-todo/pkg/shutdown"
	"os"
	"time"
)

func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	env, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	cleanup, err := run(env)
	defer cleanup()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	shutdown.Gracefully()
}

func run(env config.EnvVars) (func(), error) {
	app, cleanup, err := buildServer(env)
	if err != nil {
		return nil, err
	}

	go func() {
		app.Listen("0.0.0.0:" + env.PORT)
	}()

	return func() {
		cleanup()
		app.Shutdown()
	}, nil
}

func buildServer(env config.EnvVars) (*fiber.App, func(), error) {
	db, err := storage.BootstrapMongo(env.MONGODB_URI, env.MONGODB_NAME, 10*time.Second)
	if err != nil {
		return nil, nil, err
	}

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	todoStore := todo.NewTodoStorage(db)
	todoController := todo.NewTodoController(todoStore)
	todo.AddTodoRoutes(app, todoController)

	return app, func() {
		storage.CloseMongo(db)
	}, nil
}
