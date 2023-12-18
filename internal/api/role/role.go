package role

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nugrhrizki/buzz/internal/role"
)

type RoleApi struct {
	role *role.Repository
}

func NewRoleApi(role *role.Repository) *RoleApi {
	return &RoleApi{role: role}
}

func (ra *RoleApi) CreateRole(c *fiber.Ctx) error {
	payload := new(role.Role)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	err := ra.role.CreateRole(payload)
	if err != nil {
		return err
	}

	return nil
}

func (ra *RoleApi) GetRole(c *fiber.Ctx) error {
	payload := new(role.Role)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	role, err := ra.role.GetRoleByRolename(payload.Name)
	if err != nil {
		return err
	}

	return c.JSON(role)
}

func (ra *RoleApi) GetRoles(c *fiber.Ctx) error {
	roles, err := ra.role.GetRoles()
	if err != nil {
		return err
	}

	return c.JSON(roles)
}

func (ra *RoleApi) UpdateRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return errors.New("failed to convert id to int")
	}

	payload := new(role.Role)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	role, err := ra.role.GetRoleById(id)
	if err != nil {
		return errors.New("failed to get user by id")
	}

	role.Name = payload.Name
	role.Actions = payload.Actions

	if err := ra.role.UpdateRole(role); err != nil {
		return err
	}

	return nil
}

func (ra *RoleApi) DeleteRole(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return errors.New("failed to convert id to int")
	}

	role, err := ra.role.GetRoleById(id)
	if err != nil {
		return errors.New("failed to get role by id")
	}

	if err := ra.role.DeleteRole(role); err != nil {
		return err
	}

	return nil
}
