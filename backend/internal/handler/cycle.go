package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// CycleHandler handles siklus endpoints
type CycleHandler struct {
	db *sqlx.DB
}

func NewCycleHandler(db *sqlx.DB) *CycleHandler {
	return &CycleHandler{db: db}
}

func (h *CycleHandler) Archive(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Alasan string `json:"alasan" binding:"required,oneof=selesai gagal"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.db.Exec("UPDATE siklus SET status='arsip', alasan_arsip=$1, diarsipkan_pada=now() WHERE id=$2", req.Alasan, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not archive"})
		return
	}
	c.Status(http.StatusOK)
}

func (h *CycleHandler) GetStages(c *gin.Context) {
	id := c.Param("id")
	var stages []map[string]interface{}
	h.db.Select(&stages, "SELECT * FROM tahap WHERE id_siklus=$1 ORDER BY urutan", id)
	c.JSON(http.StatusOK, stages)
}

func (h *CycleHandler) UpdateStages(c *gin.Context) {
	id := c.Param("id")
	var arr []struct {
		NamaTahap  string `json:"nama_tahap"`
		Urutan     int    `json:"urutan"`
		DurasiHari int    `json:"durasi_hari"`
	}
	if err := c.ShouldBindJSON(&arr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tx := h.db.MustBegin()
	tx.Exec("DELETE FROM tahap WHERE id_siklus=$1", id)
	for _, s := range arr {
		tx.Exec("INSERT INTO tahap (id_siklus,nama_tahap,urutan,durasi_hari) VALUES ($1,$2,$3,$4)", id, s.NamaTahap, s.Urutan, s.DurasiHari)
	}
	tx.Commit()
	c.Status(http.StatusOK)
}

func (h *CycleHandler) CreateTask(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Judul      string `json:"judul" binding:"required"`
		Tipe       string `json:"tipe" binding:"required,oneof=one_time interval"`
		JadwalJSON string `json:"jadwal_json" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.db.Exec("INSERT INTO tugas (id_siklus, judul, tipe, jadwal_json, aktif, dibuat_pada) VALUES ($1,$2,$3,$4,true,now())", id, req.Judul, req.Tipe, req.JadwalJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create task"})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *CycleHandler) ListNotes(c *gin.Context) {
	id := c.Param("id")
	var notes []map[string]interface{}
	h.db.Select(&notes, "SELECT * FROM catatan WHERE id_siklus=$1 ORDER BY dibuat_pada DESC", id)
	c.JSON(http.StatusOK, notes)
}

func (h *CycleHandler) AddNote(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		CatatanTeks string   `json:"catatan_teks"`
		FotoURLs    []string `json:"foto_urls"`
		IdTugas     *int     `json:"id_tugas_terkait"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fotoJson := "[]"
	if len(req.FotoURLs) > 0 {
		// simple encode
		fotoJson = "[\"" + strings.Join(req.FotoURLs, "\",\"") + "\"]"
	}
	_, err := h.db.Exec("INSERT INTO catatan (id_siklus, id_pot, catatan_teks, foto_urls_json, id_tugas_terkait, dibuat_pada) VALUES ($1, (SELECT id_pot FROM siklus WHERE id=$1), $2, $3, $4, now())", id, req.CatatanTeks, fotoJson, req.IdTugas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add note"})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *CycleHandler) ListIncidents(c *gin.Context) {
	id := c.Param("id")
	var ins []map[string]interface{}
	h.db.Select(&ins, "SELECT * FROM insiden WHERE id_siklus=$1 ORDER BY dibuka_pada DESC", id)
	c.JSON(http.StatusOK, ins)
}

func (h *CycleHandler) CreateIncident(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Jenis    string   `json:"jenis" binding:"required"`
		Tingkat  int      `json:"tingkat" binding:"required,oneof=1 2 3"`
		Catatan  string   `json:"catatan"`
		FotoURLs []string `json:"foto_urls"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fotoJson := "[]"
	if len(req.FotoURLs) > 0 {
		fotoJson = "[\"" + strings.Join(req.FotoURLs, "\",\"") + "\"]"
	}
	_, err := h.db.Exec("INSERT INTO insiden (id_siklus, dibuka_pada, jenis, tingkat, status, catatan, foto_urls_json) VALUES ($1, now(), $2, $3, 'buka', $4, $5)", id, req.Jenis, req.Tingkat, req.Catatan, fotoJson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create incident"})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *CycleHandler) ListHarvest(c *gin.Context) {
	id := c.Param("id")
	var p []map[string]interface{}
	h.db.Select(&p, "SELECT * FROM panen WHERE id_siklus=$1 ORDER BY tanggal DESC", id)
	c.JSON(http.StatusOK, p)
}

func (h *CycleHandler) CreateHarvest(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Tanggal string `json:"tanggal" binding:"required"`
		Jumlah  string `json:"jumlah" binding:"required"`
		Satuan  string `json:"satuan" binding:"required"`
		Catatan string `json:"catatan"`
		Grade   string `json:"grade"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.db.Exec("INSERT INTO panen (id_siklus, tanggal, jumlah, satuan, catatan, grade) VALUES ($1,$2,$3,$4,$5,$6)", id, req.Tanggal, req.Jumlah, req.Satuan, req.Catatan, req.Grade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create harvest"})
		return
	}
	c.Status(http.StatusCreated)
}
