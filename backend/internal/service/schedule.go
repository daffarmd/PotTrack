package service

import (
	"encoding/json"
	"time"
)

// Schedule represents the JSON structure for tugas schedule

type Schedule struct {
	Type            string `json:"type"`
	DueAt           string `json:"due_at,omitempty"`
	Every           int    `json:"every,omitempty"`
	Unit            string `json:"unit,omitempty"`
	TimeOfDay       string `json:"time_of_day,omitempty"`
	StartFrom       string `json:"start_from,omitempty"`
	StageName       string `json:"stage_name,omitempty"`
	StartOffsetDays int    `json:"start_offset_days,omitempty"`
	End             *struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"end,omitempty"`
}

func ParseSchedule(jsonStr string) (*Schedule, error) {
	var s Schedule
	if err := json.Unmarshal([]byte(jsonStr), &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// NextDue computes next due given current due and schedule
func NextDue(current time.Time, s *Schedule) time.Time {
	switch s.Type {
	case "one_time":
		if s.DueAt == "" {
			return current
		}
		t, _ := time.Parse(time.RFC3339, s.DueAt)
		return t
	case "interval":
		switch s.Unit {
		case "day":
			return current.Add(time.Duration(s.Every) * 24 * time.Hour)
		case "week":
			return current.Add(time.Duration(s.Every*7) * 24 * time.Hour)
		default:
			return current
		}
	default:
		return current
	}
}
