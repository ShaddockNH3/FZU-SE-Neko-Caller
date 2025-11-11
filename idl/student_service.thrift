namespace go student

struct Student{
    1: string student_id
    2: string name
    3: optional string major
    4: optional double totalPoints = 0.0
    5: optional i64 callCount = 0
    6: optional i64 transferRights = 0  // 点名转移权
}

enum RollCallMode{
    RANDOM = 0,
    SEQUENTIAL = 1,
}

enum RandomEventType{
    NONE = 0,
    Double_Point = 1,
    CRAZY_THURSDAY = 2,
    // 后续扩展
}
