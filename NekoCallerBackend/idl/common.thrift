namespace go common

// 返回信息
struct BaseResponse{
    1: i32 code,       // 状态码，100表示成功，-1表示失败
    2: string message, // 返回信息
}

// 学生信息
struct Student{
    1: string student_id    // 学号
    2: string name      // 姓名
    3: optional string major    // 专业
}

// 班级信息
struct Class{
    1: string class_id    // 班级ID
    2: string class_name  // 班级名称
    3: optional list<string> student_ids // 班级内学生ID列表
}

// 关联信息
struct Enrollment {
    1: string enrollment_id,    // 选课记录的唯一ID
    2: string student_id,       // 关联到学生
    3: string class_id,         // 关联到班级
    4: double total_points,     // 该学生在该班级的总积分
    5: i64 call_count,         // 该学生在该班级的被点名次数
    6: i64 transfer_rights,    // 该学生在该班级的点名转移权
    7: i64 skip_rights,        // 该学生在该班级的跳过权
}

// 获取班级花名册时，返回的数据项
struct RosterItem {
    1: Student student_info,
    2: Enrollment enrollment_info,
}

// 点名模式
enum RollCallMode{
    RANDOM = 0,     // 随机
    SEQUENTIAL = 1,     // 顺序
    REVERSE_SEQUENTIAL = 2, // 反向顺序
    LOW_POINTS_FIRST = 3, // 积分低优先，按照正态分布处理
}

// 随机事件类型
enum RandomEventType{
    NONE = 0,
    Double_Point = 1,   // 双倍积分
    CRAZY_THURSDAY = 2, // 疯狂星期四
    BLESSING_1024 = 3,      // 1024 程序员福报
    SOLITUDE_PRIMES = 4, // 质数的孤独
    LUCKY_7 = 5,        // 幸运 7 大奖
    // 后续扩展
}

enum AnswerType{
    NORMAL = 0,     // 正常回答
    HELP = 1,       // 请求帮助
    SKIP = 2,       // 跳过回答
    TRANSFER = 3    // 转移回答
}
