package whatsapp

import (
	"strconv"
	"strings"

	"github.com/nugrhrizki/buzz/pkg/whatsapp/user"
	"github.com/patrickmn/go-cache"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

func (w *Whatsapp) UserToUserInfo(u *user.User) user.UserInfo {
	return user.UserInfo{
		Id:      strconv.Itoa(u.Id),
		Jid:     u.Jid,
		Webhook: u.Webhook,
		Token:   u.Token,
		Events:  u.Events,
	}
}

func (w *Whatsapp) GetCacheUserInfo(token string) (user.UserInfo, bool) {
	if x, found := w.userInfoCache.Get(token); found {
		return x.(user.UserInfo), true
	}
	return user.UserInfo{}, false
}

func (w *Whatsapp) UpdateCacheUserInfo(token string, value user.UserInfo) {
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
	phonenumber = strings.Split(phonenumber, ":")[0]
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
		return nil, ErrNoSession
	}

	if client == nil {
		return nil, ErrNotConnected
	}

	return client, nil
}
