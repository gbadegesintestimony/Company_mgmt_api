package models

type AuditLog struct {
	ID        string
	ActorID   string
	Action    string
	TargetID  string
	Metadata  string
	CreatedAt string
}
