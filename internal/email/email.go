package email

import "time"

type EncryptedKey struct {
	RecipientID         uint   `json:"recipient_id"`
	EncryptedPassphrase []byte `json:"encrypted_passphrase"`
}

type Message struct {
	ID                   uint `gorm:"primaryKey"`
	ConversationID       uint
	SenderID             uint
	RecipientsJSON       string `gorm:"type:text"` // JSON array of recipient IDs
	EncryptedBody        []byte
	EncryptedSessionKeys string `gorm:"type:text"` // JSON array of EncryptedKey
	EncryptedAttachments []byte
	Metadata             string `gorm:"type:text"` // JSON string for additional data
	Status               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	SentAt               time.Time
}
