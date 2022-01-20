package main

import (
	"context"
	"eric/database"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := fiber.New(fiber.Config{
		// Prefork: true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Promethium",
	})

	database.DbConnect()
	defer database.DBPool.Close()

	// app.Static("/", "dist")

	app.Post("/item", CreateItemHandler)

	// app.Patch("/item")

	// app.Get("/*", func(ctx *fiber.Ctx) error {
	// 	return ctx.SendFile("./dist/index.html")
	// })

	app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func CreateItemHandler(c *fiber.Ctx) error {

	conn, err := database.DBPool.Acquire(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}
	defer conn.Release()

	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			fmt.Println("rolled-back!")
			tx.Rollback(context.TODO())
		} else {
			fmt.Println("committed!")
			tx.Commit(context.TODO())
		}
	}()

	// parse body to map
	var databody map[string]interface{}
	err = c.BodyParser(&databody)
	if err != nil {
		return err
	}

	// parse body to custom struct shaped to json
	// var databody struct {
	// 	Item  models.Item `json:"item"`
	// 	User  models.User `json:"user"`
	// 	Token string      `json:"token"`
	// }
	// userBody := &databody.User
	// itemBody := &databody.Item

	// err = c.BodyParser(&databody)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(databody.Token)

	// var item models.Item
	// err = pgxscan.Get(context.Background(), tx, &item, "INSERT INTO items (product, manufacturer, device_type, serial, condition, year) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *", &itemBody.Product, &itemBody.Manufacturer, &itemBody.DeviceType, &itemBody.Serial, &itemBody.Condition, &itemBody.Year)
	// if err != nil {
	// 	return err
	// }

	// var user models.User
	// err = pgxscan.Get(context.Background(), tx, &user, "INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING *", &userBody.Username, &userBody.Password, &userBody.Role)
	// if err != nil {
	// 	return err
	// }

	// return c.Status(201).JSON(fiber.Map{"success": true, "item": item, "user": user})
	return c.Status(201).JSON(fiber.Map{"success": true})
}
