package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Pot represents a pot record

// PotHandler handles pot endpoints
type PotHandler struct {
	db *sqlx.DB
}

func NewPotHandler(db *sqlx.DB) *PotHandler {
	return &PotHandler{db: db}
}

func (h *PotHandler) List(c *gin.Context) {
	var pots []map[string]interface{}
	h.db.Select(&pots, "SELECT * FROM pot WHERE id_pengguna=$1", h.userID(c))
	c.JSON(http.StatusOK, pots)
}

func (h *PotHandler) Create(c *gin.Context) {
	var req struct {
		NamaPot      string `json:"nama_pot" binding:"required"`
		NamaTanaman  string `json:"nama_tanaman" binding:"required"`
		TanggalMulai string `json:"tanggal_mulai"`
		Varietas     string `json:"varietas"`
		UkuranPot    string `json:"ukuran_pot"`
		MediaTanam   string `json:"media_tanam"`
		Lokasi       string `json:"lokasi"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.TanggalMulai == "" {
		req.TanggalMulai = "now()"
	}
	_, err := h.db.Exec(`INSERT INTO pot (id_pengguna,nama_pot,nama_tanaman,varietas,ukuran_pot,media_tanam,lokasi,dibuat_pada,tanggal_mulai)
        VALUES ($1,$2,$3,$4,$5,$6,$7,now(),$8)`, h.userID(c), req.NamaPot, req.NamaTanaman, req.Varietas, req.UkuranPot, req.MediaTanam, req.Lokasi, req.TanggalMulai)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create pot"})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *PotHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var pot map[string]interface{}
	err := h.db.Get(&pot, "SELECT * FROM pot WHERE id=$1 AND id_pengguna=$2", id, h.userID(c))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, pot)
}

func (h *PotHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var fields map[string]interface{}
	if err := c.ShouldBindJSON(&fields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Simple dynamic update
	set := ""
	args := []interface{}{}
	i := 1
	for k, v := range fields {
		set += k + "=$" + strconv.Itoa(i) + ","
		args = append(args, v)
		i++
	}
	if set == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	set = set[:len(set)-1]
	args = append(args, id, h.userID(c))
	_, err := h.db.Exec("UPDATE pot SET "+set+" WHERE id=$"+strconv.Itoa(i)+" AND id_pengguna=$"+strconv.Itoa(i+1), args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update pot"})
		return
	}
	c.Status(http.StatusOK)
}

func (h *PotHandler) StartCycle(c *gin.Context) {
	// logic to insert siklus; basic version
	potID := c.Param("id")
	_, err := h.db.Exec("INSERT INTO siklus (id_pot,tanggal_mulai,status) VALUES ($1,now(),'aktif')", potID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not start cycle"})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *PotHandler) userID(c *gin.Context) int {
	if v, ok := c.Get("user_id"); ok {
		switch v := v.(type) {
		case float64:
			return int(v)
		case int:
			return v
		}
	}
	return 0
}
