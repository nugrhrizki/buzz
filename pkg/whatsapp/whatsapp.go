package whatsapp

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	_ "modernc.org/sqlite"

	"github.com/nugrhrizki/buzz/pkg/env"
	"github.com/nugrhrizki/buzz/pkg/utils"
	"github.com/nugrhrizki/buzz/pkg/whatsapp/user"
)

type Whatsapp struct {
	clientStore map[int]*whatsmeow.Client
	clientHttp  map[int]*resty.Client
	killchannel map[int](chan bool)

	container     *sqlstore.Container
	userInfoCache *cache.Cache
	log           *zerolog.Logger

	users *user.Repository
}

var MessageTypes = []string{
	"Message",
	"ReadReceipt",
	"Presence",
	"HistorySync",
	"ChatPresence",
	"All",
}

func New(
	users *user.Repository,
	log *zerolog.Logger,
	env *env.Env,
) *Whatsapp {
	container, err := sqlstore.New(
		env.DB_DRIVER,
		env.DB_DSN,
		nil,
	)

	if err != nil {
		panic(err)
	}

	return &Whatsapp{
		clientStore: make(map[int]*whatsmeow.Client),
		clientHttp:  make(map[int]*resty.Client),
		killchannel: make(map[int](chan bool)),

		container:     container,
		userInfoCache: cache.New(5*time.Minute, 10*time.Minute),
		log:           log,

		users: users,
	}
}

// Connects to Whatsapp Websocket on server startup if last state was connected
func (w *Whatsapp) ConnectOnStartup() {
	users, err := w.users.GetConnectedUser()
	if err != nil {
		w.log.Error().Err(err).Msg("DB Problem")
		return
	}

	for _, u := range users {
		w.log.Info().Str("token", u.Token).Msg("Connect to Whatsapp on startup")

		userInfo := user.UserInfo{
			Id:      strconv.Itoa(u.Id),
			Jid:     u.Jid,
			Webhook: u.Webhook,
			Token:   u.Token,
			Events:  u.Events,
		}

		w.userInfoCache.Set(u.Token, userInfo, cache.NoExpiration)
		// Gets and set subscription to webhook events
		eventarray := strings.Split(u.Events, ",")

		var subscribedEvents []string
		if len(eventarray) < 1 {
			if !utils.Find(subscribedEvents, "All") {
				subscribedEvents = append(subscribedEvents, "All")
			}
		} else {
			for _, arg := range eventarray {
				if !utils.Find(MessageTypes, arg) {
					w.log.Warn().Str("Type", arg).Msg("Message type discarded")
					continue
				}
				if !utils.Find(subscribedEvents, arg) {
					subscribedEvents = append(subscribedEvents, arg)
				}
			}
		}

		eventstring := strings.Join(subscribedEvents, ",")
		w.log.Info().Str("events", eventstring).Str("jid", u.Jid).Msg("Attempt to connect")
		w.killchannel[u.Id] = make(chan bool)
		go w.StartClient(u.Id, u.Jid, u.Token, subscribedEvents)
	}
}

func (w *Whatsapp) StartClient(userID int, textjid string, token string, subscriptions []string) {
	w.log.Info().
		Str("userid", strconv.Itoa(userID)).
		Str("jid", textjid).
		Msg("Starting websocket connection to Whatsapp")

	var deviceStore *store.Device
	var err error

	if w.clientStore[userID] != nil {
		isConnected := w.clientStore[userID].IsConnected()
		if isConnected {
			return
		}
	}

	if textjid != "" {
		jid, _ := w.ParseJID(textjid)
		deviceStore, err = w.container.GetDevice(jid)
		if err != nil {
			panic(err)
		}
	} else {
		w.log.Warn().Msg("No jid found. Creating new device")
		deviceStore = w.container.NewDevice()
	}

	if deviceStore == nil {
		w.log.Warn().Msg("No store found. Creating new one")
		deviceStore = w.container.NewDevice()
	}

	osName := "WADAK"
	store.DeviceProps.PlatformType = waProto.DeviceProps_CHROME.Enum()
	store.DeviceProps.Os = &osName

	wclient := whatsmeow.NewClient(deviceStore, nil)

	w.clientStore[userID] = wclient
	client := NewClient(
		wclient,
		1,
		userID,
		token,
		subscriptions,
		w,
	)
	client.SetEventHandlerID(client.WAClient.AddEventHandler(client.EventHandler))
	w.clientHttp[userID] = resty.New()
	w.clientHttp[userID].SetRedirectPolicy(resty.FlexibleRedirectPolicy(15))
	w.clientHttp[userID].SetTimeout(5 * time.Second)
	w.clientHttp[userID].SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	if wclient.Store.ID == nil {
		// No ID stored, new login

		qrChan, err := wclient.GetQRChannel(context.Background())
		if err != nil {
			// This error means that we're already logged in, so ignore it.
			if !errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
				w.log.Error().Err(err).Msg("Failed to get QR channel")
			}
		} else {
			err = wclient.Connect()
			if err != nil {
				panic(err)
			}
			for evt := range qrChan {
				switch evt.Event {
				case "code":
					// Display QR code in terminal (useful for testing/developing)
					// Store encoded/embeded base64 QR on database for retrieval with the /qr endpoint
					image, _ := qrcode.Encode(evt.Code, qrcode.Medium, 256)
					base64qrcode := "data:image/png;base64," + base64.StdEncoding.EncodeToString(image)
					err := w.users.SetQRCode(userID, base64qrcode)
					if err != nil {
						w.log.Error().Err(err).Msg("Failed to set QR code")
					}
				case "timeout":
					// Clear QR code from DB on timeout
					err = w.users.SetQRCode(userID, "")
					if err != nil {
						w.log.Error().Err(err).Msg("Failed to clear QR code")
					}
					w.log.Warn().Msg("QR timeout killing channel")
					delete(w.clientStore, userID)
					w.killchannel[userID] <- true
				case "success":
					w.log.Info().Msg("QR pairing ok!")
					// Clear QR code after pairing
					err = w.users.SetQRCode(userID, "")
					if err != nil {
						w.log.Error().Err(err).Msg("Failed to clear QR code")
					}
				default:
					w.log.Info().Str("event", evt.Event).Msg("Login event")
				}
			}
		}

	} else {
		// Already logged in, just connect
		w.log.Info().Msg("Already logged in, just connect")
		err = wclient.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Keep connected client live until disconnected/killed
	for {
		select {
		case <-w.killchannel[userID]:
			w.log.Info().Str("userid", strconv.Itoa(userID)).Msg("Received kill signal")
			wclient.Disconnect()
			delete(w.clientStore, userID)
			err = w.users.SetUserConnected(userID, 0)
			if err != nil {
				w.log.Error().Err(err).Msg("Failed to set user disconnected")
			}
			return
		default:
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
