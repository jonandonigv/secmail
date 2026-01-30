package email

import (
	"encoding/json"
	"errors"
	"secmail/internal/crypto"
	"secmail/internal/models"
	"time"

	"gorm.io/gorm"
)

// SendReply sends a reply to an existing conversation.
func SendReply(senderID uint, conversationID uint, body string, db *gorm.DB) error {
	// Get the original message to extract recipients
	var originalMsg Message
	if err := db.Where("id = ?", conversationID).First(&originalMsg).Error; err != nil {
		return err
	}

	// Parse original recipients
	var originalRecipients []uint
	if err := json.Unmarshal([]byte(originalMsg.RecipientsJSON), &originalRecipients); err != nil {
		return err
	}

	// Build recipient list: original sender + all original recipients (excluding current sender)
	recipients := []uint{originalMsg.SenderID}
	for _, r := range originalRecipients {
		if r != senderID {
			recipients = append(recipients, r)
		}
	}

	// Parse metadata for subject
	var metadata map[string]string
	if err := json.Unmarshal([]byte(originalMsg.Metadata), &metadata); err != nil {
		return err
	}
	subject := "Re: " + metadata["subject"]

	// Encrypt the body
	encryptedBody, passphrase, err := crypto.EncryptBody([]byte(body))
	if err != nil {
		return err
	}

	// Get public keys for recipients
	var users []models.User
	if err := db.Where("id IN ?", recipients).Find(&users).Error; err != nil {
		return err
	}

	// Encrypt passphrase for each recipient
	var encryptedKeys []EncryptedKey
	for _, user := range users {
		encryptedPass, err := crypto.EncryptPassphrase(passphrase, user.PublicKey)
		if err != nil {
			return err
		}
		encryptedKeys = append(encryptedKeys, EncryptedKey{
			RecipientID:         user.ID,
			EncryptedPassphrase: encryptedPass,
		})
	}

	// Marshal recipients and encrypted keys
	recipientsJSON, err := json.Marshal(recipients)
	if err != nil {
		return err
	}
	encryptedKeysJSON, err := json.Marshal(encryptedKeys)
	if err != nil {
		return err
	}

	// Create metadata
	replyMetadata := map[string]string{"subject": subject}
	metadataJSON, err := json.Marshal(replyMetadata)
	if err != nil {
		return err
	}

	// Create message
	message := Message{
		SenderID:             senderID,
		ConversationID:       conversationID,
		RecipientsJSON:       string(recipientsJSON),
		EncryptedBody:        encryptedBody,
		EncryptedSessionKeys: string(encryptedKeysJSON),
		Metadata:             string(metadataJSON),
		Status:               "sent",
		SentAt:               time.Now(),
	}
	if err := db.Create(&message).Error; err != nil {
		return err
	}

	return nil
}

// SendMessage sends an encrypted email from sender to recipients.
func SendMessage(senderID uint, recipients []uint, subject, body string, db *gorm.DB) error {
	if len(recipients) == 0 {
		return errors.New("no recipients")
	}

	// Encrypt the body
	encryptedBody, passphrase, err := crypto.EncryptBody([]byte(body))
	if err != nil {
		return err
	}

	// Get public keys for recipients
	var users []models.User
	if err := db.Where("id IN ?", recipients).Find(&users).Error; err != nil {
		return err
	}
	if len(users) != len(recipients) {
		return errors.New("some recipients not found")
	}

	// Encrypt passphrase for each recipient
	var encryptedKeys []EncryptedKey
	for _, user := range users {
		encryptedPass, err := crypto.EncryptPassphrase(passphrase, user.PublicKey)
		if err != nil {
			return err
		}
		encryptedKeys = append(encryptedKeys, EncryptedKey{
			RecipientID:         user.ID,
			EncryptedPassphrase: encryptedPass,
		})
	}

	// Marshal recipients and encrypted keys
	recipientsJSON, err := json.Marshal(recipients)
	if err != nil {
		return err
	}
	encryptedKeysJSON, err := json.Marshal(encryptedKeys)
	if err != nil {
		return err
	}

	// Create metadata
	metadata := map[string]string{"subject": subject}
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	// Create message
	message := Message{
		SenderID:             senderID,
		RecipientsJSON:       string(recipientsJSON),
		EncryptedBody:        encryptedBody,
		EncryptedSessionKeys: string(encryptedKeysJSON),
		Metadata:             string(metadataJSON),
		Status:               "sent",
		SentAt:               time.Now(),
	}
	if err := db.Create(&message).Error; err != nil {
		return err
	}

	return nil
}
