package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// TaskHandler handles tugas endpoints
type TaskHandler struct {
	db *sqlx.DB
}

func NewTaskHandler(db *sqlx.DB) *TaskHandler {
	return &TaskHandler{db: db}
}

func (h *TaskHandler) List(c *gin.Context) {
	filter := c.Query("filter")
	var tasks []map[string]interface{}
	q := "SELECT * FROM tugas WHERE 1=1"
	if filter == "today" {
		q += " AND next_due_at::date = current_date"
	} else if filter == "overdue" {
		q += " AND next_due_at < now()"
	}
	h.db.Select(&tasks, q)
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var fields map[string]interface{}
	if err := c.ShouldBindJSON(&fields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	args = append(args, id)
	_, err := h.db.Exec("UPDATE tugas SET "+set+" WHERE id=$"+strconv.Itoa(i), args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update"})
		return
	}
	c.Status(http.StatusOK)
}

func (h *TaskHandler) Complete(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		CatatanTeks string   `json:"catatan_teks"`
		FotoURLs    []string `json:"foto_urls"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fetch existing
	var tugas struct {
		ID         int       `db:"id"`
		JadwalJSON string    `db:"jadwal_json"`
		NextDueAt  time.Time `db:"next_due_at"`
		Tipe       string    `db:"tipe"`
		Aktif      bool      `db:"aktif"`
		IdSiklus   int       `db:"id_siklus"`
	}
	h.db.Get(&tugas, "SELECT * FROM tugas WHERE id=$1", id)
	// create catatan if provided
	if req.CatatanTeks != "" || len(req.FotoURLs) > 0 {
		fotoJson := "[]"
		if len(req.FotoURLs) > 0 {
			fotoJson = "[\"" + strings.Join(req.FotoURLs, "\",\"") + "\"]"
		}
		h.db.Exec("INSERT INTO catatan (id_siklus, id_pot, catatan_teks, foto_urls_json, id_tugas_terkait, dibuat_pada) VALUES ($1, (SELECT id_pot FROM siklus WHERE id=$1), $2, $3, $4, now())", tugas.IdSiklus, req.CatatanTeks, fotoJson, tugas.ID)
	}
	// advance next_due_at for interval
	if tugas.Tipe == "interval" && tugas.Aktif {
		// simplistic: add one day
		newNext := tugas.NextDueAt.Add(24 * time.Hour)
		h.db.Exec("UPDATE tugas SET last_completed_at=now(), next_due_at=$1 WHERE id=$2", newNext, tugas.ID)
	} else {
		h.db.Exec("UPDATE tugas SET last_completed_at=now() WHERE id=$1", tugas.ID)
	}
	c.Status(http.StatusOK)
}
