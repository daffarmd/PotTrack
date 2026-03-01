package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func TestTaskCompleteCreatesNoteAndAdvances(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	sqlxdb := sqlx.NewDb(db, "postgres")
	handler := NewTaskHandler(sqlxdb)

	// prepare task fetch
	mock.ExpectQuery("SELECT .* FROM tugas WHERE id=\\$1").WithArgs("1").WillReturnRows(
		sqlmock.NewRows([]string{"id", "jadwal_json", "next_due_at", "tipe", "aktif", "id_siklus"}).
			AddRow(1, `{\"type\":\"interval\",\"every\":1,\"unit\":\"day\"}`, time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), "interval", true, 10))

	// expect catatan insert (4 args: id_siklus, text, fotoJson, tugas_id)
	mock.ExpectExec("INSERT INTO catatan").WithArgs(10, "note", sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(1, 1))

	// expect tugas update
	mock.ExpectExec("UPDATE tugas SET last_completed_at=now\\(\\), next_due_at=\\$1 WHERE id=\\$2").WillReturnResult(sqlmock.NewResult(1, 1))

	// create request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/tugas/1/selesai", strings.NewReader(`{"catatan_teks":"note","foto_urls":["photo1"]}`))
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	handler.Complete(c)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 got %d body %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
