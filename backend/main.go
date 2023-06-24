package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"sync"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Done  bool   `json:"done"`
}

var (
	todos      []Todo
	lastTodoID int
	mutex      sync.Mutex
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, http://127.0.0.1:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	todos = []Todo{}

	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.JSON("200")
	})

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.JSON(todos)
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		mutex.Lock()
		lastTodoID++
		todo.ID = lastTodoID
		todos = append(todos, *todo)
		mutex.Unlock()

		return c.JSON(todos)
	})

	app.Patch("/api/todos/:id/edit", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		if err != nil {
			return c.Status(401).SendString("Invalid ID.")
		}

		todo := &Todo{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		mutex.Lock()
		for i, t := range todos {
			if t.ID == id {
				if todo.Title != "" {
					todos[i].Title = todo.Title
				}
				todos[i].Body = todo.Body
				break
			}
		}
		mutex.Unlock()

		return c.JSON(todos)
	})

	app.Patch("/api/todos/:id/done", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		if err != nil {
			return c.Status(401).SendString("Invalid ID.")
		}

		for i, t := range todos {
			if t.ID == id {
				todos[i].Done = !todos[i].Done
				break
			}
		}

		return c.JSON(todos)
	})

	app.Delete("/api/todos/:id/del", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		if err != nil {
			return c.Status(401).SendString("Invalid ID.")
		}

		index := -1
		for i, t := range todos {
			if t.ID == id {
				index = i
				break
			}
		}

		if index == -1 {
			return c.Status(401).SendString("Todo not found.")
		}

		mutex.Lock()
		todos = append(todos[:index], todos[index+1:]...)
		mutex.Unlock()

		return c.JSON(todos)
	})

	log.Fatal(app.Listen(":4000"))
}
