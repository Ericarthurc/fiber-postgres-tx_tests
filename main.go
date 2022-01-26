package main

import (
	"context"
	"eric/database"
	"eric/models"
	"fmt"
	"os"

	"github.com/georgysavva/scany/pgxscan"
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
	// var databody map[string]interface{}
	// err = c.BodyParser(&databody)
	// if err != nil {
	// 	return err
	// }

	// parse body to custom struct shaped to json
	var databody struct {
		Service models.Service `json:"service"`
		User    models.User    `json:"user"`
		Token   string         `json:"token"`
	}
	userBody := &databody.User
	serviceBody := &databody.Service

	err = c.BodyParser(&databody)
	if err != nil {
		return err
	}

	fmt.Println(userBody)
	fmt.Println(serviceBody)

	var service models.Service
	err = pgxscan.Get(context.Background(), tx, &service, "INSERT INTO services (time, seats) VALUES ($1, $2) RETURNING *", &serviceBody.Time, &serviceBody.Seats)
	if err != nil {
		return err
	}

	userBody.Servicetime = service.Time
	userBody.Serviceid = service.ID

	fmt.Println(userBody)

	var user models.User
	err = pgxscan.Get(context.Background(), tx, &user, "INSERT INTO users (name, email, userseats, servicetime, serviceid) VALUES ($1, $2, $3, $4, $5) RETURNING *", &userBody.Name, &userBody.Email, &userBody.Userseats, &userBody.Servicetime, &userBody.Serviceid)
	if err != nil {
		fmt.Println("here")
		return err
	}

	fmt.Println(user)

	return c.Status(201).JSON(fiber.Map{"success": true, "service": service, "user": user})
	// return c.Status(201).JSON(fiber.Map{"success": true})
}
