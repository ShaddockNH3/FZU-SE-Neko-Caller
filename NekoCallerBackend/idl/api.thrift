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
}

// 分数变动请求
struct SolveRollCallRequest{
    1: string enrollment_id,
    2: common.AnswerType answer_type,
    3: optional double custom_score,
    4: optional common.RandomEventType event_type = common.RandomEventType.NONE,
    5: optional string target_enrollment_id,
}

service ApiService {
    // 导入班级及学生数据
    common.BaseResponse ImportClassData(1: ImportDataRequest req) (api.post="/v1/import/class-data")

    // 班级相关接口
    common.Class GetClass(1: string class_id (api.path="class_id")) (api.get="/v1/classes/:class_id")
    list<common.Class> ListClasses() (api.get="/v1/classes")
    common.BaseResponse DeleteClass(1: string class_id (api.path="class_id")) (api.delete="/v1/classes/:class_id")

    // 学生相关接口
    common.Student GetStudent(1: string student_id (api.path="student_id")) (api.get="/v1/students/:student_id")
    list<common.Student> ListAllStudents() (api.get="/v1/students")
    common.BaseResponse DeleteStudent(1: string student_id (api.path="student_id")) (api.delete="/v1/students/:student_id")
    
    // 班级花名册相关接口
    list<common.RosterItem> GetClassRoster(1: string class_id (api.query="class_id")) (api.get="/v1/roster")
    common.BaseResponse RemoveStudentFromClass(1: string enrollment_id (api.path="enrollment_id")) (api.delete="/v1/enrollments/:enrollment_id")

    // 点名相关接口
    RollCallResponse RollCall(1: RollCallRequest req) (api.post="/v1/roll-calls")
    common.BaseResponse SolveRollCall(1: SolveRollCallRequest req) (api.post="/v1/roll-calls/solve")
}
