package service

import (
	"testing"
	"time"
)

func TestNextDue_OneTime(t *testing.T) {
	s := &Schedule{Type: "one_time", DueAt: "2026-03-10T08:00:00Z"}
	now := time.Now()
	nd := NextDue(now, s)
	if !nd.Equal(time.Date(2026, 3, 10, 8, 0, 0, 0, time.UTC)) {
		t.Errorf("expected 2026-03-10T08:00:00Z got %v", nd)
	}
}

func TestNextDue_IntervalDay(t *testing.T) {
	s := &Schedule{Type: "interval", Every: 2, Unit: "day"}
	now := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	nd := NextDue(now, s)
	expected := now.Add(48 * time.Hour)
	if !nd.Equal(expected) {
		t.Errorf("expected %v got %v", expected, nd)
	}
}

func TestNextDue_IntervalWeek(t *testing.T) {
	s := &Schedule{Type: "interval", Every: 1, Unit: "week"}
	now := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	nd := NextDue(now, s)
	expected := now.Add(7 * 24 * time.Hour)
	if !nd.Equal(expected) {
		t.Errorf("expected %v got %v", expected, nd)
	}
}
