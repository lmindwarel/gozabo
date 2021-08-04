package gobudins

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebhooksListeners struct {
	OnAccountCreated func(Account)
}

type WebhookEvent struct {
	Event string `json:"event"`
	Data  string `json:"data"` // keep as raw type (json string) to decode after depending on event type
}

func (ctrl *Controller) GinWebhookEndpoint(c *gin.Context) {
	var event WebhookEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		fmt.Printf("Failed to unmarshal created user: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	fmt.Sprintf("received webhook event %s with data %s", event.Event, event.Data)

	switch event.Event {
	case "account.post":
		ctrl.whPostAccount(c, event.Data)
	}
}

func (ctrl *Controller) whPostAccount(c *gin.Context, data string) {
	var account Account
	if err := json.Unmarshal([]byte(data), &account); err != nil {
		fmt.Printf("Failed to unmarshal account: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ctrl.whListeners.OnAccountCreated(account)
}
