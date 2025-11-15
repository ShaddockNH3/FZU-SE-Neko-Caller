package model

import (
	"time"

	"gorm.io/datatypes"
)

// ScoreEvent records every points adjustment for auditing and analytics.
type ScoreEvent struct {
	EventID      string            `gorm:"type:varchar(255);primaryKey"`
	EnrollmentID string            `gorm:"type:varchar(255);index"`
	StudentID    string            `gorm:"type:varchar(255);index"`
	ClassID      string            `gorm:"type:varchar(255);index"`
	Delta        float64           `gorm:"not null"`
	Reason       string            `gorm:"type:varchar(255);not null"`
	EventType    int32             `gorm:"not null"`
	Metadata     datatypes.JSONMap `gorm:"type:json"`
	CreatedAt    time.Time         `gorm:"autoCreateTime"`
}

func (ScoreEvent) TableName() string {
	return "score_events"
}
