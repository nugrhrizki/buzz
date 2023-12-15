package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nugrhrizki/buzz/internal/user"
)

type UserApi struct {
	user *user.Repository
}

func NewUserApi(user *user.Repository) *UserApi {
	return &UserApi{user: user}
}

func (ua *UserApi) CreateUser(c *fiber.Ctx) error {
	payload := new(user.User)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	err := ua.user.CreateUser(payload)
	if err != nil {
		return err
	}

	return nil
}

func (ua *UserApi) GetUser(c *fiber.Ctx) error {
	payload := new(user.User)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	user, err := ua.user.GetUserByUsername(payload.Username)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (ua *UserApi) GetUsers(c *fiber.Ctx) error {
	users, err := ua.user.GetUsers()
	if err != nil {
		return err
	}

	return c.JSON(users)
}

func (ua *UserApi) UpdateUser(c *fiber.Ctx) error {
	payload := new(user.User)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	err := ua.user.UpdateUser(payload)
	if err != nil {
		return err
	}

	return nil
}

func (ua *UserApi) DeleteUser(c *fiber.Ctx) error {
	payload := new(user.User)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	err := ua.user.DeleteUser(payload)
	if err != nil {
		return err
	}

	return nil
}
