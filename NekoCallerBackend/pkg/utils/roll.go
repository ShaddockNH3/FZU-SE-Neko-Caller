package utils

import (
	"FZUSENekoCaller/biz/model/common"
	"FZUSENekoCaller/pkg/errno"
)

func Roll(roster []*common.RosterItem, mode common.RollCallMode, event common.RandomEventType) (common.RosterItem, error) {
	if len(roster) == 0 {
		return common.RosterItem{}, errno.StudentNotFoundErr
	}

	switch mode {
	case common.RollCallMode_RANDOM:
		switch event {
		case common.RandomEventType_Double_Point:
			return rollWithDoublePoint(roster)
		case common.RandomEventType_CRAZY_THURSDAY:
			return rollWithCrazyThursday(roster)
		default:
			// 没有事件，使用最基本的随机点名，尚未实现
			return common.RosterItem{}, nil
		}
	case common.RollCallMode_SEQUENTIAL:
		return *roster[0], nil
	case common.RollCallMode_REVERSE_SEQUENTIAL:
		return *roster[len(roster)-1], nil
	case common.RollCallMode_LOW_POINTS_FIRST:
		// 积分低优先，按照正态分布处理，尚未实现
		return common.RosterItem{}, nil
	default:
		// 未知点名模式，采用最基本的随机点名，尚未实现
		return common.RosterItem{}, nil
	}
}

func rollWithDoublePoint(roster []*common.RosterItem) (common.RosterItem, error) {
	return common.RosterItem{}, nil
}

func rollWithCrazyThursday(roster []*common.RosterItem) (common.RosterItem, error) {
	return common.RosterItem{}, nil
}
