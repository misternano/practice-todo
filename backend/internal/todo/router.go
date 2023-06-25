package todo

import "github.com/gofiber/fiber/v2"

func AddTodoRoutes(app *fiber.App, controller *TodoController) {
	todos := app.Group("api/todos")

	todos.Get("/", controller.getAll)
	todos.Post("/", controller.create)
	todos.Patch("/:id/complete", controller.complete)
	todos.Patch("/:id/edit", controller.edit)
	todos.Delete("/:id/delete", controller.delete)

}
