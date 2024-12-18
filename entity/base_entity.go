package entity

import "time"

type BaseEntity struct {
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at;default:now()"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;default:'0001-01-01 00:00:00Z'"`
}
