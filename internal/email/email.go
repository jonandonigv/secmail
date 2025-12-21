package email

type Message struct {
	ID                   int
	ConversationID       int
	SenderID             int
	RecipientsID         int
	EncryptedBody        string
	EncryptedSessionKeys []byte
	EncryptedAttachments []byte
	Metadata             string
	CreatedAt            string
}
