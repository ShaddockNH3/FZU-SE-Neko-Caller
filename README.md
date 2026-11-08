# FZU-SE-Neko-Caller

福州大学软件工程结对编程作业：课堂随机点名系统，Go（Hertz + GORM）+ Vue 3，支持四种点名模式、积分管理和几个恶趣味随机事件。

## 快速开始

```bash
cd NekoCallerBackend && go run main.go          # :8888，首次启动自动建表
cd NekoCallerFrountend && npm install && npm run dev
```

需要 MySQL，连接信息在 `NekoCallerBackend/biz/dal/mysql/mysql.go`。前端目录名拼错成 `Frountend` 是历史遗留，没改。

## 接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/v1/import/class-data` | 上传 Excel 导入班级与学生 |
| GET | `/v1/classes` | 班级列表（含学生 ID） |
| GET | `/v1/classes/:class_id` | 单个班级 |
| DELETE | `/v1/classes/:class_id` | 删班级并清选课记录 |
| GET | `/v1/roster?class_id=` | 点名名册（学生 + 选课数据） |
| GET | `/v1/classes/:class_id/leaderboard` | 积分排行 |
| GET | `/v1/classes/:class_id/stats` | 班级统计 |
| GET | `/v1/students`、`/v1/students/:id` | 学生列表与详情 |
| DELETE | `/v1/students/:student_id` | 删学生 |
| DELETE | `/v1/enrollments/:enrollment_id` | 把学生踢出班级 |
| POST | `/v1/roll-calls` | 点一个人名，返回选中学生与实际生效的事件 |
| POST | `/v1/roll-calls/solve` | 提交作答，结算积分（支持跳过、转移） |
| POST | `/v1/roll-calls/reset` | 重置班级点名次数与积分 |

## 规则

点名模式和随机事件的规则都在 `biz/domain/rollcall`，纯函数，有单测。

| 模式 | 规则 |
| --- | --- |
| 加权随机 | 权重 `1 / (1 + 积分×0.4 + 被点次数×0.2)`，分高被点多的更难中 |
| 顺序 / 逆序 | 按名册顺序取头/尾 |
| 低分优先 | 在积分最低的三分之一里随机，避免总是同一个人 |

| 事件 | 触发条件 | 结算 |
| --- | --- | --- |
| 双倍积分 | 无条件 | 正常作答或求助，得分 ×2 |
| 疯狂星期四 | 无条件 | 同上，×1.5；点人时积分是 50 的因数权重 ×1.25 |
| 1024 福报 | 当前秒数是 2 的幂 | 答对固定 +1.024，答错免扣 |
| 质数的孤独 | 当前分钟数是质数 | 答对额外 +0.37 |
| 幸运 7 | 学号末位是 7 | 答错免扣 |

每被点名两次多一次转移权，转移权可以把这次作答转给别人（对方照常计分）。

## 项目结构

```
NekoCallerBackend/
├── biz/domain/rollcall/  点名、事件条件、积分结算（纯函数，有单测）
├── biz/service/          用例：读写库 + 调用领域规则
├── biz/handler/          HTTP 入口，Excel 导入
├── biz/dal/              GORM 生成代码与连接
└── pkg/errno/            统一错误码
NekoCallerFrountend/src/
├── theme.js              浅色调色板与 Naive UI themeOverrides，颜色只在这里改
├── utils/feedback.js     全局 message / dialog（组件外也能弹提示）
├── constants/events.js   随机事件文案与积分规则（一份，预览和结算共用）
├── composables/          点名页状态与用例（选班 / 点名 / 结算 / 重置 / 历史）
├── components/           配置面板、事件选择、学生卡片、结算表单、历史
├── views/                页面：首页 / 班级 / 学生 / 点名 / 排行榜 / 花名册
└── api/                  请求封装（错误提示统一在拦截器）
```

## 测试

```bash
cd NekoCallerBackend && go test ./biz/domain/...
```

覆盖四种抽取模式、权重随积分下降、疯狂星期四加权、事件条件（2 的幂 / 质数 / 尾号 7）、积分结算表、转移权发放。

## 已知问题

- 前端没接 CI，改动靠 `npm run build` 和手点验证。
- 前端预演的分数变化和后端结算各算一遍，改规则要同步改 `RollCall.vue`。
- 事件条件看服务器本地时间的秒/分钟，跨时区部署会不对。
- Excel 导入没有按学号去重，同一份表导两次会有重复学生。

## 更多

分层改动与取舍见 [docs/DESIGN.md](docs/DESIGN.md)。

## 致谢

[Hertz](https://github.com/cloudwego/hertz)、[GORM](https://gorm.io/)、[Naive UI](https://www.naiveui.com/)、[excelize](https://github.com/qax-os/excelize)。

## License

MIT，见 [LICENSE](LICENSE)。
