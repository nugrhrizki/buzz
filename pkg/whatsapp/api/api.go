package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"

	"github.com/nugrhrizki/buzz/pkg/utils"
	"github.com/nugrhrizki/buzz/pkg/whatsapp"
	"github.com/nugrhrizki/buzz/pkg/whatsapp/user"
	"github.com/rs/zerolog"
	"github.com/vincent-petithory/dataurl"
)

type Api struct {
	log      *zerolog.Logger
	whatsapp *whatsapp.Whatsapp
	users    *user.Repository
}

func New(
	log *zerolog.Logger,

	whatsapp *whatsapp.Whatsapp,

	users *user.Repository,
) *Api {
	return &Api{
		log:      log,
		whatsapp: whatsapp,
		users:    users,
	}
}

func (a *Api) CreateUser(payload *user.User) (*user.User, error) {
	user := user.User{
		Name:  payload.Name,
		Token: payload.Token,
	}

	if err := a.users.CreateUser(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *Api) GetUsers() ([]user.User, error) {
	return a.users.GetUsers()
}

func (a *Api) GetUserById(id int) (*user.User, error) {
	return a.users.GetUserById(id)
}

func (a *Api) DeleteUser(payload *user.User) error {
	return a.users.DeleteUser(payload)
}

func (a *Api) UpdateUser(payload *user.User) error {
	return a.users.UpdateUser(payload)
}

func (a *Api) Connect(userInfo *user.UserInfo, payload *ConnectPayload) error {
	txtid := userInfo.Id
	jid := userInfo.Jid
	token := userInfo.Token
	userId, err := strconv.Atoi(txtid)
	if err != nil {
		return err
	}

	subscribedEvents := make([]string, 0)
	if len(payload.Subscribe) < 1 {
		if !utils.Find(subscribedEvents, "All") {
			subscribedEvents = append(subscribedEvents, "All")
		}
	} else {
		for _, arg := range payload.Subscribe {
			if !utils.Find(whatsapp.MessageTypes, arg) {
				a.log.Warn().Str("Type", arg).Msg("Message type discarded")
				continue
			}
			if !utils.Find(subscribedEvents, arg) {
				subscribedEvents = append(subscribedEvents, arg)
			}
		}
	}

	eventstring := strings.Join(subscribedEvents, ",")
	if err := a.users.SetEvents(userId, eventstring); err != nil {
		a.log.Warn().Msg("Could not set events in users table")
	}

	a.log.Info().Str("events", eventstring).Msg("Setting subscribed events")
	userInfo.Events = eventstring
	a.whatsapp.UpdateCacheUserInfo(token, *userInfo)

	a.log.Info().Str("jid", jid).Msg("Attempt to connect")
	a.whatsapp.NewKillChannel(userId)

	go a.whatsapp.StartClient(userId, jid, token, subscribedEvents)

	if !payload.Immediate {
		a.log.Warn().Msg("Waiting 10 seconds")
		time.Sleep(10000 * time.Millisecond)

		client, err := a.whatsapp.GetClient(userId)
		if err != nil {
			return err
		}

		if client != nil {
			if !client.IsConnected() {
				return whatsapp.ErrFailedToConnect
			}
		} else {
			return whatsapp.ErrFailedToConnect
		}
	}

	return nil
}

func (a *Api) Disconnect(userInfo *user.UserInfo) error {
	txtid := userInfo.Id
	jid := userInfo.Jid
	token := userInfo.Token
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return err
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return err
	}

	if !client.IsConnected() {
		a.log.Warn().Str("jid", jid).Msg("Ignoring disconnect as it was not connected")
		return nil
	}

	if !client.IsLoggedIn() {
		a.log.Warn().Str("jid", jid).Msg("Ignoring disconnect as it was not logged in")
		return nil
	}

	a.log.Info().Str("jid", jid).Msg("Disconnection successfull")
	a.whatsapp.SendKillChannel(userid)

	if err := a.users.SetEvents(userid, ""); err != nil {
		a.log.Warn().Str("userid", txtid).Msg("Could not set events in users table")
	}

	userInfo.Events = ""
	a.whatsapp.UpdateCacheUserInfo(token, *userInfo)

	return nil
}

func (a *Api) GetWebhook(userInfo *user.UserInfo) (*GetWebhookResponse, error) {
	txtid := userInfo.Id
	userId, err := strconv.Atoi(txtid)
	if err != nil {
		return nil, err
	}

	user, err := a.users.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	return &GetWebhookResponse{
		Webhook:   user.Webhook,
		Subscribe: strings.Split(user.Events, ","),
	}, nil
}

func (a *Api) SetWebhook(userInfo *user.UserInfo, payload *SetWebhookPayload) error {
	txtid := userInfo.Id
	token := userInfo.Token
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return err
	}

	var webhook = payload.WebhookURL

	if err := a.users.SetWebhook(userid, webhook); err != nil {
		return err
	}

	userInfo.Webhook = webhook
	a.whatsapp.UpdateCacheUserInfo(token, *userInfo)

	return nil
}

func (a *Api) GetQR(userInfo *user.UserInfo) (string, error) {
	txtid := userInfo.Id
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return "", err
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return "", err
	}

	if !client.IsConnected() {
		return "", whatsapp.ErrNotConnected
	}

	code, err := a.users.GetQRCode(userid)
	if err != nil {
		return "", err
	}

	if client.IsLoggedIn() {
		return "", whatsapp.ErrAlreadyLoggedIn
	}

	return code, nil
}

func (a *Api) Logout(userInfo *user.UserInfo) error {
	txtid := userInfo.Id
	jid := userInfo.Jid
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return err
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return err
	}

	if client.IsLoggedIn() && client.IsConnected() {
		err := client.Logout()
		if err != nil {
			a.log.Error().Str("jid", jid).Msg("Could not perform logout")
			return errors.New("could not perform logout")
		}

		a.log.Info().Str("jid", jid).Msg("Logged out")
		a.whatsapp.SendKillChannel(userid)

		return nil
	} else if client.IsConnected() {
		a.log.Warn().Str("jid", jid).Msg("Ignoring logout as it was not logged in")
		return errors.New("could not disconnect as it was not logged in")
	}

	a.log.Warn().Str("jid", jid).Msg("Ignoring logout as it was not connected")
	return errors.New("could not disconnect as it was not connected")
}

func (a *Api) GetStatus(userInfo *user.UserInfo) (*GetStatusResponse, error) {
	txtid := userInfo.Id
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return nil, err
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return nil, err
	}

	isConnected := client.IsConnected()
	isLoggedIn := client.IsLoggedIn()

	if isConnected && isLoggedIn {
		a.users.SetUserConnected(userid, 1)
	}

	return &GetStatusResponse{
		Connected: isConnected,
		LoggedIn:  isLoggedIn,
	}, nil
}

func (a *Api) SendDocument(userInfo *user.UserInfo, payload *SendDocumentPayload) (whatsmeow.SendResponse, error) {
	txtid := userInfo.Id
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaId,
		payload.ContextInfo.Participant,
		a.whatsapp,
	)
	if err != nil {
		a.log.Error().Msg(fmt.Sprintf("%s", err))
		return whatsmeow.SendResponse{}, err
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	var uploaded whatsmeow.UploadResponse
	var filedata []byte

	if payload.Document[0:29] == "data:application/octet-stream" {
		dataURL, err := dataurl.DecodeString(payload.Document)
		if err != nil {
			return whatsmeow.SendResponse{}, errors.New("could not decode base64 encoded data from payload")
		}

		filedata = dataURL.Data
		uploaded, err = client.Upload(context.Background(), filedata, whatsmeow.MediaDocument)
		if err != nil {
			return whatsmeow.SendResponse{}, fmt.Errorf("failed to upload file: %v", err)
		}
	} else {
		return whatsmeow.SendResponse{}, errors.New("document data should start with \"data:application/octet-stream;base64,\"")
	}

	msg := &waProto.Message{DocumentMessage: &waProto.DocumentMessage{
		Url:           proto.String(uploaded.URL),
		FileName:      &payload.FileName,
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String(http.DetectContentType(filedata)),
		FileEncSha256: uploaded.FileEncSHA256,
		FileSha256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
	}}

	if payload.ContextInfo.StanzaId != nil {
		msg.ExtendedTextMessage.ContextInfo = &waProto.ContextInfo{
			StanzaId:      proto.String(*payload.ContextInfo.StanzaId),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waProto.Message{Conversation: proto.String("")},
		}
	}

	return client.SendMessage(context.Background(), recipient, msg)
}

func (a *Api) SendAudio(userInfo *user.UserInfo, payload *SendAudioPayload) (whatsmeow.SendResponse, error) {
	txtid := userInfo.Id
	userid, _ := strconv.Atoi(txtid)

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaId,
		payload.ContextInfo.Participant,
		a.whatsapp,
	)
	if err != nil {
		a.log.Error().Msg(fmt.Sprintf("%s", err))
		return whatsmeow.SendResponse{}, err
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	var uploaded whatsmeow.UploadResponse
	var filedata []byte

	if payload.Audio[0:14] != "data:audio/ogg" {
		return whatsmeow.SendResponse{}, errors.New("audio data should start with \"data:audio/ogg;base64,\"")
	}

	dataURL, err := dataurl.DecodeString(payload.Audio)
	if err != nil {
		return whatsmeow.SendResponse{}, errors.New("could not decode base64 encoded data from payload")
	}

	filedata = dataURL.Data
	uploaded, err = client.Upload(context.Background(), filedata, whatsmeow.MediaAudio)
	if err != nil {
		return whatsmeow.SendResponse{}, fmt.Errorf("failed to upload file: %v", err)
	}

	ptt := true
	mime := "audio/ogg; codecs=opus"

	msg := &waProto.Message{AudioMessage: &waProto.AudioMessage{
		Url:        proto.String(uploaded.URL),
		DirectPath: proto.String(uploaded.DirectPath),
		MediaKey:   uploaded.MediaKey,
		//Mimetype:      proto.String(http.DetectContentType(filedata)),
		Mimetype:      &mime,
		FileEncSha256: uploaded.FileEncSHA256,
		FileSha256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
		Ptt:           &ptt,
	}}

	if payload.ContextInfo.StanzaId != nil {
		msg.ExtendedTextMessage.ContextInfo = &waProto.ContextInfo{
			StanzaId:      proto.String(*payload.ContextInfo.StanzaId),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waProto.Message{Conversation: proto.String("")},
		}
	}

	return client.SendMessage(context.Background(), recipient, msg)
}

func (a *Api) SendImage(userInfo *user.UserInfo, payload *SendImagePayload) (whatsmeow.SendResponse, error) {

	txtid := userInfo.Id
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaId,
		payload.ContextInfo.Participant,
		a.whatsapp,
	)
	if err != nil {
		a.log.Error().Msg(fmt.Sprintf("%s", err))
		return whatsmeow.SendResponse{}, err
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	var uploaded whatsmeow.UploadResponse
	var filedata []byte

	if payload.Image[0:10] != "data:image" {
		return whatsmeow.SendResponse{}, errors.New("image data should start with \"data:image/png;base64,\"")
	}

	dataURL, err := dataurl.DecodeString(payload.Image)
	if err != nil {
		return whatsmeow.SendResponse{}, errors.New("could not decode base64 encoded data from payload")
	}

	filedata = dataURL.Data
	uploaded, err = client.Upload(context.Background(), filedata, whatsmeow.MediaImage)
	if err != nil {
		return whatsmeow.SendResponse{}, fmt.Errorf("failed to upload file: %v", err)
	}

	msg := &waProto.Message{ImageMessage: &waProto.ImageMessage{
		Caption:       proto.String(payload.Caption),
		Url:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String(http.DetectContentType(filedata)),
		FileEncSha256: uploaded.FileEncSHA256,
		FileSha256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
	}}

	if payload.ContextInfo.StanzaId != nil {
		msg.ExtendedTextMessage.ContextInfo = &waProto.ContextInfo{
			StanzaId:      proto.String(*payload.ContextInfo.StanzaId),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waProto.Message{Conversation: proto.String("")},
		}
	}

	return client.SendMessage(context.Background(), recipient, msg)
}

func (a *Api) SendSticker(userInfo *user.UserInfo, payload *SendStickerPayload) (whatsmeow.SendResponse, error) {
	txtid := userInfo.Id
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaId,
		payload.ContextInfo.Participant,
		a.whatsapp,
	)
	if err != nil {
		a.log.Error().Msg(fmt.Sprintf("%s", err))
		return whatsmeow.SendResponse{}, err
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	var uploaded whatsmeow.UploadResponse
	var filedata []byte

	if payload.Sticker[0:4] != "data" {
		return whatsmeow.SendResponse{}, errors.New("data should start with \"data:mime/type;base64,\"")
	}

	dataURL, err := dataurl.DecodeString(payload.Sticker)
	if err != nil {
		return whatsmeow.SendResponse{}, errors.New("could not decode base64 encoded data from payload")
	}

	filedata = dataURL.Data
	uploaded, err = client.Upload(context.Background(), filedata, whatsmeow.MediaImage)
	if err != nil {
		return whatsmeow.SendResponse{}, fmt.Errorf("failed to upload file: %v", err)
	}

	msg := &waProto.Message{StickerMessage: &waProto.StickerMessage{
		Url:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String(http.DetectContentType(filedata)),
		FileEncSha256: uploaded.FileEncSHA256,
		FileSha256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
		PngThumbnail:  payload.PngThumbnail,
	}}

	if payload.ContextInfo.StanzaId != nil {
		msg.ExtendedTextMessage.ContextInfo = &waProto.ContextInfo{
			StanzaId:      proto.String(*payload.ContextInfo.StanzaId),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waProto.Message{Conversation: proto.String("")},
		}
	}

	return client.SendMessage(context.Background(), recipient, msg)
}

func (a *Api) SendVideo(userInfo *user.UserInfo, payload *SendVideoPayload) (whatsmeow.SendResponse, error) {
	txtid := userInfo.Id
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaId,
		payload.ContextInfo.Participant,
		a.whatsapp,
	)
	if err != nil {
		a.log.Error().Msg(fmt.Sprintf("%s", err))
		return whatsmeow.SendResponse{}, err
	}

	var uploaded whatsmeow.UploadResponse
	var filedata []byte
	if payload.Video[0:4] != "data" {
		return whatsmeow.SendResponse{}, errors.New("data should start with \"data:mime/type;base64,\"")
	}

	dataURL, err := dataurl.DecodeString(payload.Video)
	if err != nil {
		return whatsmeow.SendResponse{}, errors.New("could not decode base64 encoded data from payload")
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	filedata = dataURL.Data
	uploaded, err = client.Upload(context.Background(), filedata, whatsmeow.MediaVideo)
	if err != nil {
		return whatsmeow.SendResponse{}, fmt.Errorf("failed to upload file: %v", err)
	}

	msg := &waProto.Message{VideoMessage: &waProto.VideoMessage{
		Caption:       proto.String(payload.Caption),
		Url:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String(http.DetectContentType(filedata)),
		FileEncSha256: uploaded.FileEncSHA256,
		FileSha256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(filedata))),
		JpegThumbnail: payload.JpegThumbnail,
	}}

	if payload.ContextInfo.StanzaId != nil {
		msg.ExtendedTextMessage.ContextInfo = &waProto.ContextInfo{
			StanzaId:      proto.String(*payload.ContextInfo.StanzaId),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waProto.Message{Conversation: proto.String("")},
		}
	}

	return client.SendMessage(context.Background(), recipient, msg)
}

func (a *Api) SendContact(userInfo *user.UserInfo, payload *SendContactPayload) (whatsmeow.SendResponse, error) {
	txtid := userInfo.Id
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaId,
		payload.ContextInfo.Participant,
		a.whatsapp,
	)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	msg := &waProto.Message{ContactMessage: &waProto.ContactMessage{
		DisplayName: &payload.Name,
		Vcard:       &payload.Vcard,
	}}

	if payload.ContextInfo.StanzaId != nil {
		msg.ExtendedTextMessage.ContextInfo = &waProto.ContextInfo{
			StanzaId:      proto.String(*payload.ContextInfo.StanzaId),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waProto.Message{Conversation: proto.String("")},
		}
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	return client.SendMessage(context.Background(), recipient, msg)
}

func (a *Api) SendLocation(userInfo *user.UserInfo, payload *SendLocationPayload) (whatsmeow.SendResponse, error) {
	txtid := userInfo.Id
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaId,
		payload.ContextInfo.Participant,
		a.whatsapp,
	)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	msg := &waProto.Message{LocationMessage: &waProto.LocationMessage{
		DegreesLatitude:  &payload.Latitude,
		DegreesLongitude: &payload.Longitude,
		Name:             &payload.Name,
	}}

	if payload.ContextInfo.StanzaId != nil {
		msg.ExtendedTextMessage.ContextInfo = &waProto.ContextInfo{
			StanzaId:      proto.String(*payload.ContextInfo.StanzaId),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waProto.Message{Conversation: proto.String("")},
		}
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	return client.SendMessage(context.Background(), recipient, msg)
}
func (a *Api) SendButton(userInfo *user.UserInfo, payload *SendButtonTextPayload) (whatsmeow.SendResponse, error) {
	txtid := userInfo.Id
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	recipient, ok := a.whatsapp.ParseJID(payload.Phone)
	if !ok {
		return whatsmeow.SendResponse{}, err
	}

	var buttons []*waProto.ButtonsMessage_Button
	for _, item := range payload.Buttons {
		buttons = append(buttons, &waProto.ButtonsMessage_Button{
			ButtonId: proto.String(item.ButtonId),
			ButtonText: &waProto.ButtonsMessage_Button_ButtonText{
				DisplayText: proto.String(item.ButtonText),
			},
			Type:           waProto.ButtonsMessage_Button_RESPONSE.Enum(),
			NativeFlowInfo: &waProto.ButtonsMessage_Button_NativeFlowInfo{},
		})
	}

	msg2 := &waProto.ButtonsMessage{
		ContentText: proto.String(payload.Title),
		HeaderType:  waProto.ButtonsMessage_EMPTY.Enum(),
		Buttons:     buttons,
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	return client.SendMessage(
		context.Background(),
		recipient,
		&waProto.Message{
			ViewOnceMessage: &waProto.FutureProofMessage{
				Message: &waProto.Message{
					ButtonsMessage: msg2,
				},
			},
		},
	)
}
func (a *Api) SendList(userInfo *user.UserInfo, payload *SendListPayload) (whatsmeow.SendResponse, error) {
	txtid := userInfo.Id
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	recipient, ok := a.whatsapp.ParseJID(payload.Phone)
	if !ok {
		return whatsmeow.SendResponse{}, err
	}

	var sections []*waProto.ListMessage_Section
	for _, item := range payload.Sections {
		var rows []*waProto.ListMessage_Row
		id := 1
		for _, row := range item.Rows {
			var idtext string
			if row.RowId == "" {
				idtext = strconv.Itoa(id)
			} else {
				idtext = row.RowId
			}
			rows = append(rows, &waProto.ListMessage_Row{
				RowId:       proto.String(idtext),
				Title:       proto.String(row.Title),
				Description: proto.String(row.Description),
			})
		}

		sections = append(sections, &waProto.ListMessage_Section{
			Title: proto.String(item.Title),
			Rows:  rows,
		})
	}

	msg1 := &waProto.ListMessage{
		Title:       proto.String(payload.Title),
		Description: proto.String(payload.Description),
		ButtonText:  proto.String(payload.ButtonText),
		ListType:    waProto.ListMessage_SINGLE_SELECT.Enum(),
		Sections:    sections,
		FooterText:  proto.String(payload.FooterText),
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	return client.SendMessage(
		context.Background(),
		recipient,
		&waProto.Message{
			ViewOnceMessage: &waProto.FutureProofMessage{
				Message: &waProto.Message{
					ListMessage: msg1,
				},
			},
		},
	)
}

func (a *Api) SendText(userInfo *user.UserInfo, payload *SendTextPayload) (whatsmeow.SendResponse, error) {
	txtid := userInfo.Id
	userId, err := strconv.Atoi(txtid)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	msg := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: &payload.Body,
		},
	}

	recipient, err := validateMessageFields(
		payload.Phone,
		payload.ContextInfo.StanzaId,
		payload.ContextInfo.Participant,
		a.whatsapp,
	)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	if payload.ContextInfo.StanzaId != nil {
		msg.ExtendedTextMessage.ContextInfo = &waProto.ContextInfo{
			StanzaId:      proto.String(*payload.ContextInfo.StanzaId),
			Participant:   proto.String(*payload.ContextInfo.Participant),
			QuotedMessage: &waProto.Message{Conversation: proto.String("")},
		}
	}

	client, err := a.whatsapp.GetClient(userId)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	return client.SendMessage(context.Background(), recipient, msg)
}

func (a *Api) CheckUser(userInfo *user.UserInfo, payload *CheckUserPayload) (*UserCollection, error) {
	userid, err := strconv.Atoi(userInfo.Id)
	if err != nil {
		return nil, err
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return nil, err
	}

	resp, err := client.IsOnWhatsApp(payload.Phone)
	if err != nil {
		return nil, err
	}

	uc := new(UserCollection)
	for _, item := range resp {
		if item.VerifiedName != nil {
			var msg = User{Query: item.Query, IsInWhatsapp: item.IsIn, JID: item.JID.String(), VerifiedName: item.VerifiedName.Details.GetVerifiedName()}
			uc.Users = append(uc.Users, msg)
		} else {
			var msg = User{Query: item.Query, IsInWhatsapp: item.IsIn, JID: item.JID.String(), VerifiedName: ""}
			uc.Users = append(uc.Users, msg)
		}
	}

	return uc, nil
}

func (a *Api) GetUser(userInfo *user.UserInfo, payload CheckUserPayload) (*UserInfoCollection, error) {
	txtid := userInfo.Id
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return nil, err
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return nil, err
	}

	var jids []types.JID
	for _, arg := range payload.Phone {
		jid, ok := a.whatsapp.ParseJID(arg)
		if !ok {
			return nil, whatsapp.ErrInvalidPhoneNumber
		}
		jids = append(jids, jid)
	}

	resp, err := client.GetUserInfo(jids)
	if err != nil {
		msg := fmt.Sprintf("Failed to get user info: %v", err)
		a.log.Error().Msg(msg)
		return nil, errors.New(msg)
	}

	uc := new(UserInfoCollection)
	uc.Users = make(map[types.JID]types.UserInfo)

	for jid, info := range resp {
		uc.Users[jid] = info
	}

	return uc, nil
}

func (a *Api) GetAvatar(userInfo *user.UserInfo, getAvatar *GetAvatarPayload) (*types.ProfilePictureInfo, error) {
	txtid := userInfo.Id
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return nil, err
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return nil, err
	}

	jid, ok := a.whatsapp.ParseJID(getAvatar.Phone)
	if !ok {
		return nil, whatsapp.ErrInvalidPhoneNumber
	}

	var pic *types.ProfilePictureInfo

	existingID := ""
	pic, err = client.GetProfilePictureInfo(
		jid,
		&whatsmeow.GetProfilePictureParams{
			Preview:    getAvatar.Preview,
			ExistingID: existingID,
		},
	)
	if err != nil {
		msg := fmt.Sprintf("Failed to get avatar: %v", err)
		a.log.Error().Msg(msg)
		return nil, errors.New(msg)
	}

	if pic == nil {
		return nil, whatsapp.ErrNoAvatarFound
	}

	a.log.Info().Str("id", pic.ID).Str("url", pic.URL).Msg("Got avatar")

	return pic, nil
}

func (a *Api) GetContacts(userInfo *user.UserInfo) ([]byte, error) {
	txtid := userInfo.Id
	userid, err := strconv.Atoi(txtid)
	if err != nil {
		return nil, err
	}

	client, err := a.whatsapp.GetClient(userid)
	if err != nil {
		return nil, err
	}

	result, err := client.Store.Contacts.GetAllContacts()
	if err != nil {
		return nil, err
	}

	json, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return json, nil
}

func (a *Api) SendChatPresence(userInfo *user.UserInfo, payload *ChatPresencePayload) error {
	txtid := userInfo.Id
	userId, err := strconv.Atoi(txtid)
	if err != nil {
		return err
	}

	client, err := a.whatsapp.GetClient(userId)
	if err != nil {
		return err
	}

	jid, ok := a.whatsapp.ParseJID(payload.Phone)
	if !ok {
		return whatsapp.ErrInvalidPhoneNumber
	}

	return client.SendChatPresence(
		jid,
		types.ChatPresence(payload.State),
		types.ChatPresenceMedia(payload.Media),
	)
}
