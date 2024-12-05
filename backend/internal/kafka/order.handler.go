package kafka

import (
	"net/http"
	"strconv"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/orders/types"
	"github.com/gin-gonic/gin"
)

// test handler for all orders retrieval
func (k *KafkaManager) GetOrders(c *gin.Context) {
	storeID, err := strconv.Atoi(c.Query("storeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid store ID"})
		return
	}

	topic := k.GetTopic(k.Topics.ActiveOrders)
	messages, err := k.FetchOrders(topic) // Fetch from partition 0
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter messages by storeID
	var filteredMessages []types.OrderEvent
	for _, message := range messages {
		if message.StoreID == uint(storeID) {
			filteredMessages = append(filteredMessages, message)
		}
	}

	c.JSON(http.StatusOK, filteredMessages)
}
