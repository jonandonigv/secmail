package handlers

import (
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
