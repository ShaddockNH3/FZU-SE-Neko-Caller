# Neko Caller Architecture & Delivery Plan

## 1. Goals & Scope
- 实现符合《0.作业要求文档》全部基础需求，并提供至少 2 个进阶玩法（点名转移权、节日事件加成）。
- 交付一套“前端 + Go 后端 + MySQL”课堂点名系统，支持 Excel 导入/导出、随机&顺序点名、积分排名可视化。
- 输出配套材料：原型文档、类图/流程图、PSP 统计、学习进度条、博客草稿结构等，方便后续整理博文与演示视频。

## 2. 技术栈
| 层级 | 方案 | 说明 |
| --- | --- | --- |
| 前端 | React 18 + Vite + TypeScript + Ant Design + ECharts | 桌面 Web 端，支持点名控制台、排行榜、可视化面板。
| 后端 | Go 1.22 + CloudWeGo Hertz + GORM + Wire (DI) | 提供 RESTful API、Excel 解析、RollCall 算法、数据导出。
| 数据库 | MySQL 8 (db4free / 本地 Docker) | 采用 UTF8MB4，开启外键约束和自增 ID。
| 工具 | Excelize 2.x、Golang-migrate、Swagger (Hertz 内置)、Taskfile、GitHub Actions | Excel 解析、数据库迁移、自动化测试/构建。
| 测试 | Go test + httptest + Playwright (前端 E2E) | 单元、集成、前端端到端。

## 3. 系统架构
```
┌────────────┐        HTTPS        ┌──────────────────┐       ┌───────────┐
│ React 前端 │ <────────────────> │ Hertz API Gateway │ <───> │ MySQL RDS │
└─────┬──────┘                     │ 业务 Service 层  │       └───────────┘
      │                            │ GORM DAL + Cache │
      ▼                            └────────┬────────┘
┌──────────────┐                            │
│ Figma 原型库 │                            ▼
│ Bilibili 视频│                    ┌──────────────┐
└──────────────┘                    │ Export/Import│
                                    └──────────────┘
```
关键流：
1. Excel -> 前端/后端解析 -> Import API -> GORM 批量写入 classes/students/enrollments。
2. 点名 -> RollCall API -> 权重随机算法 -> 记录 roll_call_records -> 返回前端实时展示。
3. 积分结算 -> SolveRollCall API -> 事务更新 enrollment & event log。
4. 导出 -> Export API -> 生成 CSV/Excel 流 -> 前端下载。
5. 排行 & 可视化 -> Leaderboard API -> 前端 ECharts 渲染。

## 4. 数据模型
| 表 | 关键字段 | 说明 |
| --- | --- | --- |
| classes | id(PK), class_name, created_at | 班级元数据 |
| students | id(PK), student_no, name, major | 基础学生信息 |
| enrollments | id(PK), class_id, student_id, total_points, call_count, transfer_rights, seq_index | 班级-学生关联与积分 |
| roll_call_records | id, class_id, enrollment_id, mode, event_type, called_at | 点名日志、用于统计卡片 |
| score_events | id, enrollment_id, delta, reason, event_type, metadata(json) | 记分流水，支撑性能分析 |
| attachments | id, class_id, file_name, storage_url, hash | Excel 导入/导出文件溯源（选做） |

额外索引：
- enrollments(class_id, total_points), enrollments(class_id, seq_index)
- roll_call_records(class_id, called_at desc)
- score_events(enrollment_id, created_at desc)

## 5. API 设计总览
| 功能 | Method & Path | 描述 |
| --- | --- | --- |
| Excel 导入 | `POST /v1/import/excel` (multipart) | 上传 Excel -> 解析 -> 创建班级 & 学生 |
| JSON 导入 | `POST /v1/import/class-data` | 兼容原 IDL，供测试/脚本使用 |
| 班级管理 | `GET/DELETE /v1/classes`, `GET /v1/classes/:id` | 查询/删除班级 |
| 花名册 | `GET /v1/classes/:id/roster` | 返回 roster 列表（含序号、积分、转移权） |
| 点名 | `POST /v1/roll-calls` | 支持 RANDOM/SEQUENTIAL/REVERSE/LOW_POINTS_FIRST + 事件 |
| 积分结算 | `POST /v1/roll-calls/solve` | 根据回答类型自动计算积分，记录 score_events |
| 排行榜 | `GET /v1/classes/:id/leaderboard?top=10` | Top N，总积分 & 被点次数统计 |
| 导出积分 | `GET /v1/classes/:id/export` | 生成 CSV/Excel 包含学号/姓名/专业/次数/积分 |
| 统计面板 | `GET /v1/dashboard/overview` | 今日点名次数、活跃班级、Top 5 加权等 |
| 健康检查 | `GET /ping` | 负载均衡、监控用途 |

## 6. 核心算法
1. **权重随机**：`weight = 1 / (1 + total_points * α + call_count * β)`，再结合随机事件倍率 (Double = ×2, Crazy Thursday = 指定因数)；使用累积分布 + `rand.Float64()` 实现。
2. **顺序模式**：`seq_index` 保存顺序；每次点名后自增 `next_seq_index`，反向模式倒序遍历。
3. **积分规则**：
   - 基础：被点到且签到 +1；
   - 能复述问题 +0.5，否则 -1；
   - 正确回答 +0.5~3（由 `custom_score` 控制，默认 1）；
   - 事件加成：Double -> 总增益 ×2；Crazy Thursday -> 总分若为 50 的因数则概率提升；
   - Transfer：被点两次 `call_count % 2 == 0` 则 `transfer_rights++`；`AnswerType=TRANSFER` 时可指定 `target_enrollment_id`，扣除 1 次权力。
4. **可视化统计**：`roll_call_records` 提供时序折线；`enrollments` 用于柱状图；`score_events` 支持性能分析。

## 7. 计划与里程碑
| 阶段 | 产出 | 截止 |
| --- | --- | --- |
| Day 1 | 架构设计、数据库迁移、基础 API 框架、Figma 原型草图 | D+1 |
| Day 2 | Excel 导入/导出、RollCall 算法、积分结算、Leaderboard | D+2 |
| Day 3 | React 前端（控制面板、排行榜、统计图）、接口联调 | D+3 |
| Day 4 | 单元/集成测试、性能分析、Docker 部署、README & 博客资料 | D+4 |
| Day 5 | Demo 视频脚本、录屏、博客终稿、提交任务 | D+5 |

## 8. 风险 & 应对
- **免费 MySQL 不稳定**：提供 `.env` + Docker Compose，支持本地/远程双通道；DAO 自动重试。
- **Excel 数据脏**：导入前进行 schema 校验、重复学号检测、事务回滚。
- **时间紧张**：Taskfile 自动化、模块化组件可并行开发。

## 9. 交付物清单
1. `docs/architecture.md`（本文）
2. `docs/ui-prototype/`：墨刀/Figma 链接与截图
3. `NekoCallerBackend/`：Go 服务代码、配置、测试、Taskfile
4. `web/`：React 前端源码、组件/storybook
5. `deploy/`：Dockerfile、docker-compose、SQL migration
6. `docs/blog/`：博客模板草稿（含 PSP、学习进度条等）
7. `video/`：演示脚本与 OBS 设置（文本）
