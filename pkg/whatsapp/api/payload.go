package api

import (
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

type User struct {
	Query        string `json:"query"`
	IsInWhatsapp bool   `json:"is_in_whatsapp"`
	JID          string `json:"jid"`
	VerifiedName string `json:"verified_name"`
}

type UserCollection struct {
	Users []User `json:"users"`
}

type UserInfoCollection struct {
	Users map[types.JID]types.UserInfo `json:"users"`
}

type ConnectPayload struct {
	Subscribe []string `json:"subscribe"`
	Immediate bool     `json:"immediate"`
}

type ChatPresencePayload struct {
	Phone string `json:"phone"`
	State string `json:"state"`
	Media string `json:"media"`
}

type CheckUserPayload struct {
	Phone []string `json:"phone"`
}

type SetWebhookPayload struct {
	WebhookURL string `json:"webhook_url"`
}

type SendDocumentPayload struct {
	Phone       string              `json:"phone"`
	Document    string              `json:"document"`
	FileName    string              `json:"filename"`
	Id          string              `json:"id"`
	ContextInfo waProto.ContextInfo `json:"context_info"`
}

type SendAudioPayload struct {
	Phone       string              `json:"phone"`
	Audio       string              `json:"audio"`
	Caption     string              `json:"caption"`
	Id          string              `json:"id"`
	ContextInfo waProto.ContextInfo `json:"context_info"`
}

type SendImagePayload struct {
	Phone       string              `json:"phone"`
	Image       string              `json:"image"`
	Caption     string              `json:"caption"`
	Id          string              `json:"id"`
	ContextInfo waProto.ContextInfo `json:"context_info"`
}

type SendStickerPayload struct {
	Phone        string              `json:"phone"`
	Sticker      string              `json:"sticker"`
	Id           string              `json:"id"`
	PngThumbnail []byte              `json:"png_thumbnail"`
	ContextInfo  waProto.ContextInfo `json:"context_info"`
}

type SendVideoPayload struct {
	Phone         string              `json:"phone"`
	Video         string              `json:"video"`
	Caption       string              `json:"caption"`
	Id            string              `json:"id"`
	JpegThumbnail []byte              `json:"jpeg_thumbnail"`
	ContextInfo   waProto.ContextInfo `json:"context_info"`
}

type SendContactPayload struct {
	Phone       string              `json:"phone"`
	Id          string              `json:"id"`
	Name        string              `json:"name"`
	Vcard       string              `json:"vcard"`
	ContextInfo waProto.ContextInfo `json:"context_info"`
}

type SendLocationPayload struct {
	Phone       string              `json:"phone"`
	Id          string              `json:"id"`
	Name        string              `json:"name"`
	Latitude    float64             `json:"latitude"`
	Longitude   float64             `json:"longitude"`
	ContextInfo waProto.ContextInfo `json:"context_info"`
}

type SendButtonPayload struct {
	ButtonId   string `json:"button_id"`
	ButtonText string `json:"button_text"`
}

type SendButtonTextPayload struct {
	Phone   string              `json:"phone"`
	Title   string              `json:"title"`
	Buttons []SendButtonPayload `json:"buttons"`
	Id      string              `json:"id"`
}

type SendListRowPayload struct {
	RowId       string `json:"row_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type SendListSectionPayload struct {
	Title string               `json:"title"`
	Rows  []SendListRowPayload `json:"rows"`
}

type SendListPayload struct {
	Phone       string                   `json:"phone"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	ButtonText  string                   `json:"button_text"`
	FooterText  string                   `json:"footer_text"`
	Sections    []SendListSectionPayload `json:"sections"`
	Id          string                   `json:"id"`
}

type SendTextPayload struct {
	Phone       string              `json:"phone"`
	Body        string              `json:"body"`
	Id          string              `json:"id"`
	ContextInfo waProto.ContextInfo `json:"context_info"`
}

type GetAvatarPayload struct {
	Phone   string `json:"phone"`
	Preview bool   `json:"preview"`
}

type GetStatusResponse struct {
	Connected bool `json:"connected"`
	LoggedIn  bool `json:"logged_in"`
}

type GetWebhookResponse struct {
	Webhook   string   `json:"webhook"`
	Subscribe []string `json:"subscribe"`
}
