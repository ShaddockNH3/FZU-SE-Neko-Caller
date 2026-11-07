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

## 代价

- `pkg/utils` 包被删掉，逻辑搬到 `biz/domain/rollcall`，引用路径变了。
- service 里多了一层转换（request → 领域入参），行数没省多少。
- 领域层仍然依赖 `biz/model/common` 的 thrift 结构体，不是彻底的纯净领域模型——为了不复制一套影子类型，认了。
- 前端 `RollCall.vue` 里还有一份分数预演逻辑，没合并到共享模块，规则变了要改两处。
