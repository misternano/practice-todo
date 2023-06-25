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

type createTodoResponse struct {
	ID string `json:"id"`
}

type toggleCompleteRequest struct {
	ID string `json:"id"`
}

type editTodoRequest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
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

	return c.Status(fiber.StatusCreated).JSON(createTodoResponse{
		ID: id,
	})
}

func (t *TodoController) toggleComplete(c *fiber.Ctx) error {
	var req toggleCompleteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	err := t.storage.toggleComplete(req.ID, c.Context())
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
