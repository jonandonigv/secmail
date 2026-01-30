package handlers

import (
	"fmt"
	"net/http"
	"secmail/internal/email"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SendEmailRequest struct {
	Recipients []uint `json:"recipients" binding:"required,dive,min=1,max=10"`
	Subject    string `json:"subject" binding:"required,max=100"`
	Body       string `json:"body" binding:"required,max=10000"`
}

type InboxResponse struct {
	Messages []email.DecryptedMessage `json:"messages"`
}

// SendEmail handles sending an email
func SendEmail(c *gin.Context, db *gorm.DB) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	var req SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sanitize inputs
	req.Subject = strings.TrimSpace(req.Subject)
	req.Body = strings.TrimSpace(req.Body)

	err := email.SendMessage(userID, req.Recipients, req.Subject, req.Body, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

// GetInbox handles retrieving the user's inbox
func GetInbox(c *gin.Context, db *gorm.DB) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	messages, err := email.GetInbox(userID, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := InboxResponse{Messages: messages}
	c.JSON(http.StatusOK, response)
}

type SentResponse struct {
	Messages []email.SentMessage `json:"messages"`
}

// GetSent handles retrieving the user's sent emails
func GetSent(c *gin.Context, db *gorm.DB) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	messages, err := email.GetSent(userID, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := SentResponse{Messages: messages}
	c.JSON(http.StatusOK, response)
}

// DeleteEmail handles deleting an email
func DeleteEmail(c *gin.Context, db *gorm.DB) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	// Parse message ID from URL parameter
	messageIDStr := c.Param("id")
	var messageID uint
	if _, err := fmt.Sscanf(messageIDStr, "%d", &messageID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	if err := email.DeleteMessage(userID, messageID, db); err != nil {
		if err.Error() == "unauthorized: user is not sender or recipient" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email deleted successfully"})
}

type ReplyRequest struct {
	MessageID uint   `json:"message_id" binding:"required,min=1"`
	Body      string `json:"body" binding:"required,max=10000"`
}

// ReplyEmail handles replying to an email
func ReplyEmail(c *gin.Context, db *gorm.DB) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	var req ReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sanitize input
	req.Body = strings.TrimSpace(req.Body)

	if err := email.SendReply(userID, req.MessageID, req.Body, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reply sent successfully"})
}

type ConversationResponse struct {
	Messages []email.DecryptedMessage `json:"messages"`
}

// GetConversation handles retrieving a conversation thread
func GetConversation(c *gin.Context, db *gorm.DB) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDVal.(uint)

	// Parse conversation ID from URL parameter
	conversationIDStr := c.Param("id")
	var conversationID uint
	if _, err := fmt.Sscanf(conversationIDStr, "%d", &conversationID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	messages, err := email.GetConversation(userID, conversationID, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := ConversationResponse{Messages: messages}
	c.JSON(http.StatusOK, response)
}
