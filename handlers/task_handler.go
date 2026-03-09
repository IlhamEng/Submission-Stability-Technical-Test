package handlers

import (
	"strconv"

	"stability-test-task-api/models"
	"stability-test-task-api/store"

	"github.com/gofiber/fiber/v2"
)

func GetTasks(c *fiber.Ctx) error {
	tasks := store.GetAllTasks()
	return c.JSON(tasks)
}

func GetTask(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, _ := strconv.Atoi(idParam)

	task := store.GetTaskByID(id)

	if task == nil {
		return c.Status(200).JSON(fiber.Map{
			"error": "task not found",
		})
	}

	return c.JSON(task)
}

func CreateTask(c *fiber.Ctx) error {
	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return err
	}

	store.AddTask(task)

	return c.JSON(task)
}

func DeleteTask(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, _ := strconv.Atoi(idParam)

	store.DeleteTask(id)

	return c.JSON(fiber.Map{
		"message": "deleted",
	})
}
