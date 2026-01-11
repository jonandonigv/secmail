package email

import "time"

type Message struct {
	ID                   uint
	ConversationID       uint
	SenderID             uint
	Recipients           []uint // List of recipient user IDs
	EncryptedBody        []byte
	EncryptedSessionKeys []byte
	EncryptedAttachments []byte
	Metadata             string // JSON string for additional data
	Status               string // e.g., "sent", "delivered", "read"
	CreatedAt            time.Time
	SentAt               time.Time
}
