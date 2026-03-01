package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// IncidentHandler handles insiden updates
type IncidentHandler struct {
	db *sqlx.DB
}

func NewIncidentHandler(db *sqlx.DB) *IncidentHandler {
	return &IncidentHandler{db: db}
}

func (h *IncidentHandler) Close(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Exec("UPDATE insiden SET status='selesai', ditutup_pada=now() WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not close"})
		return
	}
	c.Status(http.StatusOK)
}
