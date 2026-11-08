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

## 时钟与时区

`ClockAt(now, loc)`：先把时间换到 `loc`，再取分钟和秒。`loc` 来自 `pkg/constants.AppLocation`，
默认 `Asia/Shanghai`，可用 `APP_TIMEZONE` 覆盖。

为什么较真：`质数的孤独` 看分钟、`1024 福报` 看秒，都是学生眼里的"现在几点"。
服务器跑在 UTC 时东八区是整小时偏移，分秒恰好一样，看不出问题；但只要部署到 +05:30 这种
半小时时区，同一瞬间的分钟就不同，事件会在错误的时间触发。现在不再依赖宿主机时区设置。

`pkg/constants` 里 `import _ "time/tzdata"`，把时区库编进二进制——容器镜像常常只带 UTC，
没有这个文件 `LoadLocation` 会失败然后悄悄退回 UTC。

## 为什么把时钟做成参数

事件条件判断的是"第几秒 / 第几分"。原实现读 `time.Now()`，只能等真实时间凑巧。抽成 `Clock` 之后，"16 秒满足 1024 福报、17 秒不满足"这种断言一行就能写。

## 前端

`RollCall.vue` 原来 507 行：模板 248 行、脚本 247 行。最要命的是积分规则写了三遍——
后端 `score.go`、模板里的五个预览标签、`submitSolve` 里算历史记录的那段 if/else。

拆完：

| 文件 | 职责 | 行数 |
| --- | --- | --- |
| `views/RollCall.vue` | 两栏版面 + 连事件 | 98 |
| `composables/useRollCall.js` | 全部状态与用例 | 252 |
| `components/RollCallConfig.vue` | 班级、模式、事件、开始/重置 | 131 |
| `components/EventSelector.vue` | 事件单选与说明气泡 | 55 |
| `components/StudentResultCard.vue` | 学生信息 + 事件标签 | 67 |
| `components/SolveForm.vue` | 回答类型、分数滑块、预览、转移目标 | 162 |
| `components/RollCallHistory.vue` | 历史时间线 | 56 |
| `components/StatCard.vue` | 首页四张指标卡（原来每张写一遍） | 54 |

首页 245 行 → 195 行，四张指标卡收成一个组件 + 一个 `statCards` 数组。

### 界面

这一版把 Element Plus 换成了 **Naive UI + Phosphor 图标**，和港区那个项目用同一套色值
（底 `#f4f6fa`、卡片纯白、主色深蓝 `#1e40af`、强调金 `#c8963e`），全部写在 `theme.js` 的
`themeOverrides` 里，组件内不写死颜色。

侧栏从深灰 200px 换成白色可折叠：< 720px 收成 64px 图标栏。

顺带修的：

| 问题 | 处理 |
| --- | --- |
| 首页「今日点名 / 累计点名」是编的（学生数 × 5、再乘 0.1） | 改成按班级拉排行榜，累加 `call_count` 与 `total_points`；两张卡换成「累计点名」「平均积分」 |
| 进点名页班级是空的，按钮点不了 | `loadClasses` 里默认选第一个，从首页带 `classId` 进来优先用它 |
| Element 的标签类型 `danger` 在 Naive 里叫 `error` | `constants/events.js` 加 `naiveTagType()` 统一转 |
| 提示同时在拦截器和 catch 里弹两次 | 收到拦截器一处，其余只写日志 |

`累计点名` 现在是按班级逐个拉排行榜算的，N 个班就是 N 次请求。几十个班以内没问题，
再多就该让后端加一个汇总接口（`GetClassStats` 只到单班粒度）。

积分规则收进 `constants/events.js` 的 `scoreDelta()`，模板预览和本地历史都调它；五个 `el-tooltip` 里的说明文字变成 `EVENT_OPTIONS` 的数据，加事件只改数组。

`scoreDelta` 的 16 种组合我逐条对过后端 `ScoreDelta` 的期望值（正常/求助 × 无事件/双倍/疯四/1024/质数/幸运7/负数分），一致。

## 后端构建

`go 1.27.0`。之前写的「必须用 Go 1.24」不成立了：卡点是 `sonic v1.14` 用了 Go 1.27 已删掉的
`GoMapIterator`，把 hertz 升到 0.10.6、sonic 升到 1.15.5 就通了，代码一行没改。

## 代价

- `pkg/utils` 包被删掉，逻辑搬到 `biz/domain/rollcall`，引用路径变了。
- service 里多了一层转换（request → 领域入参），行数没省多少。
- 领域层仍然依赖 `biz/model/common` 的 thrift 结构体，不是彻底的纯净领域模型——为了不复制一套影子类型，认了。
- 前端 `RollCall.vue` 里还有一份分数预演逻辑，没合并到共享模块，规则变了要改两处。
