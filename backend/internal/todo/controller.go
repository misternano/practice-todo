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

func (t *TodoController) getAll(c *fiber.Ctx) error {
	todos, err := t.storage.getAllTodos(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get todos",
		})
	}

	return c.JSON(todos)
}
