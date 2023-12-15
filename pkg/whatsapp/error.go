package whatsapp

import "errors"

var (
	ErrNoAvatarFound          = errors.New("no avatar found")
	ErrNotConnected           = errors.New("not connected")
	ErrAlreadyLoggedIn        = errors.New("already logged in")
	ErrNoSession              = errors.New("no session")
	ErrClientAlreadyConnected = errors.New("client already connected")
	ErrClientNotFound         = errors.New("client not found")
	ErrFailedToConnect        = errors.New("failed to connect")
	ErrInvalidPhoneNumber     = errors.New("invalid phone number")
	ErrEmptyBody              = errors.New("body cannot be empty")
	ErrMissingStanzaId        = errors.New("missing stanza id in contextinfo")
	ErrMissingParticipant     = errors.New("missing participant in contextinfo")
)
