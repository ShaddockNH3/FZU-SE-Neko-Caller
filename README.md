# Neko Caller - 智能课堂点名系统

> 福州大学2025年下半学期软件工程结对编程

## 项目概述

Neko Caller 是一个功能完善的课堂点名管理系统,支持多种点名模式、积分管理、随机事件以及数据导入导出功能。

## 技术栈

### 后端
- **Go 1.24.4** - 主要编程语言
- **CloudWeGo Hertz** - 高性能HTTP框架
- **GORM** - ORM框架
- **MySQL 8** - 关系型数据库
- **Excelize** - Excel文件处理

### 前端
- **React 18** - UI框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Ant Design** - UI组件库

## 功能特性

### 基础功能
- ✅ **Excel导入导出** - 支持从Excel导入学生名单,导出积分详单
- ✅ **多种点名模式**
  - 随机点名(加权随机,积分越高概率越低)
  - 顺序点名
  - 反向顺序点名
  - 低分优先点名
- ✅ **积分管理系统**
  - 初始积分为0
  - 被点名到场 +1分
  - 请求帮助 +0.5分
  - 跳过回答 -1分
  - 自定义积分(0.5-3分)
- ✅ **积分排行榜** - 支持查询Top N排名
- ✅ **数据可视化** - 提供积分分布、点名频次统计

### 进阶功能
- ✅ **点名转移权**
  - 每被点名2次获得1次转移权
  - 可将点名机会转移给同班其他同学
  - 转移扣除0.5分
- ✅ **随机事件系统**
  - **双倍积分**: 权重×1.3, 得分×2
  - **疯狂星期四**: 积分为50的因数(1,2,5,10,25,50)的学生权重×1.25, 得分×1.5

## 快速开始

### 环境要求
- Go 1.24.4+
- MySQL 8.0+
- Node.js 18+ (可选,如果需要运行前端)

### 数据库配置

1. 创建MySQL数据库:
```sql
CREATE DATABASE neko_caller_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 配置环境变量(可选):
```bash
export DB_USER=root
export DB_PASSWORD=123456
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_NAME=neko_caller_db
```

默认配置: `root:123456@tcp(127.0.0.1:3306)/neko_caller_db`

### 后端启动

```bash
cd NekoCallerBackend
go mod tidy
go run main.go
```

服务默认运行在 `http://localhost:8888`

### 前端启动(可选)

```bash
cd web
npm install
npm run dev
```

前端默认运行在 `http://localhost:5173`

## API接口文档

### 班级管理

#### 获取班级列表
```
GET /v1/classes
```

#### 获取班级详情
```
GET /v1/classes/:class_id
```

#### 删除班级
```
DELETE /v1/classes/:class_id
```

#### 获取班级排行榜
```
GET /v1/classes/:class_id/leaderboard?top=10
```

#### 获取班级统计信息
```
GET /v1/classes/:class_id/stats
```

#### 导出班级积分详单(Excel)
```
GET /v1/classes/:class_id/export
```

### 学生管理

#### 获取学生列表
```
GET /v1/students
```

#### 获取学生详情
```
GET /v1/students/:student_id
```

#### 删除学生
```
DELETE /v1/students/:student_id
```

### 数据导入

#### JSON导入
```
POST /v1/import/class-data
Content-Type: application/json

{
  "class_name": "软工K班",
  "students": [
    {
      "student_id": "102101001",
      "name": "张三",
      "major": "软件工程"
    }
  ]
}
```

#### Excel导入
```
POST /v1/import/excel
Content-Type: multipart/form-data

file: <Excel文件>
class_name: 软工K班
```

Excel格式要求:
| 学号 | 姓名 | 专业(可选) |
|------|------|-----------|
| 102101001 | 张三 | 软件工程 |

### 花名册管理

#### 获取班级花名册
```
GET /v1/roster?class_id=xxx
```

#### 移除学生选课记录
```
DELETE /v1/enrollments/:enrollment_id
```

### 点名功能

#### 发起点名
```
POST /v1/roll-calls
Content-Type: application/json

{
  "class_id": "xxx",
  "mode": 0,  // 0-随机, 1-顺序, 2-反向, 3-低分优先
  "event_type": 0  // 0-无, 1-双倍积分, 2-疯狂星期四
}
```

#### 结算点名
```
POST /v1/roll-calls/solve
Content-Type: application/json

{
  "enrollment_id": "xxx",
  "answer_type": 0,  // 0-正常, 1-请求帮助, 2-跳过, 3-转移
  "custom_score": 2.0,  // 可选,仅answer_type=0时有效
  "event_type": 0,
  "target_enrollment_id": "yyy"  // 仅answer_type=3时必填
}
```

### 响应格式

成功响应:
```json
{
  "code": 100,
  "message": "操作成功"
}
```

错误响应:
```json
{
  "code": -1,
  "message": "错误信息"
}
```

## 数据库结构

### students - 学生表
| 字段 | 类型 | 说明 |
|------|------|------|
| student_id | VARCHAR(255) | 学号(主键) |
| name | VARCHAR(255) | 姓名 |
| major | VARCHAR(255) | 专业 |

### classes - 班级表
| 字段 | 类型 | 说明 |
|------|------|------|
| class_id | VARCHAR(255) | 班级ID(主键,UUID) |
| class_name | VARCHAR(255) | 班级名称 |

### enrollments - 选课表
| 字段 | 类型 | 说明 |
|------|------|------|
| enrollment_id | VARCHAR(255) | 选课ID(主键,UUID) |
| student_id | VARCHAR(255) | 学号(外键) |
| class_id | VARCHAR(255) | 班级ID(外键) |
| total_points | DOUBLE | 总积分 |
| call_count | BIGINT | 点名次数 |
| transfer_rights | BIGINT | 转移权次数 |

### roll_call_records - 点名记录表
| 字段 | 类型 | 说明 |
|------|------|------|
| record_id | VARCHAR(255) | 记录ID(主键,UUID) |
| class_id | VARCHAR(255) | 班级ID |
| enrollment_id | VARCHAR(255) | 选课ID |
| student_id | VARCHAR(255) | 学号 |
| mode | INT | 点名模式 |
| event_type | INT | 事件类型 |
| created_at | DATETIME | 创建时间 |

### score_events - 积分变动记录表
| 字段 | 类型 | 说明 |
|------|------|------|
| event_id | VARCHAR(255) | 事件ID(主键,UUID) |
| enrollment_id | VARCHAR(255) | 选课ID |
| student_id | VARCHAR(255) | 学号 |
| class_id | VARCHAR(255) | 班级ID |
| delta | DOUBLE | 积分变化量 |
| reason | VARCHAR(500) | 变动原因 |
| event_type | INT | 事件类型 |
| metadata | JSON | 元数据(如转移目标等) |
| created_at | DATETIME | 创建时间 |

## 核心算法

### 加权随机点名算法

基础权重计算:
```
weight = 1 / (1 + points * 0.4 + callCount * 0.2)
```

事件权重调整:
- 双倍积分: `weight *= 1.3`
- 疯狂星期四: 积分为50的因数时 `weight *= 1.25`

使用累积分布函数进行随机抽样,确保权重越高被选中概率越大。

### 低分优先算法

1. 按积分升序排序(积分相同按点名次数排序)
2. 选取前1/3的学生
3. 从这部分学生中随机选择

### 转移权机制

- **获得条件**: 每被点名2次(call_count变为偶数)自动获得1次转移权
- **使用方式**: 结算时选择"转移"答案类型,指定目标学生
- **验证规则**: 
  - 必须有转移权(transfer_rights > 0)
  - 目标必须是同班学生
  - 不能转移给自己
- **效果**: 
  - 源学生: -0.5分, 转移权-1
  - 目标学生: call_count+1, 可能获得转移权

## 项目结构

```
FZU-SE-Neko-Caller/
├── NekoCallerBackend/          # 后端代码
│   ├── main.go                 # 程序入口
│   ├── biz/                    # 业务逻辑
│   │   ├── dal/               # 数据访问层
│   │   │   ├── model/        # 数据模型
│   │   │   ├── mysql/        # MySQL初始化
│   │   │   └── query/        # GORM生成的查询代码
│   │   ├── handler/          # HTTP处理器
│   │   ├── model/            # API模型(Thrift生成)
│   │   ├── router/           # 路由注册
│   │   └── service/          # 业务服务层
│   ├── idl/                  # Thrift IDL定义
│   └── pkg/                  # 公共包
│       ├── constants/        # 常量定义
│       ├── errno/            # 错误码定义
│       └── utils/            # 工具函数
└── web/                      # 前端代码(React)
    └── src/
        ├── pages/           # 页面组件
        └── styles/          # 样式文件
```

## 测试

```bash
cd NekoCallerBackend
go test ./...
```

## 开发团队

- ShaddockNH3

## 许可证

本项目仅用于学习和教学目的。

## 更新日志

### v1.0.0 (2025-11-16)
- ✅ 完成Excel导入导出功能
- ✅ 实现多种点名模式
- ✅ 积分管理与排行榜
- ✅ 点名转移权机制
- ✅ 随机事件系统(双倍积分、疯狂星期四)
- ✅ 数据可视化统计API
- ✅ 完整的REST API接口

## 常见问题

### Q: 如何修改默认数据库配置?
A: 可以通过环境变量覆盖默认配置,参考"数据库配置"章节。

### Q: Excel导入失败怎么办?
A: 请确保Excel文件格式正确,至少包含"学号"和"姓名"两列,专业列可选。

### Q: 疯狂星期四什么时候触发?
A: 当学生的积分恰好是50的因数(1,2,5,10,25,50)时,在随机点名中权重增加25%,结算时积分变化乘以1.5倍。

### Q: 转移权如何使用?
A: 在结算点名时,选择answer_type=3(转移),并指定target_enrollment_id为同班其他学生的选课ID即可。

### Q: 前端是必须的吗?
A: 不是,后端提供完整的REST API,可以直接使用Postman等工具测试和使用所有功能。

## 联系方式

如有问题或建议,请通过GitHub Issues反馈。
