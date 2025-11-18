include "common.thrift"

namespace go api

// 导入班级请求
struct ImportDataRequest {
    1: string class_name,
    2: list<common.Student> students,
}

// 点名请求
struct RollCallRequest{
    1: string class_id,
    2: common.RollCallMode mode,
    3: optional common.RandomEventType event_type = common.RandomEventType.NONE,
}

// 点名响应
struct RollCallResponse{
    1: common.BaseResponse base_response,
    2: optional common.RosterItem roster_item,
    3: optional common.RandomEventType actual_event_type = common.RandomEventType.NONE, // 实际生效的事件类型
}

// 分数变动请求
struct SolveRollCallRequest{
    1: string enrollment_id,
    2: common.AnswerType answer_type,
    3: optional double custom_score,
    4: optional common.RandomEventType event_type = common.RandomEventType.NONE,
    5: optional string target_enrollment_id,
}

// 排行榜项
struct LeaderboardItem {
    1: i32 rank,
    2: string student_id,
    3: string name,
    4: optional string major,
    5: double total_points,
    6: i64 call_count,
}

// 班级统计信息
struct ClassStats {
    1: i32 total_students,
    2: i64 total_calls,
    3: double average_points,
    4: map<string, i32> points_distribution,  // 积分区间分布
    5: map<string, i64> call_frequency,       // 点名次数分布
}

service ApiService {
    // ==================== 导入接口 ====================
    // 导入班级及学生数据（JSON格式）
    common.BaseResponse ImportClassData(1: ImportDataRequest req) (api.post="/v1/import/class-data")
    
    // 导入班级及学生数据（Excel文件）
    // 注意：此接口使用multipart/form-data，无法在thrift中完整定义，需手动实现
    // POST /v1/import/excel
    
    // ==================== 班级相关接口 ====================
    // 获取班级信息
    common.Class GetClass(1: string class_id (api.path="class_id")) (api.get="/v1/classes/:class_id")
    
    // 获取所有班级列表
    list<common.Class> ListClasses() (api.get="/v1/classes")
    
    // 删除班级
    common.BaseResponse DeleteClass(1: string class_id (api.path="class_id")) (api.delete="/v1/classes/:class_id")
    
    // 导出班级花名册为Excel
    // 注意：此接口返回文件流，无法在thrift中完整定义，需手动实现
    // GET /v1/classes/:class_id/export
    
    // 获取班级积分排行榜
    // 注意：top参数通过查询参数传递，hz工具不支持多参数函数
    list<LeaderboardItem> GetLeaderboard(1: string class_id (api.path="class_id")) (api.get="/v1/classes/:class_id/leaderboard")
    
    // 获取班级统计信息
    ClassStats GetClassStats(1: string class_id (api.path="class_id")) (api.get="/v1/classes/:class_id/stats")

    // ==================== 学生相关接口 ====================
    // 获取学生信息
    common.Student GetStudent(1: string student_id (api.path="student_id")) (api.get="/v1/students/:student_id")
    
    // 获取所有学生列表
    list<common.Student> ListAllStudents() (api.get="/v1/students")
    
    // 删除学生
    common.BaseResponse DeleteStudent(1: string student_id (api.path="student_id")) (api.delete="/v1/students/:student_id")
    
    // ==================== 班级花名册相关接口 ====================
    // 获取班级花名册
    list<common.RosterItem> GetClassRoster(1: string class_id (api.query="class_id")) (api.get="/v1/roster")
    
    // 从班级中移除学生
    common.BaseResponse RemoveStudentFromClass(1: string enrollment_id (api.path="enrollment_id")) (api.delete="/v1/enrollments/:enrollment_id")

    // ==================== 点名相关接口 ====================
    // 执行点名
    RollCallResponse RollCall(1: RollCallRequest req) (api.post="/v1/roll-calls")
    
    // 点名结算
    common.BaseResponse SolveRollCall(1: SolveRollCallRequest req) (api.post="/v1/roll-calls/solve")
    
    // 重置点名状态（用于顺序/逆序点名）
    common.BaseResponse ResetRollCall(1: string class_id (api.body="class_id")) (api.post="/v1/roll-calls/reset")
}
