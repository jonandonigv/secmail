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

type SentMessage struct {
	ID         uint
	Recipients []uint
	Subject    string
	Body       string
	Status     string
	SentAt     time.Time
}

// GetSent retrieves all messages sent by the given user.
func GetSent(userID uint, db *gorm.DB) ([]SentMessage, error) {
	// Query messages where user is the sender
	var messages []Message
	if err := db.Where("sender_id = ?", userID).Order("sent_at DESC").Find(&messages).Error; err != nil {
		return nil, err
	}

	var sentMessages []SentMessage
	for _, msg := range messages {
		// Parse recipients
		var recipients []uint
		if err := json.Unmarshal([]byte(msg.RecipientsJSON), &recipients); err != nil {
			return nil, err
		}

		// Parse metadata for subject
		var metadata map[string]string
		if err := json.Unmarshal([]byte(msg.Metadata), &metadata); err != nil {
			return nil, err
		}
		subject := metadata["subject"]

		// For sent messages, we need to decrypt using the first recipient's key
		// In a real system, we'd store the plaintext or use a different approach
		// For now, we'll mark the body as encrypted
		sentMessages = append(sentMessages, SentMessage{
			ID:         msg.ID,
			Recipients: recipients,
			Subject:    subject,
			Body:       "[Encrypted - view in sent folder not yet implemented]",
			Status:     msg.Status,
			SentAt:     msg.SentAt,
		})
	}

	return sentMessages, nil
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

// GetConversation retrieves all messages in a conversation thread.
func GetConversation(userID uint, conversationID uint, db *gorm.DB) ([]DecryptedMessage, error) {
	// Get user to access private key
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// Query all messages in the conversation
	var messages []Message
	if err := db.Where("id = ? OR conversation_id = ?", conversationID, conversationID).
		Order("sent_at ASC").Find(&messages).Error; err != nil {
		return nil, err
	}

	var decryptedMessages []DecryptedMessage
	for _, msg := range messages {
		// Check if user is sender or recipient
		isAuthorized := false
		if msg.SenderID == userID {
			isAuthorized = true
		} else {
			// Check if user is a recipient
			var recipients []uint
			if err := json.Unmarshal([]byte(msg.RecipientsJSON), &recipients); err != nil {
				return nil, err
			}
			for _, r := range recipients {
				if r == userID {
					isAuthorized = true
					break
				}
			}
		}

		if !isAuthorized {
			continue // Skip messages user is not authorized to see
		}

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

		// If user is sender, we need to handle differently
		// For now, skip messages where we can't decrypt
		if !found && msg.SenderID != userID {
			continue
		}

		var body string
		if found {
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
			body = string(bodyBytes)
		} else {
			body = "[Sent by you - content not stored in plaintext]"
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
			Body:           body,
			Status:         msg.Status,
			SentAt:         msg.SentAt,
		})
	}

	return decryptedMessages, nil
}

// DeleteMessage soft-deletes a message if the user is the sender or recipient.
func DeleteMessage(userID uint, messageID uint, db *gorm.DB) error {
	// Get the message
	var msg Message
	if err := db.Where("id = ?", messageID).First(&msg).Error; err != nil {
		return err
	}

	// Check if user is the sender
	if msg.SenderID == userID {
		return db.Delete(&msg).Error
	}

	// Check if user is a recipient
	var recipients []uint
	if err := json.Unmarshal([]byte(msg.RecipientsJSON), &recipients); err != nil {
		return err
	}

	for _, recipientID := range recipients {
		if recipientID == userID {
			return db.Delete(&msg).Error
		}
	}

	return errors.New("unauthorized: user is not sender or recipient")
}
