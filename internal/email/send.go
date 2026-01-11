package email

import (
	"encoding/json"
	"errors"
	"secmail/internal/crypto"
	"secmail/internal/models"
	"time"

	"gorm.io/gorm"
)

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
