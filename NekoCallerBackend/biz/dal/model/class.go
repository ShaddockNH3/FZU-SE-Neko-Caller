package model

type Class struct {
	ClassID   string `gorm:"type:varchar(255);primaryKey"` // 班级ID，是主键
	ClassName string `gorm:"type:varchar(255);not null"`   // 班级名称
}

func (Class) TableName() string {
	return "classes"
}
