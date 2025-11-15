package model

import "time"

// RollCallRecord stores every roll-call action for later statistics.
type RollCallRecord struct {
	RecordID     string    `gorm:"type:varchar(255);primaryKey"`
	ClassID      string    `gorm:"type:varchar(255);index"`
	EnrollmentID string    `gorm:"type:varchar(255);index"`
	StudentID    string    `gorm:"type:varchar(255);index"`
	Mode         int32     `gorm:"not null"`
	EventType    int32     `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (RollCallRecord) TableName() string {
	return "roll_call_records"
}
