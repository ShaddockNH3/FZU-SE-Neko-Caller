package rollcall

import (
	"testing"

	"FZUSENekoCaller/biz/model/common"
)

func roster() []*common.RosterItem {
	return []*common.RosterItem{
		{EnrollmentInfo: &common.Enrollment{TotalPoints: 0, CallCount: 0}},
		{EnrollmentInfo: &common.Enrollment{TotalPoints: 10, CallCount: 0}},
		{EnrollmentInfo: &common.Enrollment{TotalPoints: 0, CallCount: 3}},
	}
}

func TestPickWithEmptyRoster(t *testing.T) {
	picker := NewPicker(1)
	if _, err := picker.Pick(nil, common.RollCallMode_RANDOM, common.RandomEventType_NONE); err == nil {
		t.Error("空名册应该报错")
	}
}

func TestBaseWeight(t *testing.T) {
	cases := []struct {
		name      string
		points    float64
		callCount int64
		want      float64
	}{
		{"零分零次", 0, 0, 1.0},
		{"高分压低权重", 10, 0, 0.2},
		{"被点多次也压权重", 0, 5, 0.5},
		{"负分按 0 算", -3, 0, 1.0},
	}

	for _, testCase := range cases {
		item := &common.RosterItem{
			EnrollmentInfo: &common.Enrollment{TotalPoints: testCase.points, CallCount: testCase.callCount},
		}
		got := BaseWeight(item)
		diff := got - testCase.want
		if diff > 1e-9 || diff < -1e-9 {
			t.Errorf("%s: 期望 %v，实际 %v", testCase.name, testCase.want, got)
		}
	}
}

func TestWeightsBoostOnEvents(t *testing.T) {
	base := Weights(roster(), common.RandomEventType_NONE)
	doubled := Weights(roster(), common.RandomEventType_Double_Point)

	for index := range base {
		if doubled[index] <= base[index] {
			t.Errorf("双倍积分应该抬高 %d 号权重", index)
		}
	}

	// 积分 25 是 50 的因数，疯狂星期四该被加成
	lucky := []*common.RosterItem{
		{EnrollmentInfo: &common.Enrollment{TotalPoints: 25}},
		{EnrollmentInfo: &common.Enrollment{TotalPoints: 7}},
	}
	boosted := Weights(lucky, common.RandomEventType_CRAZY_THURSDAY)
	plain := Weights(lucky, common.RandomEventType_NONE)
	if boosted[0] <= plain[0] {
		t.Error("积分 25 在疯狂星期四应该被加成")
	}
	if boosted[1] != plain[1] {
		t.Error("积分 7 不是 50 的因数，不该被加成")
	}
}

func TestPickModes(t *testing.T) {
	items := roster()
	picker := NewPicker(42)

	first, err := picker.Pick(items, common.RollCallMode_SEQUENTIAL, common.RandomEventType_NONE)
	if err != nil || first != *items[0] {
		t.Error("顺序模式应该抽第一个")
	}

	last, err := picker.Pick(items, common.RollCallMode_REVERSE_SEQUENTIAL, common.RandomEventType_NONE)
	if err != nil || last != *items[len(items)-1] {
		t.Error("逆序模式应该抽最后一个")
	}
}

func TestLowPointsFirstPicksFromBottomThird(t *testing.T) {
	items := make([]*common.RosterItem, 0, 9)
	for index := 0; index < 9; index++ {
		items = append(items, &common.RosterItem{
			EnrollmentInfo: &common.Enrollment{
				StudentID:   string(rune('A' + index)),
				TotalPoints: float64(index * 10),
				CallCount:   0,
			},
		})
	}

	// 固定 seed 下只该在最低的三个人里挑
	for seed := int64(0); seed < 20; seed++ {
		picked, err := NewPicker(seed).Pick(items, common.RollCallMode_LOW_POINTS_FIRST, common.RandomEventType_NONE)
		if err != nil {
			t.Fatalf("不该报错: %v", err)
		}
		if picked.EnrollmentInfo.TotalPoints > 20 {
			t.Errorf("seed=%d 抽到了 %v 分的学生，超出最低三分之一",
				seed, picked.EnrollmentInfo.TotalPoints)
		}
	}
}

func TestPickIsDeterministicWithSeed(t *testing.T) {
	items := roster()
	a, _ := NewPicker(7).Pick(items, common.RollCallMode_RANDOM, common.RandomEventType_NONE)
	b, _ := NewPicker(7).Pick(items, common.RollCallMode_RANDOM, common.RandomEventType_NONE)
	if a != b {
		t.Error("同一个 seed 应该出同一个人")
	}
}

func TestIsFactorOf50(t *testing.T) {
	yes := []float64{1, 2, 5, 10, 25, 50}
	no := []float64{0, -5, 3, 7, 51, 100}

	for _, value := range yes {
		if !IsFactorOf50(value) {
			t.Errorf("%v 应该是 50 的因数", value)
		}
	}
	for _, value := range no {
		if IsFactorOf50(value) {
			t.Errorf("%v 不该算 50 的因数", value)
		}
	}
}
