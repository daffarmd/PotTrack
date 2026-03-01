package service

// Cycle represents minimal data for archiving logic

type Cycle struct {
	Status      string
	AlasanArsip string
}

// ArchiveCycle marks a cycle as archived with reason
func ArchiveCycle(c *Cycle, reason string) {
	c.Status = "arsip"
	c.AlasanArsip = reason
}
