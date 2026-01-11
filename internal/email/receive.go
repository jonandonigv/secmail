package email

import (
	"encoding/json"
	"errors"
	"secmail/internal/crypto"
	"secmail/internal/models"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type DecryptedMessage struct {
	ID             uint
	ConversationID uint
	SenderID       uint
	Subject        string
	Body           string
	Status         string
	SentAt         time.Time
}

// GetInbox retrieves and decrypts messages for the given user.
func GetInbox(userID uint, db *gorm.DB) ([]DecryptedMessage, error) {
	// Get user to access private key
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// Query messages where user is recipient (simple LIKE check)
	userIDStr := "\"" + strconv.Itoa(int(userID)) + "\""
	var messages []Message
	if err := db.Where("recipients_json LIKE ?", "%"+userIDStr+"%").Find(&messages).Error; err != nil {
		return nil, err
	}

	var decryptedMessages []DecryptedMessage
	for _, msg := range messages {
		// Parse encrypted keys
		var encryptedKeys []EncryptedKey
		if err := json.Unmarshal([]byte(msg.EncryptedSessionKeys), &encryptedKeys); err != nil {
			return nil, err
		}

		// Find the key for this user
		var encryptedPass []byte
		found := false
		for _, key := range encryptedKeys {
			if key.RecipientID == userID {
				encryptedPass = key.EncryptedPassphrase
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("session key not found for user")
		}

		// Decrypt passphrase
		passphrase, err := crypto.DecryptPassphrase(encryptedPass, user.PrivateKey)
		if err != nil {
			return nil, err
		}

		// Decrypt body
		bodyBytes, err := crypto.DecryptBody(msg.EncryptedBody, passphrase)
		if err != nil {
			return nil, err
		}

		// Parse metadata for subject
		var metadata map[string]string
		if err := json.Unmarshal([]byte(msg.Metadata), &metadata); err != nil {
			return nil, err
		}
		subject := metadata["subject"]

		decryptedMessages = append(decryptedMessages, DecryptedMessage{
			ID:             msg.ID,
			ConversationID: msg.ConversationID,
			SenderID:       msg.SenderID,
			Subject:        subject,
			Body:           string(bodyBytes),
			Status:         msg.Status,
			SentAt:         msg.SentAt,
		})
	}

	return decryptedMessages, nil
}
