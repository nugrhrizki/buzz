package whatsapp

import (
	"encoding/json"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/nugrhrizki/buzz/pkg/utils"
	"github.com/nugrhrizki/buzz/pkg/whatsapp/user"
	"github.com/patrickmn/go-cache"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

var historySyncID int32

type Client struct {
	WAClient       *whatsmeow.Client
	eventHandlerID uint32
	userID         int
	token          string
	subscriptions  []string
	whatsapp       *Whatsapp
}

func NewClient(
	waClient *whatsmeow.Client,
	eventHandlerID uint32,
	userID int,
	token string,
	subscriptions []string,

	whatsapp *Whatsapp,
) Client {
	return Client{
		WAClient:       waClient,
		eventHandlerID: eventHandlerID,
		userID:         userID,
		token:          token,
		subscriptions:  subscriptions,
		whatsapp:       whatsapp,
	}
}

func (c *Client) SetClient(waClient *whatsmeow.Client) {
	c.WAClient = waClient
}

func (c *Client) SetEventHandlerID(eventHandlerID uint32) {
	c.eventHandlerID = eventHandlerID
}

func (c *Client) SetUserID(userID int) {
	c.userID = userID
}

func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) SetSubscriptions(subscriptions []string) {
	c.subscriptions = subscriptions
}

func (c *Client) EventHandler(rawEvt interface{}) {
	txtid := strconv.Itoa(c.userID)
	postmap := make(map[string]interface{})
	postmap["event"] = rawEvt
	dowebhook := 0
	path := ""

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	switch evt := rawEvt.(type) {
	case *events.AppStateSyncComplete:
		if len(c.WAClient.Store.PushName) > 0 && evt.Name == appstate.WAPatchCriticalBlock {
			err := c.WAClient.SendPresence(types.PresenceAvailable)
			if err != nil {
				c.whatsapp.log.Warn().Err(err).Msg("Failed to send available presence")
			} else {
				c.whatsapp.log.Info().Msg("Marked self as available")
			}
		}
	case *events.Connected, *events.PushNameSetting:
		if len(c.WAClient.Store.PushName) == 0 {
			return
		}
		// Send presence available when connecting and when the pushname is changed.
		// This makes sure that outgoing messages always have the right pushname.
		err := c.WAClient.SendPresence(types.PresenceAvailable)
		if err != nil {
			c.whatsapp.log.Warn().Err(err).Msg("Failed to send available presence")
		} else {
			c.whatsapp.log.Info().Msg("Marked self as available")
		}
		c.whatsapp.log.Info().Msg("Setting up status connection")
		err = c.whatsapp.users.SetUserConnected(c.userID, 1)
		if err != nil {
			c.whatsapp.log.Error().Err(err).Msg("Failed to set user connected")
			return
		}
	case *events.PairSuccess:
		c.whatsapp.log.Info().Str("userid", strconv.Itoa(c.userID)).Str("token", c.token).Str("ID", evt.ID.String()).Str("BusinessName", evt.BusinessName).Str("Platform", evt.Platform).Msg("QR Pair Success")
		jid := evt.ID
		err = c.whatsapp.users.SetUserJid(c.userID, jid)
		if err != nil {
			c.whatsapp.log.Error().Err(err).Msg("Failed to set user jid")
			return
		}

		userInfo, found := c.whatsapp.userInfoCache.Get(c.token)
		if !found {
			c.whatsapp.log.Warn().Msg("No user info cached on pairing?")
		} else {
			txtid := userInfo.(user.UserInfo).Id
			token := userInfo.(user.UserInfo).Token
			newUserInfo := user.UserInfo{
				Id:      txtid,
				Jid:     jid.String(),
				Webhook: userInfo.(user.UserInfo).Webhook,
				Token:   token,
				Events:  userInfo.(user.UserInfo).Events,
			}
			c.whatsapp.userInfoCache.Set(token, newUserInfo, cache.NoExpiration)
			c.whatsapp.log.Info().Str("jid", jid.String()).Str("userid", txtid).Str("token", token).Msg("User information set")
		}
	case *events.StreamReplaced:
		c.whatsapp.log.Info().Msg("Received StreamReplaced event")
		return
	case *events.Message:
		postmap["type"] = "Message"
		dowebhook = 1
		metaParts := []string{fmt.Sprintf("pushname: %s", evt.Info.PushName), fmt.Sprintf("timestamp: %s", evt.Info.Timestamp)}
		if evt.Info.Type != "" {
			metaParts = append(metaParts, fmt.Sprintf("type: %s", evt.Info.Type))
		}
		if evt.Info.Category != "" {
			metaParts = append(metaParts, fmt.Sprintf("category: %s", evt.Info.Category))
		}
		if evt.IsViewOnce {
			metaParts = append(metaParts, "view once")
		}
		if evt.IsViewOnce {
			metaParts = append(metaParts, "ephemeral")
		}

		c.whatsapp.log.Info().Str("id", evt.Info.ID).Str("source", evt.Info.SourceString()).Str("parts", strings.Join(metaParts, ", ")).Msg("Message Received")

		// try to get Image if any
		img := evt.Message.GetImageMessage()
		if img != nil {

			// check/creates user directory for files
			userDirectory := fmt.Sprintf("%s/files/user_%s", exPath, txtid)
			_, err := os.Stat(userDirectory)
			if os.IsNotExist(err) {
				errDir := os.MkdirAll(userDirectory, 0751)
				if errDir != nil {
					c.whatsapp.log.Error().Err(errDir).Msg("Could not create user directory")
					return
				}
			}

			data, err := c.WAClient.Download(img)
			if err != nil {
				c.whatsapp.log.Error().Err(err).Msg("Failed to download image")
				return
			}
			exts, _ := mime.ExtensionsByType(img.GetMimetype())
			path = fmt.Sprintf("%s/%s%s", userDirectory, evt.Info.ID, exts[0])
			err = os.WriteFile(path, data, 0600)
			if err != nil {
				c.whatsapp.log.Error().Err(err).Msg("Failed to save image")
				return
			}
			c.whatsapp.log.Info().Str("path", path).Msg("Image saved")
		}

		// try to get Audio if any
		audio := evt.Message.GetAudioMessage()
		if audio != nil {

			// check/creates user directory for files
			userDirectory := fmt.Sprintf("%s/files/user_%s", exPath, txtid)
			_, err := os.Stat(userDirectory)
			if os.IsNotExist(err) {
				errDir := os.MkdirAll(userDirectory, 0751)
				if errDir != nil {
					c.whatsapp.log.Error().Err(errDir).Msg("Could not create user directory")
					return
				}
			}

			data, err := c.WAClient.Download(audio)
			if err != nil {
				c.whatsapp.log.Error().Err(err).Msg("Failed to download audio")
				return
			}
			exts, _ := mime.ExtensionsByType(audio.GetMimetype())
			path = fmt.Sprintf("%s/%s%s", userDirectory, evt.Info.ID, exts[0])
			err = os.WriteFile(path, data, 0600)
			if err != nil {
				c.whatsapp.log.Error().Err(err).Msg("Failed to save audio")
				return
			}
			c.whatsapp.log.Info().Str("path", path).Msg("Audio saved")
		}

		// try to get Document if any
		document := evt.Message.GetDocumentMessage()
		if document != nil {

			// check/creates user directory for files
			userDirectory := fmt.Sprintf("%s/files/user_%s", exPath, txtid)
			_, err := os.Stat(userDirectory)
			if os.IsNotExist(err) {
				errDir := os.MkdirAll(userDirectory, 0751)
				if errDir != nil {
					c.whatsapp.log.Error().Err(errDir).Msg("Could not create user directory")
					return
				}
			}

			data, err := c.WAClient.Download(document)
			if err != nil {
				c.whatsapp.log.Error().Err(err).Msg("Failed to download document")
				return
			}
			extension := ""
			exts, err := mime.ExtensionsByType(document.GetMimetype())
			if err != nil {
				extension = exts[0]
			} else {
				filename := document.FileName
				extension = filepath.Ext(*filename)
			}
			path = fmt.Sprintf("%s/%s%s", userDirectory, evt.Info.ID, extension)
			err = os.WriteFile(path, data, 0600)
			if err != nil {
				c.whatsapp.log.Error().Err(err).Msg("Failed to save document")
				return
			}
			c.whatsapp.log.Info().Str("path", path).Msg("Document saved")
		}
	case *events.Receipt:
		postmap["type"] = "ReadReceipt"
		dowebhook = 1
		switch evt.Type {
		case types.ReceiptTypeRead, types.ReceiptTypeReadSelf:
			c.whatsapp.log.Info().Strs("id", evt.MessageIDs).Str("source", evt.SourceString()).Str("timestamp", fmt.Sprintf("%v", evt.Timestamp)).Msg("Message was read")
			if evt.Type == types.ReceiptTypeRead {
				postmap["state"] = "Read"
			} else {
				postmap["state"] = "ReadSelf"
			}
		case types.ReceiptTypeDelivered:
			postmap["state"] = "Delivered"
			c.whatsapp.log.Info().Str("id", evt.MessageIDs[0]).Str("source", evt.SourceString()).Str("timestamp", fmt.Sprintf("%v", evt.Timestamp)).Msg("Message delivered")
		default:
			// Discard webhooks for inactive or other delivery types
			return
		}
	case *events.Presence:
		postmap["type"] = "Presence"
		dowebhook = 1
		if evt.Unavailable {
			postmap["state"] = "offline"
			if evt.LastSeen.IsZero() {
				c.whatsapp.log.Info().Str("from", evt.From.String()).Msg("User is now offline")
			} else {
				c.whatsapp.log.Info().Str("from", evt.From.String()).Str("lastSeen", fmt.Sprintf("%v", evt.LastSeen)).Msg("User is now offline")
			}
		} else {
			postmap["state"] = "online"
			c.whatsapp.log.Info().Str("from", evt.From.String()).Msg("User is now online")
		}
	case *events.HistorySync:
		postmap["type"] = "HistorySync"
		dowebhook = 1

		// check/creates user directory for files
		userDirectory := fmt.Sprintf("%s/files/user_%s", exPath, txtid)
		_, err := os.Stat(userDirectory)
		if os.IsNotExist(err) {
			errDir := os.MkdirAll(userDirectory, 0751)
			if errDir != nil {
				c.whatsapp.log.Error().Err(errDir).Msg("Could not create user directory")
				return
			}
		}

		id := atomic.AddInt32(&historySyncID, 1)
		fileName := fmt.Sprintf("%s/history-%d.json", userDirectory, id)
		file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			c.whatsapp.log.Error().Err(err).Msg("Failed to open file to write history sync")
			return
		}
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		err = enc.Encode(evt.Data)
		if err != nil {
			c.whatsapp.log.Error().Err(err).Msg("Failed to write history sync")
			return
		}
		c.whatsapp.log.Info().Str("filename", fileName).Msg("Wrote history sync")
		_ = file.Close()
	case *events.AppState:
		c.whatsapp.log.Info().Str("index", fmt.Sprintf("%+v", evt.Index)).Str("actionValue", fmt.Sprintf("%+v", evt.SyncActionValue)).Msg("App state event received")
	case *events.LoggedOut:
		c.whatsapp.log.Info().Str("reason", evt.Reason.String()).Msg("Logged out")
		c.whatsapp.killchannel[c.userID] <- true
		err := c.whatsapp.users.SetUserConnected(c.userID, 0)
		if err != nil {
			c.whatsapp.log.Error().Err(err).Msg("Failed to set user disconnected")
			return
		}
	case *events.ChatPresence:
		postmap["type"] = "ChatPresence"
		dowebhook = 1
		c.whatsapp.log.Info().Str("state", string(evt.State)).Str("media", string(evt.Media)).Str("chat", evt.MessageSource.Chat.String()).Str("sender", evt.MessageSource.Sender.String()).Msg("Chat Presence received")
	case *events.CallOffer:
		c.whatsapp.log.Info().Str("event", fmt.Sprintf("%+v", evt)).Msg("Got call offer")
	case *events.CallAccept:
		c.whatsapp.log.Info().Str("event", fmt.Sprintf("%+v", evt)).Msg("Got call accept")
	case *events.CallTerminate:
		c.whatsapp.log.Info().Str("event", fmt.Sprintf("%+v", evt)).Msg("Got call terminate")
	case *events.CallOfferNotice:
		c.whatsapp.log.Info().Str("event", fmt.Sprintf("%+v", evt)).Msg("Got call offer notice")
	case *events.CallRelayLatency:
		c.whatsapp.log.Info().Str("event", fmt.Sprintf("%+v", evt)).Msg("Got call relay latency")
	default:
		c.whatsapp.log.Warn().Str("event", fmt.Sprintf("%+v", evt)).Msg("Unhandled event")
	}

	if dowebhook == 1 {
		// call webhook
		webhookurl := ""
		userInfo, found := c.whatsapp.userInfoCache.Get(c.token)
		if !found {
			c.whatsapp.log.Warn().
				Str("token", c.token).
				Msg("Could not call webhook as there is no user for this token")
		} else {
			webhookurl = userInfo.(user.UserInfo).Webhook
		}

		if !utils.Find(c.subscriptions, postmap["type"].(string)) &&
			!utils.Find(c.subscriptions, "All") {
			c.whatsapp.log.Warn().
				Str("type", postmap["type"].(string)).
				Msg("Skipping webhook. Not subscribed for this type")
			return
		}

		if webhookurl != "" {
			c.whatsapp.log.Info().Str("url", webhookurl).Msg("Calling webhook")
			values, _ := json.Marshal(postmap)
			if path == "" {
				data := make(map[string]string)
				data["jsonData"] = string(values)
				data["token"] = c.token
				go c.whatsapp.CallHook(webhookurl, data, c.userID)
			} else {
				data := make(map[string]string)
				data["jsonData"] = string(values)
				data["token"] = c.token
				go c.whatsapp.CallHookFile(webhookurl, data, c.userID, path)
			}
		} else {
			c.whatsapp.log.Warn().Str("userid", strconv.Itoa(c.userID)).Msg("No webhook set for user")
		}
	}
}
