# 代码结构

技术栈没动：Hertz + GORM + Vue 3。动的是规则住在哪儿。

## 原来什么样

| 位置 | 问题 |
| --- | --- |
| `pkg/utils/roll.go` | 包级随机源 `rand.New(rand.NewSource(time.Now()))`，没有注入口子，抽中谁无法复现，测试只能断言"不报错" |
| `biz/service/api_service.go` | 权重、低分优先、事件条件、积分结算全塞在 811 行的 service 里 |
| `checkEventCondition` | 内部直接 `time.Now()`，测试想造"第 16 秒"没法造 |
| `calculateScoreDelta` | 入参是 thrift 请求结构，规则和传输格式绑死 |
| 调试输出 | `fmt.Printf("[DEBUG ...")` 混在业务分支里 |

## 现在怎么分

```
handler（HTTP / Excel）→ service（读写库、调规则）
                              ↓
                     biz/domain/rollcall（纯规则）
```

| 文件 | 内容 |
| --- | --- |
| `picker.go` | `Picker` 持有一个 `*rand.Rand`，`Pick` 按模式抽人；`Weights` / `BaseWeight` / `IsFactorOf50` 可单独断言 |
| `event.go` | `Clock{Minute, Second}` + `EventMet`，时间从外面传进来；`IsPowerOfTwo` / `IsPrime` |
| `score.go` | `ScoreDelta(作答方式, 事件, 自定义分)`，脱开 thrift 指针；`GrantsTransferRight` |

service 只留三个适配点：包级 `picker`、`rollcall.ClockAt(time.Now())`、把 `req` 的三个字段摊给 `ScoreDelta`。

## 为什么把随机源做成字段

原实现每次调用都读全局 `rng`，测试只能断言"没崩"。现在 `NewPicker(seed)` 固定种子就固定结果，`TestPickIsDeterministicWithSeed` 才能成立；线上用 `NewPicker(time.Now().UnixNano())`，行为不变。

## 为什么把时钟做成参数

事件条件判断的是"第几秒 / 第几分"。原实现读 `time.Now()`，只能等真实时间凑巧。抽成 `Clock` 之后，"16 秒满足 1024 福报、17 秒不满足"这种断言一行就能写。

## 前端

`RollCall.vue` 原来 507 行：模板 248 行、脚本 247 行。最要命的是积分规则写了三遍——
后端 `score.go`、模板里的五个预览标签、`submitSolve` 里算历史记录的那段 if/else。

拆完：

| 文件 | 职责 | 行数 |
| --- | --- | --- |
| `views/RollCall.vue` | 两栏版面 + 连事件 | 91 |
| `composables/useRollCall.js` | 全部状态与用例 | 250 |
| `components/RollCallConfig.vue` | 班级、模式、事件、开始/重置 | 130 |
| `components/EventSelector.vue` | 事件单选与说明气泡 | 50 |
| `components/StudentResultCard.vue` | 学生信息 + 事件标签 | 60 |
| `components/SolveForm.vue` | 回答类型、分数滑块、预览、转移目标 | 150 |
| `components/RollCallHistory.vue` | 历史时间线 | 45 |
| `components/StatCard.vue` | 首页四张指标卡（原来每张写一遍） | 70 |

首页 245 行 → 200 行，四张指标卡收成一个组件 + 一个 `statCards` 数组。

积分规则收进 `constants/events.js` 的 `scoreDelta()`，模板预览和本地历史都调它；五个 `el-tooltip` 里的说明文字变成 `EVENT_OPTIONS` 的数据，加事件只改数组。

`scoreDelta` 的 16 种组合我逐条对过后端 `ScoreDelta` 的期望值（正常/求助 × 无事件/双倍/疯四/1024/质数/幸运7/负数分），一致。

## 代价

- `pkg/utils` 包被删掉，逻辑搬到 `biz/domain/rollcall`，引用路径变了。
- service 里多了一层转换（request → 领域入参），行数没省多少。
- 领域层仍然依赖 `biz/model/common` 的 thrift 结构体，不是彻底的纯净领域模型——为了不复制一套影子类型，认了。
- 前端 `RollCall.vue` 里还有一份分数预演逻辑，没合并到共享模块，规则变了要改两处。
