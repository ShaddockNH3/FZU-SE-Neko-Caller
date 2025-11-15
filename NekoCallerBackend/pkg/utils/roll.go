package utils

import (
	"FZUSENekoCaller/biz/model/common"
	"FZUSENekoCaller/pkg/errno"
)

func Roll(roster []*common.RosterItem, mode common.RollCallMode, event common.RandomEventType) string {
	if len(roster) == 0 {
		return errno.StudentNotFoundErr.Error()
	}

	switch mode {
	case common.RollCallMode_RANDOM:
		if event == common.RandomEventType_CRAZY_THURSDAY{
			return rollWithCrazyThursday(roster)
		}
		return ""
	case common.RollCallMode_SEQUENTIAL:
		return roster[0].StudentInfo.StudentID
	case common.RollCallMode_REVERSE_SEQUENTIAL:
		return roster[len(roster)-1].StudentInfo.StudentID
	case common.RollCallMode_LOW_POINTS_FIRST:
		// 积分低优先，按照正态分布处理
		return ""
	default:
		return ""
	}
}

func rollWithCrazyThursday(roster []*common.RosterItem) string {
	return ""
}
