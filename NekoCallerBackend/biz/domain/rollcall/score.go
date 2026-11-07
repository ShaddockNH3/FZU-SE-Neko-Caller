package rollcall

import (
	"math"

	"FZUSENekoCaller/biz/model/common"
)

// ScoreDelta 算一次作答带来的积分变化。
//
// 入参只有三样：作答方式、随机事件、自定义分数，所以能直接列表测。
// 返回保留三位小数，和原来 math.Round(base*1000)/1000 一致。
func ScoreDelta(answer common.AnswerType, event common.RandomEventType, customScore *float64) float64 {
	base, correct := baseScore(answer, customScore)
	base = applyEvent(base, answer, event, customScore, correct)
	return math.Round(base*1000) / 1000
}

// baseScore 是没算事件加成的原始分，第二个返回值表示"算不算答对"。
func baseScore(answer common.AnswerType, customScore *float64) (float64, bool) {
	switch answer {
	case common.AnswerType_NORMAL:
		if customScore == nil {
			return 0, false
		}
		score := *customScore
		return score, score >= 0
	case common.AnswerType_HELP:
		// 请人帮忙复述也算答对，给半分
		return 0.5, true
	case common.AnswerType_SKIP:
		// 有跳过权就不扣分
		return 0, false
	case common.AnswerType_TRANSFER:
		// 有转移权就不扣分
		return 0, false
	default:
		return 0, false
	}
}

// applyEvent 叠随机事件。只有正常作答和求助吃加成，跳过和转移不吃。
func applyEvent(
	base float64,
	answer common.AnswerType,
	event common.RandomEventType,
	customScore *float64,
	correct bool,
) float64 {
	switch event {
	case common.RandomEventType_BLESSING_1024:
		// 答对固定 +1.024，答错免扣
		if answer == common.AnswerType_NORMAL && customScore != nil {
			switch {
			case *customScore > 0:
				return 1.024
			case *customScore < 0:
				return 0
			}
		}
	case common.RandomEventType_SOLITUDE_PRIMES:
		// 答对多给 0.37
		if correct {
			return base + 0.37
		}
	case common.RandomEventType_LUCKY_7:
		// 答错免扣（幸运闪避）
		if answer == common.AnswerType_NORMAL && customScore != nil && *customScore < 0 {
			return 0
		}
	case common.RandomEventType_Double_Point:
		if eatsBonus(answer) {
			return base * 2
		}
	case common.RandomEventType_CRAZY_THURSDAY:
		if eatsBonus(answer) {
			return base * 1.5
		}
	}
	return base
}

func eatsBonus(answer common.AnswerType) bool {
	return answer == common.AnswerType_NORMAL || answer == common.AnswerType_HELP
}

// GrantsTransferRight 判断这次点名后是否多给一次转移权：每被点两次给一次。
func GrantsTransferRight(nextCallCount int64) bool {
	return nextCallCount > 0 && nextCallCount%2 == 0
}
