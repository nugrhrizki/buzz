package api

import (
	"github.com/nugrhrizki/buzz/pkg/whatsapp"
	"go.mau.fi/whatsmeow/types"
)

func validateMessageFields(
	phone string,
	stanzaid *string,
	participant *string,
	wa *whatsapp.Whatsapp,
) (types.JID, error) {
	recipient, ok := wa.ParseJID(phone)
	if !ok {
		return types.NewJID("", types.DefaultUserServer), whatsapp.ErrInvalidPhoneNumber
	}

	if stanzaid != nil {
		if participant == nil {
			return types.NewJID(
				"",
				types.DefaultUserServer,
			), whatsapp.ErrMissingParticipant
		}
	}

	if participant != nil {
		if stanzaid == nil {
			return types.NewJID(
				"",
				types.DefaultUserServer,
			), whatsapp.ErrMissingStanzaId
		}
	}

	return recipient, nil
}
