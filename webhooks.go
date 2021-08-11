package gozabo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type WebhooksListeners struct {
	OnAccountCreated func(Account, url.Values)
}

type WebhookEvent struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"` // keep as raw type (json string) to decode after depending on event type
}

func (ctrl *Controller) GinWebhookEndpoint(c *gin.Context) {
	var event WebhookEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		fmt.Printf("Failed to unmarshal webhook payload: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	fmt.Sprintf("received webhook event %s with data %s", event.Event, event.Data)

	switch event.Event {
	case "account.post":
		ctrl.whPostAccount(c, event.Data)
	}
}

func (ctrl *Controller) whPostAccount(c *gin.Context, data []byte) {
	var err error

	var account Account
	if err = json.Unmarshal(data, &account); err != nil {
		fmt.Printf("Failed to unmarshal account: %s\nraw account is: %+v", err, string(data))
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// meta
	metaHeader := c.GetHeader("X-Connect-Meta")
	var meta url.Values
	if metaHeader != "" {
		meta, err = url.ParseQuery(metaHeader)
		if err != nil {
			fmt.Printf("Failed parse meta content: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}

	ctrl.whListeners.OnAccountCreated(account, meta)
}
