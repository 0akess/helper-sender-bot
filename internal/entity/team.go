package entity

import "github.com/google/uuid"

type Team struct {
	Name    string
	Token   uuid.UUID
	LeadEID string
}

type AuthMeta struct {
	Team  string
	Token uuid.UUID
}
