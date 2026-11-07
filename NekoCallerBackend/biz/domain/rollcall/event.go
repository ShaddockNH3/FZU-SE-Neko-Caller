package rollcall

import (
	"time"

	"FZUSENekoCaller/biz/model/common"
)

// Clock 是判事件条件用到的"现在几点"。抽出来是为了测试能随便捏时间。
type Clock struct {
	Minute int
	Second int
}

// ClockAt 把真实时间截到分钟和秒。
func ClockAt(now time.Time) Clock {
	return Clock{Minute: now.Minute(), Second: now.Second()}
}

// EventMet 判断随机事件的条件是否满足。
//
// 不满足时不生效（调用方会把事件降级成 NONE），已有行为没变，
// 只是把判断从"读 time.Now()"改成"读传进来的时钟"。
func EventMet(event common.RandomEventType, studentID string, clock Clock) bool {
	switch event {
	case common.RandomEventType_BLESSING_1024:
		// 1024 福报：秒数是 2 的幂
		return IsPowerOfTwo(clock.Second)
	case common.RandomEventType_SOLITUDE_PRIMES:
		// 质数的孤独：分钟数是质数
		return IsPrime(clock.Minute)
	case common.RandomEventType_LUCKY_7:
		// 幸运 7：学号末位是 7
		return len(studentID) > 0 && studentID[len(studentID)-1] == '7'
	default:
		return true
	}
}

// IsPowerOfTwo 判断正整数是不是 2 的幂。
func IsPowerOfTwo(n int) bool {
	if n <= 0 {
		return false
	}
	return (n & (n - 1)) == 0
}

// IsPrime 判断 0..59 之间的数是不是质数（只用得到分钟，表比试除好读）。
func IsPrime(n int) bool {
	switch n {
	case 2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59:
		return true
	default:
		return false
	}
}
