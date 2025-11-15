package model

type Student struct {
	StudentID string `gorm:"type:varchar(255);primaryKey"` // 学号，是主键
	Name      string `gorm:"type:varchar(255);not null"`   // 姓名
	Major     string `gorm:"type:varchar(255)"`            // 专业
}

func (Student) TableName() string {
	return "students"
}
