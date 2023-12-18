package auth

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nugrhrizki/buzz/internal/role"
	"github.com/nugrhrizki/buzz/internal/user"
	"github.com/nugrhrizki/buzz/pkg/password"
	"github.com/rs/zerolog"
)

type HashParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
	saltLength  uint32
}

type AuthApi struct {
	user     *user.Repository
	role     *role.Repository
	log      *zerolog.Logger
	password *password.Password
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func NewAuthApi(user *user.Repository, role *role.Repository, log *zerolog.Logger, password *password.Password) *AuthApi {
	return &AuthApi{
		user:     user,
		role:     role,
		log:      log,
		password: password,
	}
}

func (a *AuthApi) IdentifyUser(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)

	user, err := a.user.GetUserById(int(claims["uid"].(float64)))
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to get user")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"title":   "Unauthorized",
			"message": "You are not authorized",
		})
	}

	user.Password = ""
	return c.JSON(fiber.Map{
		"status":  "success",
		"title":   "Authorized",
		"message": "You are authorized",
		"data":    user,
	})
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *AuthApi) Login(c *fiber.Ctx) error {
	request := new(LoginPayload)
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"title":   "Oops, something went wrong",
			"message": err.Error(),
		})
	}

	user, err := a.user.GetUserByUsername(request.Username)
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to get user")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"title":   "Login Failed",
			"message": "Username or password is incorrent",
		})
	}

	isMatch, err := a.password.CompareHashPassword(request.Password, user.Password)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"title":   "Login Failed",
			"message": "Username or password is incorrent",
		})
	}

	if !isMatch {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"title":   "Login Failed",
			"message": "Username or password is incorrent",
		})
	}

	claims := jwt.MapClaims{
		"uid":      user.Id,
		"nama":     user.Name,
		"username": user.Username,
		"rid":      user.RoleId,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := "secret"
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"title":   "Oops, something went wrong",
			"message": "Internal server error",
		})
	}

	cookie := fiber.Cookie{
		Name:     "auth-token",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "lax",
	}

	c.Cookie(&cookie)

	user.Password = ""

	return c.JSON(fiber.Map{
		"status":  "success",
		"title":   "Login success",
		"message": "Were redirecting you to dashboard in a second",
		"data":    user,
	})
}

func (a *AuthApi) CreateUser(c *fiber.Ctx) error {
	payload := new(user.User)

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"title":   "Oops, something went wrong",
			"message": err.Error(),
		})
	}

	hashedPassword, err := a.password.GenerateHashPassword(payload.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"title":   "Oops, something went wrong",
			"message": err.Error(),
		})
	}

	newUser := user.User{
		Name:     payload.Name,
		Username: payload.Username,
		Password: hashedPassword,
		RoleId:   payload.RoleId,
	}

	if err := a.user.CreateUser(&newUser); err != nil {
		a.log.Error().Err(err).Msg("Failed to create user")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"title":   "Failed register user",
			"message": err.Error(),
		})
	}

	newUser.Password = ""

	return c.JSON(fiber.Map{
		"status":  "success",
		"title":   "Register success",
		"message": "A new user is registered",
		"data":    newUser,
	})
}

func (a *AuthApi) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "auth-token",
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		SameSite: "lax",
	})
	return c.JSON(fiber.Map{
		"status":  "success",
		"title":   "Logout successfully",
		"message": "Have a nice day",
	})
}
