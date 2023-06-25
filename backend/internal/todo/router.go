package todo

import "github.com/gofiber/fiber/v2"

func AddTodoRoutes(app *fiber.App, controller *TodoController) {
	todos := app.Group("api/todos")

	todos.Post("/", controller.create)
	todos.Get("/", controller.getAll)
	todos.Patch("/:id/complete", controller.toggleComplete)
	todos.Patch("/:id/edit", controller.edit)
}
