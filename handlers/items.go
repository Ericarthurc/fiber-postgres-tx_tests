package handlers

import (
	"eric/database"
	"eric/models"

	"github.com/gofiber/fiber/v2"
)

func GetItemsHandler(c *fiber.Ctx) error {
	items, err := models.GetItems(database.DBPool)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": items})
}

func GetItemsBySearchHandler(c *fiber.Ctx) error {
	searchQuery := c.Query("q")

	items, err := models.GetItemsBySearch(database.DBPool, searchQuery)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "data": items})
}

func GetItemHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")

	item, err := models.GetItem(database.DBPool, idParam)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "data": item})
}

func CreateItemHandler(c *fiber.Ctx) error {
	var itemBody models.Item

	err := c.BodyParser(&itemBody)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	createdItem, err := models.CreateItem(database.DBPool, itemBody)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "data": createdItem})
}

func UpdateItemHandler(c *fiber.Ctx) error {
	// idParam := c.Params("id")

	var err error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": err})
}

func DeleteItemHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")

	if !models.DeleteItem(database.DBPool, idParam) {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": nil})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": nil})
}
