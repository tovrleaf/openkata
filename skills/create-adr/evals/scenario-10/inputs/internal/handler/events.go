package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// EventHandler handles inbound event submissions from external data sources.
type EventHandler struct{}

// Submit accepts a batch of events and synchronously writes them to the database.
// This is a direct write path — no queue in between — which has started causing
// latency spikes under load from high-volume data sources.
func (h *EventHandler) Submit(c *gin.Context) {
	var payload struct {
		SourceID string                   `json:"source_id" binding:"required"`
		Events   []map[string]interface{} `json:"events" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: decouple ingestion from write path — high-volume sources are causing
	// p99 latency to spike when the DB is under write pressure.
	c.JSON(http.StatusAccepted, gin.H{"accepted": len(payload.Events)})
}
