package service

import "time"

// Task represents minimal task state
type Task struct {
	NextDueAt time.Time
	Tipe      string
	Aktif     bool
}

// CompleteTask marks completion and returns updated next due
func CompleteTask(t *Task, schedule *Schedule) time.Time {
	if schedule != nil && t.Tipe == "interval" && t.Aktif {
		newDue := NextDue(t.NextDueAt, schedule)
		t.NextDueAt = newDue
		return newDue
	}
	return t.NextDueAt
}
