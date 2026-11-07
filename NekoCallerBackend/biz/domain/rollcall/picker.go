// Package rollcall 放点名系统的纯规则：怎么抽人、事件条件怎么判、积分怎么算。
//
// 这些逻辑原来混在 biz/service/api_service.go 里，靠 time.Now() 和包级随机源，
// 想测都没法测。搬出来之后随机源和时间都是参数。
package rollcall

import (
	"math"
	"math/rand"
	"sort"

	"FZUSENekoCaller/biz/model/common"
	"FZUSENekoCaller/pkg/errno"
)

// Picker 是带随机源的点名器。同一个 seed 出同一个结果，方便测试和复现。
type Picker struct {
	rng *rand.Rand
}

func NewPicker(seed int64) *Picker {
	return &Picker{rng: rand.New(rand.NewSource(seed))}
}

// NewPickerWith 传入外部随机源，服务层用全局那一个，测试用固定 seed 的。
func NewPickerWith(rng *rand.Rand) *Picker {
	return &Picker{rng: rng}
}

// Pick 按模式抽一个学生。
func (p *Picker) Pick(roster []*common.RosterItem, mode common.RollCallMode, event common.RandomEventType) (common.RosterItem, error) {
	if len(roster) == 0 {
		return common.RosterItem{}, errno.StudentNotFoundErr
	}

	switch mode {
	case common.RollCallMode_SEQUENTIAL:
		return *roster[0], nil
	case common.RollCallMode_REVERSE_SEQUENTIAL:
		return *roster[len(roster)-1], nil
	case common.RollCallMode_LOW_POINTS_FIRST:
		return p.lowPointsFirst(roster), nil
	case common.RollCallMode_RANDOM:
		return p.weightedPick(roster, event)
	default:
		return p.weightedPick(roster, event)
	}
}

// weightedPick 按权重抽：分低、被点少的权重大；事件会额外加权。
func (p *Picker) weightedPick(roster []*common.RosterItem, event common.RandomEventType) (common.RosterItem, error) {
	weights := Weights(roster, event)
	total := 0.0
	for _, weight := range weights {
		total += weight
	}
	if total <= 0 {
		return *roster[p.rng.Intn(len(roster))], nil
	}

	point := p.rng.Float64() * total
	for index, weight := range weights {
		if point <= weight {
			return *roster[index], nil
		}
		point -= weight
	}
	return *roster[len(roster)-1], nil
}

// lowPointsFirst 在积分最低的三分之一里随机——纯粹按最低分抽会让同一个学生反复中招。
func (p *Picker) lowPointsFirst(roster []*common.RosterItem) common.RosterItem {
	sorted := make([]*common.RosterItem, len(roster))
	copy(sorted, roster)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].EnrollmentInfo.TotalPoints == sorted[j].EnrollmentInfo.TotalPoints {
			return sorted[i].EnrollmentInfo.CallCount < sorted[j].EnrollmentInfo.CallCount
		}
		return sorted[i].EnrollmentInfo.TotalPoints < sorted[j].EnrollmentInfo.TotalPoints
	})
	limit := int(math.Max(1, math.Ceil(float64(len(sorted))/3)))
	return *sorted[p.rng.Intn(limit)]
}

// Weights 返回每个学生的权重，暴露出来是为了能直接断言"分越高权重越小"。
func Weights(roster []*common.RosterItem, event common.RandomEventType) []float64 {
	weights := make([]float64, len(roster))
	for index, item := range roster {
		weight := BaseWeight(item)
		switch event {
		case common.RandomEventType_Double_Point:
			// 双倍积分只是乘分数，不动中签率；这里保留原实现，避免改变手感
			weight *= 1.3
		case common.RandomEventType_CRAZY_THURSDAY:
			if IsFactorOf50(item.EnrollmentInfo.TotalPoints) {
				weight *= 1.25
			}
		}
		weights[index] = weight
	}
	return weights
}

// BaseWeight 分越高、被点越多，权重越低。
func BaseWeight(item *common.RosterItem) float64 {
	points := math.Max(item.EnrollmentInfo.TotalPoints, 0)
	callCount := math.Max(float64(item.EnrollmentInfo.CallCount), 0)
	return 1 / (1 + points*0.4 + callCount*0.2)
}

// IsFactorOf50 判断积分是不是 50 的因数（1、2、5、10、25、50），疯狂星期四用。
func IsFactorOf50(points float64) bool {
	intPoints := int(points)
	if intPoints <= 0 || intPoints > 50 {
		return false
	}
	return 50%intPoints == 0
}
