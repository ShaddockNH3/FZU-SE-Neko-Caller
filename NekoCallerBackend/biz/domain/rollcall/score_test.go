package rollcall

import (
	"testing"

	"FZUSENekoCaller/biz/model/common"
)

func ptr(value float64) *float64 { return &value }

func TestScoreDelta(t *testing.T) {
	cases := []struct {
		name   string
		answer common.AnswerType
		event  common.RandomEventType
		custom *float64
		want   float64
	}{
		{"正常作答用自定义分", common.AnswerType_NORMAL, common.RandomEventType_NONE, ptr(2), 2},
		{"答错是负分", common.AnswerType_NORMAL, common.RandomEventType_NONE, ptr(-1), -1},
		{"没传分数按 0 算", common.AnswerType_NORMAL, common.RandomEventType_NONE, nil, 0},
		{"求助固定半分", common.AnswerType_HELP, common.RandomEventType_NONE, nil, 0.5},
		{"跳过不扣分", common.AnswerType_SKIP, common.RandomEventType_NONE, nil, 0},
		{"转移不扣分", common.AnswerType_TRANSFER, common.RandomEventType_NONE, nil, 0},

		{"1024 福报答对固定 1.024", common.AnswerType_NORMAL, common.RandomEventType_BLESSING_1024, ptr(3), 1.024},
		{"1024 福报答错免扣", common.AnswerType_NORMAL, common.RandomEventType_BLESSING_1024, ptr(-2), 0},
		{"质数的孤独答对加 0.37", common.AnswerType_NORMAL, common.RandomEventType_SOLITUDE_PRIMES, ptr(1), 1.37},
		{"质数的孤独答错不加", common.AnswerType_NORMAL, common.RandomEventType_SOLITUDE_PRIMES, ptr(-1), -1},
		{"幸运 7 答错免扣", common.AnswerType_NORMAL, common.RandomEventType_LUCKY_7, ptr(-3), 0},
		{"幸运 7 答对照常", common.AnswerType_NORMAL, common.RandomEventType_LUCKY_7, ptr(1), 1},
		{"双倍积分翻倍", common.AnswerType_NORMAL, common.RandomEventType_Double_Point, ptr(1), 2},
		{"疯狂星期四 1.5 倍", common.AnswerType_NORMAL, common.RandomEventType_CRAZY_THURSDAY, ptr(2), 3},
		{"求助也吃双倍", common.AnswerType_HELP, common.RandomEventType_Double_Point, nil, 1},
		{"跳过不吃双倍", common.AnswerType_SKIP, common.RandomEventType_Double_Point, nil, 0},
		{"结果保留三位小数", common.AnswerType_NORMAL, common.RandomEventType_CRAZY_THURSDAY, ptr(0.1111), 0.167},
	}

	for _, testCase := range cases {
		got := ScoreDelta(testCase.answer, testCase.event, testCase.custom)
		if got != testCase.want {
			t.Errorf("%s: 期望 %v，实际 %v", testCase.name, testCase.want, got)
		}
	}
}

func TestGrantsTransferRight(t *testing.T) {
	cases := map[int64]bool{
		0: false,
		1: false,
		2: true,
		3: false,
		4: true,
		5: false,
	}

	for callCount, want := range cases {
		if got := GrantsTransferRight(callCount); got != want {
			t.Errorf("被点 %d 次: 期望 %v，实际 %v", callCount, want, got)
		}
	}
}
