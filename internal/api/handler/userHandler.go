package handler

import (
	"strconv"

	"github.com/OzyKleyton/studio-api/internal/api/middleware"
	"github.com/OzyKleyton/studio-api/internal/api/router"
	"github.com/OzyKleyton/studio-api/internal/model"
	"github.com/OzyKleyton/studio-api/internal/model/user"
	"github.com/OzyKleyton/studio-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) Routes() router.Router {
	return func(route fiber.Router) {
		route.Post("/login", uh.LoginHandler)
		route.Post("/users", uh.CreateUserHandler)

		protected := route.Group("", middleware.AuthJwt())
		users := protected.Group("/users")

		users.Get("/", uh.FindAllUsersHandler)
		users.Get("/:email", uh.FindUserByEmailHandler)
		users.Put("/:id", uh.UpdateUserHandler)
		users.Delete("/:id", uh.DeleteUserHandler)
		// users.Get("/me", uh.GetCurrentUserHandler)
	}
}

func (uh *UserHandler) CreateUserHandler(c *fiber.Ctx) error {
	userReq := new(user.UserReq)
	if err := c.BodyParser(userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.NewErrorResponse(err, fiber.ErrBadRequest))
	}

	res := uh.service.CreateUser(userReq)

	return c.Status(res.Status).JSON(res)
}

func (uh *UserHandler) FindAllUsersHandler(c *fiber.Ctx) error {
	res := uh.service.FindAllUsers()

	return c.Status(res.Status).JSON(res)
}

func (uh *UserHandler) FindUserByEmailHandler(c *fiber.Ctx) error {
	userEmail := c.Params("email")

	res := uh.service.FindUserByEmail(userEmail)

	return c.Status(res.Status).JSON(res)
}

func (uh *UserHandler) UpdateUserHandler(c *fiber.Ctx) error {
	userReq := new(user.UserReq)
	if err := c.BodyParser(userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.NewErrorResponse(err, fiber.ErrBadRequest))
	}

	userID, _ := strconv.Atoi(c.Params("id", "0"))

	res := uh.service.UpdateUser(uint(userID), userReq)

	return c.Status(res.Status).JSON(res)
}

func (uh *UserHandler) DeleteUserHandler(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("id", "0"))

	res := uh.service.DeleteUser(uint(userID))

	return c.Status(res.Status).JSON(res)
}

func (uh *UserHandler) LoginHandler(c *fiber.Ctx) error {
	userReq := new(user.Login)
	if err := c.BodyParser(userReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.NewErrorResponse(err, fiber.ErrBadRequest))
	}

	res := uh.service.Login(*userReq)

	return c.Status(res.Status).JSON(res)
}
