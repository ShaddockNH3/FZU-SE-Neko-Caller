package utils

import (
	"math"
	"math/rand"
	"sort"
	"time"

	"FZUSENekoCaller/biz/model/common"
	"FZUSENekoCaller/pkg/errno"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func Roll(roster []*common.RosterItem, mode common.RollCallMode, event common.RandomEventType) (common.RosterItem, error) {
	if len(roster) == 0 {
		return common.RosterItem{}, errno.StudentNotFoundErr
	}

	switch mode {
	case common.RollCallMode_RANDOM:
		return randomRoll(roster, event)
	case common.RollCallMode_SEQUENTIAL:
		return *roster[0], nil
	case common.RollCallMode_REVERSE_SEQUENTIAL:
		return *roster[len(roster)-1], nil
	case common.RollCallMode_LOW_POINTS_FIRST:
		return lowPointsFirst(roster), nil
	default:
		return randomRoll(roster, event)
	}
}

func randomRoll(roster []*common.RosterItem, event common.RandomEventType) (common.RosterItem, error) {
	weights := make([]float64, len(roster))
	total := 0.0
	for i, item := range roster {
		w := baseWeight(item)
		switch event {
		case common.RandomEventType_Double_Point:
			w *= 1.3
		case common.RandomEventType_CRAZY_THURSDAY:
			// 疯狂星期四: 积分为50的因数的学生权重增加
			if isFactorOf50(item.EnrollmentInfo.TotalPoints) {
				w *= 1.25
			}
		}
		weights[i] = w
		total += w
	}

	if total <= 0 {
		return *roster[rng.Intn(len(roster))], nil
	}

	randPoint := rng.Float64() * total
	for i, weight := range weights {
		if randPoint <= weight {
			return *roster[i], nil
		}
		randPoint -= weight
	}

	return *roster[len(roster)-1], nil
}

func lowPointsFirst(roster []*common.RosterItem) common.RosterItem {
	copyRoster := make([]*common.RosterItem, len(roster))
	copy(copyRoster, roster)
	sort.Slice(copyRoster, func(i, j int) bool {
		if copyRoster[i].EnrollmentInfo.TotalPoints == copyRoster[j].EnrollmentInfo.TotalPoints {
			return copyRoster[i].EnrollmentInfo.CallCount < copyRoster[j].EnrollmentInfo.CallCount
		}
		return copyRoster[i].EnrollmentInfo.TotalPoints < copyRoster[j].EnrollmentInfo.TotalPoints
	})
	limit := int(math.Max(1, math.Ceil(float64(len(copyRoster))/3)))
	return *copyRoster[rng.Intn(limit)]
}

func baseWeight(item *common.RosterItem) float64 {
	points := math.Max(item.EnrollmentInfo.TotalPoints, 0)
	callCount := math.Max(float64(item.EnrollmentInfo.CallCount), 0)
	return 1 / (1 + points*0.4 + callCount*0.2)
}

func rollWithDoublePoint(roster []*common.RosterItem) (common.RosterItem, error) {
	return randomRoll(roster, common.RandomEventType_Double_Point)
}

func rollWithCrazyThursday(roster []*common.RosterItem) (common.RosterItem, error) {
	return randomRoll(roster, common.RandomEventType_CRAZY_THURSDAY)
}

// isFactorOf50 检查一个积分值是否为50的因数
// 50的因数有: 1, 2, 5, 10, 25, 50
func isFactorOf50(points float64) bool {
	// 取整数部分
	intPoints := int(points)
	if intPoints <= 0 || intPoints > 50 {
		return false
	}
	// 检查是否为50的因数
	return 50%intPoints == 0
}
