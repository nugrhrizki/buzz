package whatsapp

import (
	"strings"

	"github.com/nugrhrizki/buzz/pkg/whatsapp/whatsapp_user"
	"github.com/patrickmn/go-cache"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

// Update entry in User map
func (w *Whatsapp) UpdateUserInfo(values interface{}, field string, value string) interface{} {
	w.log.Debug().Str("field", field).Str("value", value).Msg("User info updated")
	// values.(Values).m[field] = value
	return values
}

func (w *Whatsapp) UpdateCacheUserInfo(token string, value whatsapp_user.WhatsappUserInfo) {
	w.userInfoCache.Set(token, value, cache.NoExpiration)
}

// parseJID parses a JID string and returns a JID struct
func (w *Whatsapp) ParseJID(arg string) (types.JID, bool) {
	if arg == "" {
		return types.NewJID("", types.DefaultUserServer), false
	}
	if arg[0] == '+' {
		arg = arg[1:]
	}

	// Basic only digit check for recipient phone number, we want to remove @server and .session
	phonenumber := ""
	phonenumber = strings.Split(arg, "@")[0]
	phonenumber = strings.Split(phonenumber, ".")[0]
	b := true
	for _, c := range phonenumber {
		if c < '0' || c > '9' {
			b = false
			break
		}
	}
	if !b {
		w.log.Warn().Msg("Bad jid format, return empty")
		recipient, _ := types.ParseJID("")
		return recipient, false
	}

	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			w.log.Error().Err(err).Str("jid", arg).Msg("Invalid jid")
			return recipient, false
		} else if recipient.User == "" {
			w.log.Error().Err(err).Str("jid", arg).Msg("Invalid jid. No server specified")
			return recipient, false
		}
		return recipient, true
	}
}

func (w *Whatsapp) GetClient(userId int) (*whatsmeow.Client, error) {
	client, ok := w.clientStore[userId]

	if !ok {
		return nil, ErrClientNotFound
	}

	if client == nil {
		return nil, ErrNotConnected
	}

	return client, nil
}
