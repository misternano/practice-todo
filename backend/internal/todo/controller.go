package todo

import "github.com/gofiber/fiber/v2"

type TodoController struct {
	storage *TodoStorage
}

func NewTodoController(storage *TodoStorage) *TodoController {
	return &TodoController{
		storage: storage,
	}
}

type createTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type editTodoRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type todoRequestResponse struct {
	ID string `json:"id"`
}

func (t *TodoController) getAll(c *fiber.Ctx) error {
	todos, err := t.storage.getAllTodos(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get todos",
		})
	}

	return c.JSON(todos)
}

func (t *TodoController) create(c *fiber.Ctx) error {
	var req createTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	id, err := t.storage.createTodo(req.Title, req.Description, false, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create todo",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(todoRequestResponse{
		ID: id,
	})
}

func (t *TodoController) complete(c *fiber.Ctx) error {
	var req todoRequestResponse
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	err := t.storage.completeTodo(req.ID, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to toggle todo completion",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo completion toggled successfully",
	})
}

func (t *TodoController) edit(c *fiber.Ctx) error {
	var req editTodoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	err := t.storage.editTodo(req.ID, req.Title, req.Description, req.Completed, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to edit todo",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo edited successfully",
	})
}

func (t *TodoController) delete(c *fiber.Ctx) error {
	var req todoRequestResponse
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	err := t.storage.delete(req.ID, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to delete todo",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo deleted successfully",
	})
}
