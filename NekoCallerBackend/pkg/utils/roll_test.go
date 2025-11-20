package utils

import (
	"testing"

	"FZUSENekoCaller/biz/model/common"
)

// TestRollWithEmptyRoster 测试空名册边界情况
func TestRollWithEmptyRoster(t *testing.T) {
	roster := []*common.RosterItem{}
	_, err := Roll(roster, common.RollCallMode_RANDOM, common.RandomEventType_NONE)
	if err == nil {
		t.Error("Expected error for empty roster, got nil")
	}
}

// TestBaseWeight 测试基础权重计算
func TestBaseWeight(t *testing.T) {
	tests := []struct {
		name      string
		points    float64
		callCount int64
		want      float64
	}{
		{
			name:      "零积分零次数",
			points:    0,
			callCount: 0,
			want:      1.0, // 1/(1+0+0) = 1
		},
		{
			name:      "高积分学生",
			points:    10,
			callCount: 0,
			want:      0.2, // 1/(1+10*0.4+0) = 1/5 = 0.2
		},
		{
			name:      "高点名次数学生",
			points:    0,
			callCount: 5,
			want:      0.5, // 1/(1+0+5*0.2) = 1/2 = 0.5
		},
		{
			name:      "高积分高次数",
			points:    10,
			callCount: 5,
			want:      0.167, // 1/(1+4+1) ≈ 0.167
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := &common.RosterItem{
				EnrollmentInfo: &common.Enrollment{
					TotalPoints: tt.points,
					CallCount:   tt.callCount,
				},
			}
			got := baseWeight(item)
			// 允许0.01的误差
			if got < tt.want-0.01 || got > tt.want+0.01 {
				t.Errorf("baseWeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestLowPointsFirst 测试低分优先算法
func TestLowPointsFirst(t *testing.T) {
	// 构造测试数据：5个学生，积分从0到40
	roster := []*common.RosterItem{
		{
			StudentInfo:    &common.Student{StudentID: "001", Name: "学生1"},
			EnrollmentInfo: &common.Enrollment{TotalPoints: 0, CallCount: 0},
		},
		{
			StudentInfo:    &common.Student{StudentID: "002", Name: "学生2"},
			EnrollmentInfo: &common.Enrollment{TotalPoints: 10, CallCount: 1},
		},
		{
			StudentInfo:    &common.Student{StudentID: "003", Name: "学生3"},
			EnrollmentInfo: &common.Enrollment{TotalPoints: 20, CallCount: 2},
		},
		{
			StudentInfo:    &common.Student{StudentID: "004", Name: "学生4"},
			EnrollmentInfo: &common.Enrollment{TotalPoints: 30, CallCount: 3},
		},
		{
			StudentInfo:    &common.Student{StudentID: "005", Name: "学生5"},
			EnrollmentInfo: &common.Enrollment{TotalPoints: 40, CallCount: 4},
		},
	}

	// 运行100次，确保只从前1/3中选择（即前2个学生）
	validStudents := map[string]bool{"001": true, "002": true}
	for i := 0; i < 100; i++ {
		result := lowPointsFirst(roster)
		if !validStudents[result.StudentInfo.StudentID] {
			t.Errorf("lowPointsFirst() selected student %s, should only select from top 1/3 (001 or 002)",
				result.StudentInfo.StudentID)
		}
	}
}

// TestIsFactorOf50 测试疯狂星期四因数判断
func TestIsFactorOf50(t *testing.T) {
	tests := []struct {
		points float64
		want   bool
	}{
		{1, true},   // 50的因数
		{2, true},   // 50的因数
		{5, true},   // 50的因数
		{10, true},  // 50的因数
		{25, true},  // 50的因数
		{50, true},  // 50的因数
		{3, false},  // 不是50的因数
		{7, false},  // 不是50的因数
		{0, false},  // 边界值
		{-1, false}, // 负数
		{51, false}, // 超过50
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := isFactorOf50(tt.points)
			if got != tt.want {
				t.Errorf("isFactorOf50(%v) = %v, want %v", tt.points, got, tt.want)
			}
		})
	}
}

// TestWeightedRoll 测试加权随机点名
func TestWeightedRoll(t *testing.T) {
	roster := []*common.RosterItem{
		{
			StudentInfo:    &common.Student{StudentID: "001", Name: "低分学生"},
			EnrollmentInfo: &common.Enrollment{TotalPoints: 0, CallCount: 0},
		},
		{
			StudentInfo:    &common.Student{StudentID: "002", Name: "高分学生"},
			EnrollmentInfo: &common.Enrollment{TotalPoints: 50, CallCount: 10},
		},
	}

	// 运行1000次统计，低分学生应该被选中更多次
	count := make(map[string]int)
	for i := 0; i < 1000; i++ {
		result, err := Roll(roster, common.RollCallMode_RANDOM, common.RandomEventType_NONE)
		if err != nil {
			t.Fatalf("Roll() error = %v", err)
		}
		count[result.StudentInfo.StudentID]++
	}

	// 低分学生应该被选中更多次（至少60%的概率）
	if count["001"] < 600 {
		t.Errorf("Low points student selected %d times, expected at least 600 times in 1000 runs", count["001"])
	}

	t.Logf("Selection distribution: 低分学生=%d次, 高分学生=%d次", count["001"], count["002"])
}

// TestRollModes 测试不同点名模式
func TestRollModes(t *testing.T) {
	roster := []*common.RosterItem{
		{
			StudentInfo:    &common.Student{StudentID: "001", Name: "学生1"},
			EnrollmentInfo: &common.Enrollment{TotalPoints: 10, CallCount: 1},
		},
		{
			StudentInfo:    &common.Student{StudentID: "002", Name: "学生2"},
			EnrollmentInfo: &common.Enrollment{TotalPoints: 20, CallCount: 2},
		},
		{
			StudentInfo:    &common.Student{StudentID: "003", Name: "学生3"},
			EnrollmentInfo: &common.Enrollment{TotalPoints: 30, CallCount: 3},
		},
	}

	// 测试顺序点名
	result, err := Roll(roster, common.RollCallMode_SEQUENTIAL, common.RandomEventType_NONE)
	if err != nil {
		t.Fatalf("Sequential roll error: %v", err)
	}
	if result.StudentInfo.StudentID != "001" {
		t.Errorf("Sequential mode should select first student, got %s", result.StudentInfo.StudentID)
	}

	// 测试逆序点名
	result, err = Roll(roster, common.RollCallMode_REVERSE_SEQUENTIAL, common.RandomEventType_NONE)
	if err != nil {
		t.Fatalf("Reverse sequential roll error: %v", err)
	}
	if result.StudentInfo.StudentID != "003" {
		t.Errorf("Reverse sequential mode should select last student, got %s", result.StudentInfo.StudentID)
	}
}
