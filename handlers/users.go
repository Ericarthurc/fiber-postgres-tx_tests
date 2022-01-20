package handlers

import (
	"eric/database"
	"eric/models"
	"eric/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

/*
getUserByCredentials |
@Desc: calls models.GetUserByUsername() and then checks passwords for unhashed match
*/
func getUserByCredentials(c *fiber.Ctx, username string, password string) (*models.User, error) {
	user, err := models.GetUserByUsername(database.DBPool, username)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("password incorrect")
	}
	return user, nil
}

/*
GetUsers |
@Desc: Get all users |
@Method: GET |
@Route: "api/v1/users" |
@Auth: Public
*/
func GetUsersHandler(c *fiber.Ctx) error {
	var err error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": nil})
}

/*
GetUser |
@Desc: Get user by id |
@Method: GET |
@Route: "api/v1/users/:id" |
@Auth: Public
*/
func GetUserHandler(c *fiber.Ctx) error {
	var err error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": nil})
}

/*
GetMe |
@Desc: Get me by jwt token |
@Method: GET |
@Route: "api/v1/users/me" |
@Middleware: Authenticate
@Auth: Private
*/
func GetMe(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	token := c.Locals("authToken")
	if userID == nil || token == nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": "missing local"})
	}

	if !models.CheckForToken(database.DBPool, userID.(string), token.(string)) {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": "token not on user"})
	}

	user, err := models.GetUserByID(database.DBPool, userID.(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	jwtToken, err := utils.CreateJWTToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	user.Tokens = append(user.Tokens, jwtToken)
	updatedUser, err := models.UpdateUserTokens(database.DBPool, *user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	utils.SendAuthCookie(c, jwtToken)

	return c.Status(201).JSON(fiber.Map{"success": true, "data": updatedUser})
}

/*
CreateUser |
@Desc: Create new user |
@Method: POST |
@Route: "api/v1/users" |
@Middleware: Authenticate
@Auth: Private
*/
func CreateUserHandler(c *fiber.Ctx) error {
	var userBody models.User

	err := c.BodyParser(&userBody)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	id := uuid.New()
	jwtToken, err := utils.CreateJWTToken(id.String())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	userBody.ID = id.String()
	userBody.Tokens = append(userBody.Tokens, jwtToken)

	hashedPassword, err := utils.HashPassword(userBody.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}
	userBody.Password = hashedPassword

	createdUser, err := models.CreateUser(database.DBPool, userBody)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	utils.SendAuthCookie(c, jwtToken)

	return c.Status(201).JSON(fiber.Map{"success": true, "data": createdUser})
}

/*
LoginUser |
@Desc: Login user by username and password |
@Method: POST |
@Route: "api/v1/users/login" |
@Auth: Public
*/
func LoginUserHandler(c *fiber.Ctx) error {
	var userBody models.User

	err := c.BodyParser(&userBody)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	if userBody.Password == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "no password entered", "data": nil})
	}

	unhashedUser, err := getUserByCredentials(c, userBody.Username, userBody.Password)
	if unhashedUser == nil || err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	jwtToken, err := utils.CreateJWTToken(unhashedUser.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	unhashedUser.Tokens = append(unhashedUser.Tokens, jwtToken)
	updatedUser, err := models.UpdateUserTokens(database.DBPool, *unhashedUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	utils.SendAuthCookie(c, jwtToken)

	return c.Status(201).JSON(fiber.Map{"success": true, "data": updatedUser})
}

/*
LogoutUser |
@Desc: Logout user |
@Method: POST |
@Route: "api/v1/users/logout" |
@Middleware: Authenticate
@Auth: Private
*/
func LogoutUserHandler(c *fiber.Ctx) error {

	userID := c.Locals("userID")
	token := c.Locals("authToken")
	if userID == nil || token == nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": "missing local"})
	}

	if !models.CheckForToken(database.DBPool, userID.(string), token.(string)) {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": "token not on user"})
	}

	err := models.LogoutUser(database.DBPool, userID.(string), token.(string))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}

	c.ClearCookie("authToken")

	return c.Status(201).JSON(fiber.Map{"success": true, "data": "user logged out"})
}

/*
UpdateUser |
@Desc: Update user by id |
@Method: PUT |
@Route: "api/v1/users/:id" |
@Auth: Private
*/
func UpdateUserHandler(c *fiber.Ctx) error {
	var err error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": nil})
}

/*
DeleteUser |
@Desc: Delete user by id |
@Method: DELETE |
@Route: "api/v1/users/:id" |
@Auth: Private
*/
func DeleteUserHandler(c *fiber.Ctx) error {
	var err error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "data": nil})
}
