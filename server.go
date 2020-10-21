package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Todo struct!
type Todo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var todos = []Todo{
	{ID: 1, Name: "abc", Completed: false},
	{ID: 2, Name: "def", Completed: true},
}

func main() {

	app := fiber.New()
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/todo", getTodo)
	app.Post("/todo", postTodo)
	app.Get("/todo/:id", getSingleTodo)
	app.Delete("/todo/:id", deleteTodo)
	app.Patch("/todo/:id", patchTodo)

	app.Listen(":5000")
}

func getTodo(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(todos)
}

func getSingleTodo(c *fiber.Ctx) error {
	paramID := c.Params("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID invalid.",
		})
	}

	for _, todo := range todos {
		if todo.ID == id {
			return c.Status(fiber.StatusFound).JSON(todo)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Record not found"})
}

func deleteTodo(c *fiber.Ctx) error {
	paramID := c.Params("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid id.",
		})
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[0:i], todos[i+1:]...)
			return c.Status(fiber.StatusOK).JSON(todo)
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Record not found"})
}

func postTodo(c *fiber.Ctx) error {
	type request struct {
		Name      string `json:"name"`
		Completed bool   `json:"completed"`
	}

	var body request

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	todo := Todo{
		ID:        len(todos) + 1,
		Name:      body.Name,
		Completed: body.Completed,
	}

	todos = append(todos, todo)

	return c.Status(fiber.StatusCreated).JSON(todo)
}

func patchTodo(c *fiber.Ctx) error {
	type request struct {
		Name      string `json:"name"`
		Completed bool   `json:"completed"`
	}
	var body request

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	paramID := c.Params("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid id.",
		})
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos[i] = Todo{
				ID:        id,
				Name:      body.Name,
				Completed: body.Completed,
			}
			return c.Status(fiber.StatusOK).JSON(todos[i])
		}
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Record not found"})
}
