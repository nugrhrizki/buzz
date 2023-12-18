package whatsapp

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nugrhrizki/buzz/pkg/whatsapp"
	"github.com/nugrhrizki/buzz/pkg/whatsapp/api"
	"github.com/nugrhrizki/buzz/pkg/whatsapp/user"
	"github.com/rs/zerolog"
)

type WhatsappAPI struct {
	whatsapp *whatsapp.Whatsapp
	api      *api.Api
	log      *zerolog.Logger
	user     *user.Repository
}

func NewWhatsappAPI(api *api.Api, wa *whatsapp.Whatsapp, log *zerolog.Logger, user *user.Repository) *WhatsappAPI {
	return &WhatsappAPI{
		whatsapp: wa,
		api:      api,
		log:      log,
		user:     user,
	}
}

func (w *WhatsappAPI) UserInfo(c *fiber.Ctx) error {
	// Get token from headers or uri parameters
	token := c.Get("token")
	if token == "" {
		token = c.Query("token")
	}

	userInfo, found := w.whatsapp.GetCacheUserInfo(token)

	if !found {
		w.log.Info().Msg("Looking for user information in DB")
		user, err := w.user.GetUserByToken(token)
		if err != nil {
			w.log.Error().Err(err).Msg("failed to get user by token")
			return err
		}

		userInfo = w.whatsapp.UserToUserInfo(user)
		w.whatsapp.UpdateCacheUserInfo(token, userInfo)
		c.Locals("userinfo", userInfo)
		return c.Next()
	}

	userid, err := strconv.Atoi(userInfo.Id)
	if err != nil {
		w.log.Error().Err(err).Msg("failed to convert user id to int")
		return err
	}

	if userid == 0 {
		return fiber.ErrUnauthorized
	}

	c.Locals("userinfo", userInfo)
	return c.Next()
}

func (wa *WhatsappAPI) CreateUser(c *fiber.Ctx) error {
	payload := new(user.User)
	if err := c.BodyParser(payload); err != nil {
		return errors.New("failed to parse body")
	}

	_, err := wa.api.CreateUser(payload)
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func (wa *WhatsappAPI) GetWhatsappUser(c *fiber.Ctx) error {
	user, err := wa.api.GetUsers()
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (wa *WhatsappAPI) GetWhatsappUserById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return errors.New("failed to convert id to int")
	}

	user, err := wa.api.GetUserById(id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (wa *WhatsappAPI) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return errors.New("failed to convert id to int")
	}

	payload := new(user.User)
	if err := c.BodyParser(payload); err != nil {
		return errors.New("failed to parse body")
	}

	user, err := wa.api.GetUserById(id)
	if err != nil {
		return errors.New("failed to get user by id")
	}

	user.Name = payload.Name

	if err := wa.api.UpdateUser(user); err != nil {
		return errors.New("failed to update user")
	}

	return nil
}

func (wa *WhatsappAPI) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return errors.New("failed to convert id to int")
	}

	user, err := wa.api.GetUserById(id)
	if err != nil {
		return errors.New("failed to get user by id")
	}

	if err := wa.api.DeleteUser(user); err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}

func (wa *WhatsappAPI) Connect(c *fiber.Ctx) error {
	payload := new(api.ConnectPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	err := wa.api.Connect(&userInfo, payload)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "failed to connect",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "success connect",
		"data": fiber.Map{
			"webhook": userInfo.Webhook,
			"jid":     userInfo.Jid,
			"events":  userInfo.Events,
			"details": "connected",
		},
	})
}

func (wa *WhatsappAPI) Disconnect(c *fiber.Ctx) error {
	userInfo := c.Locals("userinfo").(user.UserInfo)

	err := wa.api.Disconnect(&userInfo)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) GetWebhook(c *fiber.Ctx) error {
	userInfo := c.Locals("userinfo").(user.UserInfo)

	hook, err := wa.api.GetWebhook(&userInfo)
	if err != nil {
		return err
	}

	return c.JSON(hook)
}

func (wa *WhatsappAPI) SetWebhook(c *fiber.Ctx) error {
	payload := new(api.SetWebhookPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	err := wa.api.SetWebhook(&userInfo, payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) GetQR(c *fiber.Ctx) error {
	userInfo := c.Locals("userinfo").(user.UserInfo)

	qrcode, err := wa.api.GetQR(&userInfo)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "failed to get qrcode",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"qrcode":  qrcode,
	})
}

func (wa *WhatsappAPI) Logout(c *fiber.Ctx) error {
	userInfo := c.Locals("userinfo").(user.UserInfo)

	err := wa.api.Logout(&userInfo)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) GetStatus(c *fiber.Ctx) error {
	userInfo := c.Locals("userinfo").(user.UserInfo)

	status, err := wa.api.GetStatus(&userInfo)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"message": "failed to get status",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "success get status",
		"data":    status,
	})
}

func (wa *WhatsappAPI) SendDocument(c *fiber.Ctx) error {
	payload := new(api.SendDocumentPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	_, err := wa.api.SendDocument(&userInfo, payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) SendAudio(c *fiber.Ctx) error {
	payload := new(api.SendAudioPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	_, err := wa.api.SendAudio(&userInfo, payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) SendImage(c *fiber.Ctx) error {
	payload := new(api.SendImagePayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	_, err := wa.api.SendImage(&userInfo, payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) SendSticker(c *fiber.Ctx) error {
	payload := new(api.SendStickerPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	_, err := wa.api.SendSticker(&userInfo, payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) SendVideo(c *fiber.Ctx) error {
	payload := new(api.SendVideoPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	_, err := wa.api.SendVideo(&userInfo, payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) SendContact(c *fiber.Ctx) error {
	payload := new(api.SendContactPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	_, err := wa.api.SendContact(&userInfo, payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) SendLocation(c *fiber.Ctx) error {
	payload := new(api.SendLocationPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	_, err := wa.api.SendLocation(&userInfo, payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) SendButton(c *fiber.Ctx) error {
	payload := new(api.SendButtonTextPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	_, err := wa.api.SendButton(&userInfo, payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) SendList(c *fiber.Ctx) error {
	payload := new(api.SendListPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	_, err := wa.api.SendList(&userInfo, payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) SendText(c *fiber.Ctx) error {
	payload := new(api.SendTextPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	_, err := wa.api.SendText(&userInfo, payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) CheckUser(c *fiber.Ctx) error {
	payload := new(api.CheckUserPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	user, err := wa.api.CheckUser(&userInfo, payload)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (wa *WhatsappAPI) GetUser(c *fiber.Ctx) error {
	payload := new(api.CheckUserPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	user, err := wa.api.GetUser(&userInfo, *payload)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (wa *WhatsappAPI) GetAvatar(c *fiber.Ctx) error {
	payload := new(api.GetAvatarPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	avatar, err := wa.api.GetAvatar(&userInfo, payload)
	if err != nil {
		return err
	}

	return c.JSON(avatar)
}

func (wa *WhatsappAPI) GetContacts(c *fiber.Ctx) error {
	userInfo := c.Locals("userinfo").(user.UserInfo)

	contacts, err := wa.api.GetContacts(&userInfo)
	if err != nil {
		return err
	}

	c.Response().Header.Set("Content-Type", "application/json")
	return c.Send(contacts)
}

func (wa *WhatsappAPI) SendChatPresence(c *fiber.Ctx) error {
	payload := new(api.ChatPresencePayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	userInfo := c.Locals("userinfo").(user.UserInfo)

	err := wa.api.SendChatPresence(&userInfo, payload)
	if err != nil {
		return err
	}

	return nil
}
