package whatsapp

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nugrhrizki/buzz/pkg/whatsapp/api"
	"github.com/nugrhrizki/buzz/pkg/whatsapp/whatsapp_user"
)

type WhatsappAPI struct {
	api *api.Api
}

func NewWhatsappAPI(api *api.Api) *WhatsappAPI {
	return &WhatsappAPI{api: api}
}

func (wa *WhatsappAPI) CreateUser(c *fiber.Ctx) error {
	payload := new(whatsapp_user.WhatsappUser)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	_, err := wa.api.CreateUser(payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) Connect(c *fiber.Ctx) error {
	payload := new(api.ConnectPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	err := wa.api.Connect(&whatsapp_user.WhatsappUserInfo{}, payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) Disconnect(c *fiber.Ctx) error {
	err := wa.api.Disconnect(&whatsapp_user.WhatsappUserInfo{})
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) GetWebhook(c *fiber.Ctx) error {
	hook, err := wa.api.GetWebhook(&whatsapp_user.WhatsappUserInfo{})
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

	err := wa.api.SetWebhook(&whatsapp_user.WhatsappUserInfo{}, payload)
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) GetQR(c *fiber.Ctx) error {
	qrcode, err := wa.api.GetQR(&whatsapp_user.WhatsappUserInfo{})
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"qrcode": qrcode,
	})
}

func (wa *WhatsappAPI) Logout(c *fiber.Ctx) error {
	err := wa.api.Logout(&whatsapp_user.WhatsappUserInfo{})
	if err != nil {
		return err
	}

	return nil
}

func (wa *WhatsappAPI) GetStatus(c *fiber.Ctx) error {
	status, err := wa.api.GetStatus(&whatsapp_user.WhatsappUserInfo{})
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"status": status,
	})
}

func (wa *WhatsappAPI) SendDocument(c *fiber.Ctx) error {
	payload := new(api.SendDocumentPayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	_, err := wa.api.SendDocument(&whatsapp_user.WhatsappUserInfo{}, payload)
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

	_, err := wa.api.SendAudio(&whatsapp_user.WhatsappUserInfo{}, payload)
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

	_, err := wa.api.SendImage(&whatsapp_user.WhatsappUserInfo{}, payload)
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

	_, err := wa.api.SendSticker(&whatsapp_user.WhatsappUserInfo{}, payload)
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

	_, err := wa.api.SendVideo(&whatsapp_user.WhatsappUserInfo{}, payload)
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

	_, err := wa.api.SendContact(&whatsapp_user.WhatsappUserInfo{}, payload)
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

	_, err := wa.api.SendLocation(&whatsapp_user.WhatsappUserInfo{}, payload)
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

	_, err := wa.api.SendButton(&whatsapp_user.WhatsappUserInfo{}, payload)
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

	_, err := wa.api.SendList(&whatsapp_user.WhatsappUserInfo{}, payload)
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

	_, err := wa.api.SendText(&whatsapp_user.WhatsappUserInfo{}, payload)
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

	user, err := wa.api.CheckUser(&whatsapp_user.WhatsappUserInfo{}, payload)
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

	user, err := wa.api.GetUser(&whatsapp_user.WhatsappUserInfo{}, *payload)
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

	avatar, err := wa.api.GetAvatar(&whatsapp_user.WhatsappUserInfo{}, payload)
	if err != nil {
		return err
	}

	return c.JSON(avatar)
}

func (wa *WhatsappAPI) GetContacts(c *fiber.Ctx) error {
	contacts, err := wa.api.GetContacts(&whatsapp_user.WhatsappUserInfo{})
	if err != nil {
		return err
	}

	return c.JSON(contacts)
}

func (wa *WhatsappAPI) SendChatPresence(c *fiber.Ctx) error {
	payload := new(api.ChatPresencePayload)
	if err := c.BodyParser(payload); err != nil {
		return err
	}

	err := wa.api.SendChatPresence(&whatsapp_user.WhatsappUserInfo{}, payload)
	if err != nil {
		return err
	}

	return nil
}
