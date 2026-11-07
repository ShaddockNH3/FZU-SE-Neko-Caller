// 随机事件与积分预览规则。
//
// 原来这套规则在三处各写一遍：模板里的预览标签、submitSolve 里算历史记录、
// 还有后端 biz/domain/rollcall/score.go。这里收成一份，前端的预览和结算
// 都从这里取；后端的结算逻辑以后端为准（前端只做展示）。

export const EVENT_OPTIONS = [
  { value: 0, label: '无事件' },
  {
    value: 1,
    label: '双倍积分',
    tip: {
      title: '双倍积分 (Double Point)',
      body: '触发简单的幸运事件，被点名同学若回答正确，本次所获基础积分翻倍。',
      keywords: '运气，翻倍奖励'
    }
  },
  {
    value: 2,
    label: '疯狂星期四',
    tip: {
      title: '疯狂星期四 (Crazy Thursday)',
      body: '日期限定彩蛋。检测系统时间是否为星期四，触发「V我50」相关机制（如积分为50因数的同学权重增加，或获得特殊Buff）。',
      keywords: '日期检测，玩梗，特定人群Buff'
    }
  },
  {
    value: 3,
    label: '1024福报',
    tip: {
      title: '1024 程序员福报 (Blessing 1024)',
      body: '当点名发生的系统时间（秒）为 2 的幂次方（1, 2, 4, 8...）时触发。回答积分固定为 1.024 分，且答错不扣分。',
      keywords: '二进制，免死金牌，极客文化'
    }
  },
  {
    value: 4,
    label: '质数的孤独',
    tip: {
      title: '质数的孤独 (Solitude of Primes)',
      body: '当点名发生的系统时间（分）为质数（2, 3, 5, 7...）时触发，额外加成 0.37 分。',
      keywords: '数学规律，额外津贴，人文关怀'
    }
  },
  {
    value: 5,
    label: '幸运7大奖',
    tip: {
      title: '幸运 7 大奖 (Lucky 7 Jackpot)',
      body: '基于学号的被动触发事件。若被选中的同学学号末位是 7，回答错误不扣分。',
      keywords: '字符串匹配，豁免权，强运'
    }
  }
]

export const ANSWER_TYPE = Object.freeze({
  NORMAL: 0,
  HELP: 1,
  SKIP: 2,
  TRANSFER: 3
})

const round3 = (value) => Math.round(value * 1000) / 1000

/**
 * 一次作答的积分变化，和后端 ScoreDelta 保持同一套规则。
 * 注意：最终以服务端结算为准，这里只用来显示预览和写本地历史。
 */
export function scoreDelta({ answerType, eventType, customScore = 0 }) {
  if (answerType === ANSWER_TYPE.NORMAL) {
    let delta = customScore

    if (eventType === 1) delta *= 2
    else if (eventType === 2) delta *= 1.5
    else if (eventType === 3) delta = customScore > 0 ? 1.024 : 0
    else if (eventType === 4 && customScore >= 0) delta += 0.37
    else if (eventType === 5 && customScore < 0) delta = 0

    return round3(delta)
  }

  if (answerType === ANSWER_TYPE.HELP) {
    let delta = 0.5
    if (eventType === 1) delta *= 2
    else if (eventType === 2) delta *= 1.5
    else if (eventType === 4) delta += 0.37
    return round3(delta)
  }

  // 跳过、转移都不计分
  return 0
}

/** 结算表单里的那句预览标签。返回 null 表示不显示。 */
export function scorePreview(eventType, customScore) {
  switch (eventType) {
    case 1:
      return { text: `双倍后：${(customScore * 2).toFixed(1)}`, type: 'warning' }
    case 2:
      return { text: `疯四后：${(customScore * 1.5).toFixed(1)}`, type: 'danger' }
    case 3:
      return {
        text: `1024福报：${customScore > 0 ? '1.024' : '0'}${customScore < 0 ? '(免扣分)' : ''}`,
        type: 'success'
      }
    case 4:
      return {
        text: `质数孤独：${customScore >= 0 ? (customScore + 0.37).toFixed(2) : customScore.toFixed(1)}`,
        type: 'primary'
      }
    case 5:
      return { text: `幸运7：${customScore < 0 ? '0(免扣分)' : customScore.toFixed(1)}`, color: '#f56c6c' }
    default:
      return { text: `无事件：${customScore.toFixed(1)}`, type: 'info' }
  }
}

/** 学生卡片上那个事件标签的样式。 */
export function eventBadge(eventType) {
  switch (eventType) {
    case 1:
      return { label: '双倍积分', type: 'warning', tip: '双倍积分: 回答正确积分翻倍' }
    case 2:
      return { label: '疯狂星期四', type: 'danger', tip: '疯狂星期四: V我50特殊机制' }
    case 3:
      return { label: '1024福报', type: 'success', tip: '1024福报: +1.024分且答错免扣分' }
    case 4:
      return { label: '质数的孤独', type: 'primary', tip: '质数的孤独: 回答正确额外+0.37分' }
    case 5:
      return { label: '幸运7大奖', color: '#f56c6c', tip: '幸运7大奖: 答错不扣分' }
    default:
      return null
  }
}
