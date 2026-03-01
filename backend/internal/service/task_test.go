package service

import (
	"testing"
	"time"
)

func TestCompleteTaskInterval(t *testing.T) {
	tsk := &Task{NextDueAt: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), Tipe: "interval", Aktif: true}
	sched := &Schedule{Type: "interval", Every: 1, Unit: "day"}
	newDue := CompleteTask(tsk, sched)
	expected := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)
	if !newDue.Equal(expected) {
		t.Errorf("expected %v got %v", expected, newDue)
	}
}
