package user

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nugrhrizki/buzz/internal/user"
	"github.com/nugrhrizki/buzz/pkg/password"
)

type UserApi struct {
	user     *user.Repository
	password *password.Password
}

func NewUserApi(user *user.Repository, password *password.Password) *UserApi {
	return &UserApi{user: user, password: password}
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
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return errors.New("failed to convert id to int")
	}

	payload := new(user.User)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	user, err := ua.user.GetUserById(id)
	if err != nil {
		return errors.New("failed to get user by id")
	}

	user.Name = payload.Name

	if user.Username != payload.Username {
		if _, err := ua.user.GetUserByUsername(payload.Username); err == nil {
			return errors.New("username already used")
		}

		user.Username = payload.Username
	}

	if payload.Password != "" {
		hashedPassword, err := ua.password.GenerateHashPassword(payload.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	user.Confirmed = payload.Confirmed
	user.Whatsapp = payload.Whatsapp
	user.Email = payload.Email
	user.RoleId = payload.RoleId

	if err := ua.user.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (ua *UserApi) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return errors.New("failed to convert id to int")
	}

	user, err := ua.user.GetUserById(id)
	if err != nil {
		return errors.New("failed to get user by id")
	}

	if err := ua.user.DeleteUser(user); err != nil {
		return err
	}

	return nil
}
