package gobudins

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebhookData struct {
	Event string `json:"event"`
	Data  string `json:"data"` // keep as raw type (json string) to decode after depending on event type
}

type WebhookPostAccount struct {
	Event   string  `json:"event"`
	Account Account `json:"data"` // keep as raw type (json string) to decode after depending on event type
}

func (ctrl *Controller) GinWebhookEndpoint(c *gin.Context) {

	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Printf("Failed to unmarshal created user: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctrl.listeners.OnUserCreated(user)
}

func (ctrl *Controller) whAccountSynced(c *gin.Context) {
	var account SyncedAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		fmt.Printf("Failed to unmarshal synced account: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctrl.listeners.OnAccountSynced(account)
}

func (ctrl *Controller) whAccountDisabled(c *gin.Context) {
	var account SyncedAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		fmt.Printf("Failed to unmarshal synced account: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctrl.listeners.OnAccountDisabled(account)
}

func (ctrl *Controller) whConnectionDeleted(c *gin.Context) {
	var connection Connection
	if err := c.ShouldBindJSON(&connection); err != nil {
		fmt.Printf("Failed to unmarshal connection: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctrl.listeners.OnConnectionDeleted(connection)
}
