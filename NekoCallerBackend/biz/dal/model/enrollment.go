package model

type Enrollment struct {
	EnrollmentID   string  `gorm:"type:varchar(255);primaryKey"` // 记录的唯一ID
	StudentID      string  `gorm:"type:varchar(255);not null"`   // 关联到学生ID
	ClassID        string  `gorm:"type:varchar(255);not null"`   // 关联到班级ID
	TotalPoints    float64 `gorm:"default:0"`                    // 该学生在该班级的总积分
	CallCount      int64   `gorm:"default:0"`                    // 该学生在该班级的被点名次数
	TransferRights int64   `gorm:"default:0"`                    // 该学生在该班级的点名转移权

	Student Student `gorm:"foreignKey:StudentID"` // 外键查询
	Class   Class   `gorm:"foreignKey:ClassID"`
}

func (Enrollment) TableName() string {
	return "enrollments"
}
