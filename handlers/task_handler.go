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
	
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id parameter must be integer",
		})
	}

	task := store.GetTaskByID(id)

	if task == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "task not found",
		})
	}

	return c.JSON(task)
}

func CreateTask(c *fiber.Ctx) error {
	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	if task.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "title field is required",
		})
	}

	createdTask := store.AddTask(task)

	return c.Status(fiber.StatusCreated).JSON(createdTask)
}

func UpdateTask(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id parameter must be integer",
		})
	}

	var taskData models.Task
	if err := c.BodyParser(&taskData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	if taskData.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "title field is required",
		})
	}

	updatedTask := store.UpdateTask(id, taskData)
	if updatedTask == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "task not found",
		})
	}

	return c.JSON(updatedTask)
}

func DeleteTask(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id parameter must be integer",
		})
	}

	isDeleted := store.DeleteTask(id)
	if !isDeleted {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "task not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "deleted",
	})
}
