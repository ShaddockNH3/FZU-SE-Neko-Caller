# Neko Caller API 接口文档

## 基础信息

- **Base URL**: `http://localhost:8888`
- **Content-Type**: `application/json` (除特殊说明外)
- **字符编码**: `UTF-8`

## 响应格式说明

### 成功响应
```json
{
  "code": 100,
  "message": "操作成功"
}
```

### 错误响应
```json
{
  "code": -1,
  "message": "错误信息描述"
}
```

## 状态码说明

| 状态码 | 说明 |
|--------|------|
| 100 | 操作成功 |
| -1 | 操作失败 |

---

## 一、数据导入接口

### 1.1 JSON格式导入班级数据

**接口地址**: `POST /v1/import/class-data`

**请求头**:
```
Content-Type: application/json
```

**请求体**:
```json
{
  "class_name": "软工K班",
  "students": [
    {
      "student_id": "102101001",
      "name": "张三",
      "major": "软件工程"
    },
    {
      "student_id": "102101002",
      "name": "李四",
      "major": "计算机科学"
    }
  ]
}
```

**参数说明**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| class_name | string | 是 | 班级名称 |
| students | array | 是 | 学生列表 |
| students[].student_id | string | 是 | 学号 |
| students[].name | string | 是 | 姓名 |
| students[].major | string | 否 | 专业 |

**响应示例**:
```json
{
  "code": 100,
  "message": "导入成功"
}
```

**说明**:
- 如果学生已存在，将使用FirstOrCreate逻辑（更新姓名，专业仅在非空时更新）
- 重复的选课记录会自动去重
- 所有操作在事务中执行，失败自动回滚

---

### 1.2 Excel文件导入班级数据

**接口地址**: `POST /v1/import/excel`

**请求头**:
```
Content-Type: multipart/form-data
```

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| file | file | 是 | Excel文件(.xlsx) |
| class_name | string | 是 | 班级名称 |

**Excel格式要求**:

| 列A(学号) | 列B(姓名) | 列C(专业,可选) |
|-----------|-----------|----------------|
| 102101001 | 张三 | 软件工程 |
| 102101002 | 李四 | 计算机科学 |

**说明**:
- 第一行为标题行，会被跳过
- 至少需要学号和姓名两列
- 专业列可选，如果为空则不设置专业

**响应示例**:
```json
{
  "code": 100,
  "message": "成功导入班级 软工K班，共 30 名学生"
}
```

**Postman示例**:
```
POST http://localhost:8888/v1/import/excel
Body:
  - form-data
    - file: [选择Excel文件]
    - class_name: 软工K班
```

**cURL示例**:
```bash
curl -X POST http://localhost:8888/v1/import/excel \
  -F "file=@/path/to/students.xlsx" \
  -F "class_name=软工K班"
```

---

## 二、班级管理接口

### 2.1 获取班级列表

**接口地址**: `GET /v1/classes`

**请求参数**: 无

**响应示例**:
```json
[
  {
    "class_id": "550e8400-e29b-41d4-a716-446655440000",
    "class_name": "软工K班",
    "student_ids": [
      "102101001",
      "102101002",
      "102101003"
    ]
  },
  {
    "class_id": "660e8400-e29b-41d4-a716-446655440001",
    "class_name": "计科A班",
    "student_ids": [
      "102201001",
      "102201002"
    ]
  }
]
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| class_id | string | 班级ID(UUID) |
| class_name | string | 班级名称 |
| student_ids | array | 班级内所有学生的学号列表 |

---

### 2.2 获取班级详情

**接口地址**: `GET /v1/classes/:class_id`

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| class_id | string | 是 | 班级ID |

**请求示例**:
```
GET http://localhost:8888/v1/classes/550e8400-e29b-41d4-a716-446655440000
```

**响应示例**:
```json
{
  "class_id": "550e8400-e29b-41d4-a716-446655440000",
  "class_name": "软工K班",
  "student_ids": [
    "102101001",
    "102101002",
    "102101003"
  ]
}
```

---

### 2.3 删除班级

**接口地址**: `DELETE /v1/classes/:class_id`

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| class_id | string | 是 | 班级ID |

**请求示例**:
```
DELETE http://localhost:8888/v1/classes/550e8400-e29b-41d4-a716-446655440000
```

**响应示例**:
```json
{
  "code": 100,
  "message": "删除班级成功"
}
```

**说明**:
- 删除班级会同时删除所有相关的选课记录(enrollments)
- 操作在事务中执行

---

### 2.4 获取班级排行榜

**接口地址**: `GET /v1/classes/:class_id/leaderboard`

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| class_id | string | 是 | 班级ID |

**查询参数**:

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| top | int | 否 | 10 | 返回前N名 |

**请求示例**:
```
GET http://localhost:8888/v1/classes/550e8400-e29b-41d4-a716-446655440000/leaderboard?top=10
```

**响应示例**:
```json
[
  {
    "rank": 1,
    "student_id": "102101001",
    "name": "张三",
    "major": "软件工程",
    "total_points": 25.5,
    "call_count": 12
  },
  {
    "rank": 2,
    "student_id": "102101002",
    "name": "李四",
    "major": "计算机科学",
    "total_points": 23.0,
    "call_count": 10
  },
  {
    "rank": 3,
    "student_id": "102101003",
    "name": "王五",
    "total_points": 20.5,
    "call_count": 11
  }
]
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| rank | int | 排名(从1开始) |
| student_id | string | 学号 |
| name | string | 姓名 |
| major | string | 专业(可能为null) |
| total_points | float | 总积分 |
| call_count | int | 被点名次数 |

**排序规则**:
1. 按total_points降序
2. 积分相同时按call_count降序

---

### 2.5 获取班级统计信息

**接口地址**: `GET /v1/classes/:class_id/stats`

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| class_id | string | 是 | 班级ID |

**请求示例**:
```
GET http://localhost:8888/v1/classes/550e8400-e29b-41d4-a716-446655440000/stats
```

**响应示例**:
```json
{
  "total_students": 30,
  "total_calls": 150,
  "average_points": 15.5,
  "points_distribution": {
    "0-10": 5,
    "10-20": 12,
    "20-30": 10,
    "30-40": 3
  },
  "call_frequency": {
    "3": 2,
    "4": 5,
    "5": 8,
    "6": 10,
    "7": 5
  }
}
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| total_students | int | 班级总人数 |
| total_calls | int | 累计点名总次数 |
| average_points | float | 平均积分(保留2位小数) |
| points_distribution | object | 积分区间分布(key为区间,value为人数) |
| call_frequency | object | 点名次数分布(key为次数,value为人数) |

**说明**:
- 积分区间按每10分划分: 0-10, 10-20, 20-30...
- 可用于生成直方图和分布图

---

### 2.6 导出班级积分详单(Excel)

**接口地址**: `GET /v1/classes/:class_id/export`

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| class_id | string | 是 | 班级ID |

**请求示例**:
```
GET http://localhost:8888/v1/classes/550e8400-e29b-41d4-a716-446655440000/export
```

**响应**:
- Content-Type: `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- Content-Disposition: `attachment; filename*=UTF-8''软工K班-积分详单.xlsx`
- Body: Excel文件二进制流

**Excel内容格式**:

| 学号 | 姓名 | 专业 | 随机点名次数 | 总积分 |
|------|------|------|--------------|--------|
| 102101001 | 张三 | 软件工程 | 12 | 25.5 |
| 102101002 | 李四 | 计算机科学 | 10 | 23.0 |

**浏览器使用**:
- 直接在浏览器访问此URL即可下载Excel文件

**cURL示例**:
```bash
curl http://localhost:8888/v1/classes/550e8400-e29b-41d4-a716-446655440000/export \
  -o 积分详单.xlsx
```

---

## 三、学生管理接口

### 3.1 获取学生列表

**接口地址**: `GET /v1/students`

**请求参数**: 无

**响应示例**:
```json
[
  {
    "student_id": "102101001",
    "name": "张三",
    "major": "软件工程"
  },
  {
    "student_id": "102101002",
    "name": "李四",
    "major": "计算机科学"
  },
  {
    "student_id": "102101003",
    "name": "王五"
  }
]
```

**响应字段说明**:

| 字段名 | 类型 | 说明 |
|--------|------|------|
| student_id | string | 学号 |
| name | string | 姓名 |
| major | string | 专业(可能为null) |

---

### 3.2 获取学生详情

**接口地址**: `GET /v1/students/:student_id`

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| student_id | string | 是 | 学号 |

**请求示例**:
```
GET http://localhost:8888/v1/students/102101001
```

**响应示例**:
```json
{
  "student_id": "102101001",
  "name": "张三",
  "major": "软件工程"
}
```

---

### 3.3 删除学生

**接口地址**: `DELETE /v1/students/:student_id`

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| student_id | string | 是 | 学号 |

**请求示例**:
```
DELETE http://localhost:8888/v1/students/102101001
```

**响应示例**:
```json
{
  "code": 100,
  "message": "删除学生成功"
}
```

**说明**:
- 删除学生会同时删除所有相关的选课记录(enrollments)
- 操作在事务中执行

---

## 四、花名册管理接口

### 4.1 获取班级花名册

**接口地址**: `GET /v1/roster`

**查询参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| class_id | string | 是 | 班级ID |

**请求示例**:
```
GET http://localhost:8888/v1/roster?class_id=550e8400-e29b-41d4-a716-446655440000
```

**响应示例**:
```json
[
  {
    "student_info": {
      "student_id": "102101001",
      "name": "张三",
      "major": "软件工程"
    },
    "enrollment_info": {
      "enrollment_id": "770e8400-e29b-41d4-a716-446655440000",
      "student_id": "102101001",
      "class_id": "550e8400-e29b-41d4-a716-446655440000",
      "total_points": 25.5,
      "call_count": 12,
      "transfer_rights": 2
    }
  },
  {
    "student_info": {
      "student_id": "102101002",
      "name": "李四",
      "major": "计算机科学"
    },
    "enrollment_info": {
      "enrollment_id": "880e8400-e29b-41d4-a716-446655440000",
      "student_id": "102101002",
      "class_id": "550e8400-e29b-41d4-a716-446655440000",
      "total_points": 23.0,
      "call_count": 10,
      "transfer_rights": 1
    }
  }
]
```

**响应字段说明**:

| 字段路径 | 类型 | 说明 |
|----------|------|------|
| student_info | object | 学生基本信息 |
| student_info.student_id | string | 学号 |
| student_info.name | string | 姓名 |
| student_info.major | string | 专业(可能为null) |
| enrollment_info | object | 选课信息 |
| enrollment_info.enrollment_id | string | 选课记录ID |
| enrollment_info.student_id | string | 学号 |
| enrollment_info.class_id | string | 班级ID |
| enrollment_info.total_points | float | 该学生在该班级的总积分 |
| enrollment_info.call_count | int | 该学生在该班级的被点名次数 |
| enrollment_info.transfer_rights | int | 该学生在该班级的转移权次数 |

**排序规则**:
1. 按call_count升序(点名次数少的在前)
2. call_count相同时按student_id升序

---

### 4.2 移除学生的选课记录

**接口地址**: `DELETE /v1/enrollments/:enrollment_id`

**路径参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| enrollment_id | string | 是 | 选课记录ID |

**请求示例**:
```
DELETE http://localhost:8888/v1/enrollments/770e8400-e29b-41d4-a716-446655440000
```

**响应示例**:
```json
{
  "code": 100,
  "message": "移除学生成功"
}
```

**说明**:
- 仅删除该学生在该班级的选课记录
- 不会删除学生本身

---

## 五、点名功能接口

### 5.1 发起点名

**接口地址**: `POST /v1/roll-calls`

**请求头**:
```
Content-Type: application/json
```

**请求体**:
```json
{
  "class_id": "550e8400-e29b-41d4-a716-446655440000",
  "mode": 0,
  "event_type": 0
}
```

**参数说明**:

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| class_id | string | 是 | - | 班级ID |
| mode | int | 是 | - | 点名模式 |
| event_type | int | 否 | 0 | 随机事件类型 |

**点名模式(mode)**:

| 值 | 模式 | 说明 |
|----|------|------|
| 0 | RANDOM | 加权随机(积分越高概率越低) |
| 1 | SEQUENTIAL | 顺序点名(第一个) |
| 2 | REVERSE_SEQUENTIAL | 反向顺序(最后一个) |
| 3 | LOW_POINTS_FIRST | 低分优先(从最低1/3中随机) |

**随机事件类型(event_type)**:

| 值 | 事件 | 说明 |
|----|------|------|
| 0 | NONE | 无事件 |
| 1 | Double_Point | 双倍积分(选中权重×1.3,得分×2) |
| 2 | CRAZY_THURSDAY | 疯狂星期四(积分为50因数的权重×1.25,得分×1.5) |

**响应示例**:
```json
{
  "base_response": {
    "code": 100,
    "message": "点名成功"
  },
  "roster_item": {
    "student_info": {
      "student_id": "102101001",
      "name": "张三",
      "major": "软件工程"
    },
    "enrollment_info": {
      "enrollment_id": "770e8400-e29b-41d4-a716-446655440000",
      "student_id": "102101001",
      "class_id": "550e8400-e29b-41d4-a716-446655440000",
      "total_points": 25.5,
      "call_count": 13,
      "transfer_rights": 2
    }
  }
}
```

**说明**:
- 点名后会自动增加该学生的call_count
- 每2次点名(call_count变为偶数)会自动获得1次转移权
- 响应中的roster_item包含最新的点名次数

**加权随机算法说明**:
```
基础权重 = 1 / (1 + total_points * 0.4 + call_count * 0.2)
```
- 积分越高,权重越低,被点到的概率越小
- 点名次数越多,权重也会降低

**疯狂星期四触发条件**:
- 学生的total_points恰好是50的因数: 1, 2, 5, 10, 25, 50
- 例如: 积分为10.0分的学生在疯狂星期四事件下权重×1.25

---

### 5.2 结算点名

**接口地址**: `POST /v1/roll-calls/solve`

**请求头**:
```
Content-Type: application/json
```

**请求体**:
```json
{
  "enrollment_id": "770e8400-e29b-41d4-a716-446655440000",
  "answer_type": 0,
  "custom_score": 2.0,
  "event_type": 0,
  "target_enrollment_id": null
}
```

**参数说明**:

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| enrollment_id | string | 是 | - | 被点名学生的选课记录ID |
| answer_type | int | 是 | - | 回答类型 |
| custom_score | float | 否 | - | 自定义分数(仅answer_type=0时有效) |
| event_type | int | 否 | 0 | 随机事件类型(需与点名时一致) |
| target_enrollment_id | string | 否 | - | 转移目标(仅answer_type=3时必填) |

**回答类型(answer_type)**:

| 值 | 类型 | 基础分数 | 说明 |
|----|------|----------|------|
| 0 | NORMAL | 1.0或custom_score | 正常回答,可指定0.5-3分 |
| 1 | HELP | 0.5 | 请求帮助 |
| 2 | SKIP | -1.0 | 跳过回答 |
| 3 | TRANSFER | -0.5 | 转移给其他同学 |

**随机事件积分加成**:
- Double_Point: 积分×2
- CRAZY_THURSDAY: 积分×1.5

**响应示例**:
```json
{
  "code": 100,
  "message": "结算成功"
}
```

**转移示例**:
```json
{
  "enrollment_id": "770e8400-e29b-41d4-a716-446655440000",
  "answer_type": 3,
  "event_type": 0,
  "target_enrollment_id": "880e8400-e29b-41d4-a716-446655440000"
}
```

**转移规则**:
1. 源学生必须有转移权(transfer_rights > 0)
2. 目标必须是同班学生
3. 不能转移给自己
4. 转移后:
   - 源学生: -0.5分, transfer_rights-1
   - 目标学生: call_count+1, 可能获得转移权

**计分示例**:

| 场景 | 计算 | 最终得分 |
|------|------|----------|
| 正常回答(无事件) | 1.0 | 1.0 |
| 正常回答(双倍积分) | 1.0 × 2 | 2.0 |
| 正常回答(疯狂星期四) | 1.0 × 1.5 | 1.5 |
| 自定义2分(双倍积分) | 2.0 × 2 | 4.0 |
| 请求帮助(双倍积分) | 0.5 × 2 | 1.0 |
| 跳过(双倍积分) | -1.0 × 2 | -2.0 |
| 转移(无事件) | -0.5 | -0.5 |

**说明**:
- 所有积分变动会记录在score_events表中
- 转移操作会在metadata字段记录target_enrollment_id和target_student_id

---

## 六、数据模型说明

### 6.1 Student (学生)
```json
{
  "student_id": "string",  // 学号(主键)
  "name": "string",        // 姓名
  "major": "string|null"   // 专业(可选)
}
```

### 6.2 Class (班级)
```json
{
  "class_id": "string",        // 班级ID(UUID)
  "class_name": "string",      // 班级名称
  "student_ids": ["string"]    // 学生ID列表
}
```

### 6.3 Enrollment (选课记录)
```json
{
  "enrollment_id": "string",   // 选课ID(UUID)
  "student_id": "string",      // 学号
  "class_id": "string",        // 班级ID
  "total_points": "float",     // 总积分
  "call_count": "int",         // 点名次数
  "transfer_rights": "int"     // 转移权次数
}
```

### 6.4 RosterItem (花名册项)
```json
{
  "student_info": {
    "student_id": "string",
    "name": "string",
    "major": "string|null"
  },
  "enrollment_info": {
    "enrollment_id": "string",
    "student_id": "string",
    "class_id": "string",
    "total_points": "float",
    "call_count": "int",
    "transfer_rights": "int"
  }
}
```

---

## 七、完整使用流程示例

### 步骤1: 导入班级数据
```bash
curl -X POST http://localhost:8888/v1/import/excel \
  -F "file=@students.xlsx" \
  -F "class_name=软工K班"
```

### 步骤2: 查看班级列表
```bash
curl http://localhost:8888/v1/classes
```
获取class_id: `550e8400-e29b-41d4-a716-446655440000`

### 步骤3: 查看花名册
```bash
curl "http://localhost:8888/v1/roster?class_id=550e8400-e29b-41d4-a716-446655440000"
```

### 步骤4: 发起随机点名
```bash
curl -X POST http://localhost:8888/v1/roll-calls \
  -H "Content-Type: application/json" \
  -d '{
    "class_id": "550e8400-e29b-41d4-a716-446655440000",
    "mode": 0,
    "event_type": 0
  }'
```
获取被点名学生的enrollment_id: `770e8400-e29b-41d4-a716-446655440000`

### 步骤5: 结算点名(正常回答)
```bash
curl -X POST http://localhost:8888/v1/roll-calls/solve \
  -H "Content-Type: application/json" \
  -d '{
    "enrollment_id": "770e8400-e29b-41d4-a716-446655440000",
    "answer_type": 0,
    "custom_score": 2.0,
    "event_type": 0
  }'
```

### 步骤6: 查看排行榜
```bash
curl "http://localhost:8888/v1/classes/550e8400-e29b-41d4-a716-446655440000/leaderboard?top=10"
```

### 步骤7: 导出积分详单
```bash
curl "http://localhost:8888/v1/classes/550e8400-e29b-41d4-a716-446655440000/export" \
  -o 积分详单.xlsx
```

---

## 八、错误处理

### 常见错误示例

#### 1. 班级不存在
```json
{
  "code": -1,
  "message": "班级不存在: record not found"
}
```

#### 2. 学生不存在
```json
{
  "code": -1,
  "message": "学生不存在"
}
```

#### 3. 转移权不足
```json
{
  "code": -1,
  "message": "转移权不足"
}
```

#### 4. 无效的转移目标
```json
{
  "code": -1,
  "message": "目标学生不存在或不在同一班级"
}
```

#### 5. Excel格式错误
```json
{
  "code": -1,
  "message": "解析Excel文件失败: invalid format"
}
```

---

## 九、Postman测试集合

### 环境变量设置
```json
{
  "base_url": "http://localhost:8888",
  "class_id": "",
  "student_id": "",
  "enrollment_id": ""
}
```

### 测试顺序建议
1. 导入班级数据 (POST /v1/import/excel)
2. 获取班级列表 (GET /v1/classes) → 保存class_id
3. 查看花名册 (GET /v1/roster)
4. 发起点名 (POST /v1/roll-calls) → 保存enrollment_id
5. 结算点名 (POST /v1/roll-calls/solve)
6. 查看排行榜 (GET /v1/classes/:class_id/leaderboard)
7. 查看统计 (GET /v1/classes/:class_id/stats)
8. 导出Excel (GET /v1/classes/:class_id/export)

---

## 十、注意事项

1. **UUID格式**: 所有ID字段使用UUID v4格式
2. **事务处理**: 导入、删除等操作都在事务中执行，失败会自动回滚
3. **外键约束**: 数据库禁用了外键约束以避免迁移问题，但逻辑层保证数据一致性
4. **积分精度**: 积分计算结果保留2位小数
5. **转移权机制**: 每2次点名自动获得1次，使用时自动扣除
6. **事件类型**: 点名和结算时的event_type应保持一致
7. **Excel编码**: 导入导出都使用UTF-8编码

---

## 附录: 数据库表结构

### students
```sql
CREATE TABLE `students` (
  `student_id` varchar(255) PRIMARY KEY,
  `name` varchar(255) NOT NULL,
  `major` varchar(255)
);
```

### classes
```sql
CREATE TABLE `classes` (
  `class_id` varchar(255) PRIMARY KEY,
  `class_name` varchar(255) NOT NULL
);
```

### enrollments
```sql
CREATE TABLE `enrollments` (
  `enrollment_id` varchar(255) PRIMARY KEY,
  `student_id` varchar(255) NOT NULL,
  `class_id` varchar(255) NOT NULL,
  `total_points` double DEFAULT 0,
  `call_count` bigint DEFAULT 0,
  `transfer_rights` bigint DEFAULT 0
);
```

### roll_call_records
```sql
CREATE TABLE `roll_call_records` (
  `record_id` varchar(255) PRIMARY KEY,
  `class_id` varchar(255),
  `enrollment_id` varchar(255),
  `student_id` varchar(255),
  `mode` int NOT NULL,
  `event_type` int NOT NULL,
  `created_at` datetime
);
```

### score_events
```sql
CREATE TABLE `score_events` (
  `event_id` varchar(255) PRIMARY KEY,
  `enrollment_id` varchar(255),
  `student_id` varchar(255),
  `class_id` varchar(255),
  `delta` double NOT NULL,
  `reason` varchar(255) NOT NULL,
  `event_type` int NOT NULL,
  `metadata` json,
  `created_at` datetime
);
```

---

**文档版本**: v1.0.0  
**最后更新**: 2025-11-16  
**联系方式**: GitHub Issues
