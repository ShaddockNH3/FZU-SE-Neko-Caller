package rollcall

import (
	"testing"
	"time"

	"FZUSENekoCaller/biz/model/common"
)

func TestIsPowerOfTwo(t *testing.T) {
	for _, value := range []int{1, 2, 4, 8, 16, 32} {
		if !IsPowerOfTwo(value) {
			t.Errorf("%d 是 2 的幂", value)
		}
	}
	for _, value := range []int{0, -2, 3, 6, 17, 60} {
		if IsPowerOfTwo(value) {
			t.Errorf("%d 不是 2 的幂", value)
		}
	}
}

func TestIsPrime(t *testing.T) {
	for _, value := range []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59} {
		if !IsPrime(value) {
			t.Errorf("%d 是质数", value)
		}
	}
	for _, value := range []int{0, 1, 4, 6, 9, 21, 25, 49, 60} {
		if IsPrime(value) {
			t.Errorf("%d 不是质数", value)
		}
	}
}

func TestEventMet(t *testing.T) {
	clock := Clock{Minute: 13, Second: 16}

	cases := []struct {
		name      string
		event     common.RandomEventType
		studentID string
		clock     Clock
		want      bool
	}{
		{"普通事件不设条件", common.RandomEventType_NONE, "20230001", clock, true},
		{"1024 福报：16 秒满足", common.RandomEventType_BLESSING_1024, "20230001", Clock{Second: 16}, true},
		{"1024 福报：17 秒不满足", common.RandomEventType_BLESSING_1024, "20230001", Clock{Second: 17}, false},
		{"质数的孤独：13 分满足", common.RandomEventType_SOLITUDE_PRIMES, "20230001", Clock{Minute: 13}, true},
		{"质数的孤独：14 分不满足", common.RandomEventType_SOLITUDE_PRIMES, "20230001", Clock{Minute: 14}, false},
		{"幸运 7：末位是 7", common.RandomEventType_LUCKY_7, "20230007", clock, true},
		{"幸运 7：末位不是 7", common.RandomEventType_LUCKY_7, "20230001", clock, false},
		{"幸运 7：空学号不满足", common.RandomEventType_LUCKY_7, "", clock, false},
	}

	for _, testCase := range cases {
		got := EventMet(testCase.event, testCase.studentID, testCase.clock)
		if got != testCase.want {
			t.Errorf("%s: 期望 %v，实际 %v", testCase.name, testCase.want, got)
		}
	}
}

func TestClockAtUsesGivenLocation(t *testing.T) {
	shanghai, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Skipf("找不到时区数据: %v", err)
	}

	// UTC 04:23:10 == 北京时间 12:23:10，东八区是整小时偏移，分秒不变
	clock := ClockAt(time.Date(2025, 5, 1, 4, 23, 10, 0, time.UTC), shanghai)
	if clock.Minute != 23 || clock.Second != 10 {
		t.Errorf("时钟取值不对: %+v", clock)
	}

	// 加德满都 +05:45，分钟会真的变：UTC 04:23:10 → 当地 10:08:10
	kathmandu, err := time.LoadLocation("Asia/Kathmandu")
	if err != nil {
		t.Skipf("找不到时区数据: %v", err)
	}

	clock = ClockAt(time.Date(2025, 5, 1, 4, 23, 10, 0, time.UTC), kathmandu)
	if clock.Minute != 8 || clock.Second != 10 {
		t.Errorf("跨时区取分钟不对: %+v", clock)
	}
}

// loc 传 nil 时退回服务器本地时区，不应该 panic
func TestClockAtNilLocation(t *testing.T) {
	clock := ClockAt(time.Date(2025, 5, 1, 12, 23, 10, 0, time.Local), nil)
	if clock.Minute != 23 || clock.Second != 10 {
		t.Errorf("时钟取值不对: %+v", clock)
	}
}
