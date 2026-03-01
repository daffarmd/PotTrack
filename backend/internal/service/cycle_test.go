package service

import "testing"

func TestArchiveCycle(t *testing.T) {
	c := &Cycle{Status: "aktif"}
	ArchiveCycle(c, "selesai")
	if c.Status != "arsip" {
		t.Errorf("expected status arsip, got %s", c.Status)
	}
	if c.AlasanArsip != "selesai" {
		t.Errorf("expected alasan selesai, got %s", c.AlasanArsip)
	}
}
