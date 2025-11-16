package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"FZUSENekoCaller/biz/dal/model"
	"FZUSENekoCaller/biz/dal/mysql"
	"FZUSENekoCaller/biz/dal/query"
	"FZUSENekoCaller/biz/model/api"
	"FZUSENekoCaller/biz/model/common"
	"FZUSENekoCaller/pkg/errno"
	"FZUSENekoCaller/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type APIService struct {
	ctx context.Context
	db  *gorm.DB
}

func NewAPIService(ctx context.Context) *APIService {
	return &APIService{
		ctx: ctx,
		db:  mysql.GetDB(),
	}
}

func (s *APIService) GetClass(classID string) (*common.Class, error) {
	q := query.Use(s.db)
	cls, err := q.Class.WithContext(s.ctx).
		Where(q.Class.ClassID.Eq(classID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.StudentNotFoundErr
		}
		return nil, errno.ServiceErr
	}

	var studentIDs []string
	if err := q.Enrollment.WithContext(s.ctx).
		Where(q.Enrollment.ClassID.Eq(classID)).
		Pluck(q.Enrollment.StudentID, &studentIDs); err != nil {
		return nil, errno.QueryStudentIDsErr
	}

	return &common.Class{
		ClassID:    cls.ClassID,
		ClassName:  cls.ClassName,
		StudentIds: studentIDs,
	}, nil
}

func (s *APIService) ListClasses() ([]*common.Class, error) {
	q := query.Use(s.db)
	classes, err := q.Class.WithContext(s.ctx).Find()
	if err != nil {
		return nil, errno.ServiceErr
	}
	if len(classes) == 0 {
		return []*common.Class{}, nil
	}

	classIDs := make([]string, 0, len(classes))
	for _, cls := range classes {
		classIDs = append(classIDs, cls.ClassID)
	}

	enrollments, err := q.Enrollment.WithContext(s.ctx).
		Where(q.Enrollment.ClassID.In(classIDs...)).
		Find()
	if err != nil {
		return nil, errno.QueryStudentIDsErr
	}

	studentsByClass := make(map[string][]string, len(classes))
	for _, en := range enrollments {
		studentsByClass[en.ClassID] = append(studentsByClass[en.ClassID], en.StudentID)
	}

	result := make([]*common.Class, 0, len(classes))
	for _, cls := range classes {
		result = append(result, &common.Class{
			ClassID:    cls.ClassID,
			ClassName:  cls.ClassName,
			StudentIds: studentsByClass[cls.ClassID],
		})
	}

	return result, nil
}

func (s *APIService) DeleteClass(classID string) error {
	q := query.Use(s.db)
	return q.Transaction(func(tx *query.Query) error {
		if _, err := tx.Enrollment.WithContext(s.ctx).
			Where(tx.Enrollment.ClassID.Eq(classID)).
			Delete(&model.Enrollment{}); err != nil {
			return err
		}

		if _, err := tx.Class.WithContext(s.ctx).
			Where(tx.Class.ClassID.Eq(classID)).
			Delete(&model.Class{}); err != nil {
			return err
		}

		return nil
	})
}

func (s *APIService) GetStudent(studentID string) (*common.Student, error) {
	q := query.Use(s.db)
	stu, err := q.Student.WithContext(s.ctx).
		Where(q.Student.StudentID.Eq(studentID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.StudentNotFoundErr
		}
		return nil, errno.ServiceErr
	}
	major := stu.Major
	var majorPtr *string
	if major != "" {
		majorPtr = &major
	}
	return &common.Student{
		StudentID: stu.StudentID,
		Name:      stu.Name,
		Major:     majorPtr,
	}, nil
}

func (s *APIService) ListAllStudents() ([]*common.Student, error) {
	q := query.Use(s.db)
	students, err := q.Student.WithContext(s.ctx).Find()
	if err != nil {
		return nil, errno.ServiceErr
	}
	resp := make([]*common.Student, 0, len(students))
	for _, stu := range students {
		major := stu.Major
		var majorPtr *string
		if major != "" {
			majorPtr = &major
		}
		resp = append(resp, &common.Student{
			StudentID: stu.StudentID,
			Name:      stu.Name,
			Major:     majorPtr,
		})
	}
	return resp, nil
}

func (s *APIService) DeleteStudent(studentID string) error {
	q := query.Use(s.db)
	return q.Transaction(func(tx *query.Query) error {
		// 删除所有选课记录
		if _, err := tx.Enrollment.WithContext(s.ctx).
			Where(tx.Enrollment.StudentID.Eq(studentID)).
			Delete(); err != nil {
			return err
		}

		// 删除学生
		if _, err := tx.Student.WithContext(s.ctx).
			Where(tx.Student.StudentID.Eq(studentID)).
			Delete(); err != nil {
			return err
		}

		return nil
	})
}

func (s *APIService) GetClassRoster(classID string) ([]*common.RosterItem, error) {
	return s.loadRoster(classID)
}

func (s *APIService) RemoveStudentFromClass(enrollmentID string) error {
	q := query.Use(s.db)
	_, err := q.Enrollment.WithContext(s.ctx).
		Where(q.Enrollment.EnrollmentID.Eq(enrollmentID)).
		Delete()
	return err
}

func (s *APIService) RollCall(req *api.RollCallRequest) (common.RosterItem, error) {
	roster, err := s.loadRoster(req.ClassID)
	if err != nil {
		return common.RosterItem{}, err
	}

	selected, err := utils.Roll(roster, req.Mode, req.EventType)
	if err != nil {
		return common.RosterItem{}, err
	}

	if err := s.afterRollCall(selected, req); err != nil {
		return common.RosterItem{}, err
	}

	return selected, nil
}

func (s *APIService) SolveRollCall(req *api.SolveRollCallRequest) error {
	delta := calculateScoreDelta(req)
	q := query.Use(s.db)
	return q.Transaction(func(tx *query.Query) error {
		// 查找选课记录
		enrollment, err := tx.Enrollment.WithContext(s.ctx).
			Where(tx.Enrollment.EnrollmentID.Eq(req.EnrollmentID)).
			First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errno.StudentNotFoundErr
			}
			return err
		}

		var transferTarget *model.Enrollment
		if req.AnswerType == common.AnswerType_TRANSFER {
			var err error
			transferTarget, err = s.processTransfer(tx, enrollment, req)
			if err != nil {
				return err
			}
		}

		// 处理跳过权
		if req.AnswerType == common.AnswerType_SKIP {
			if enrollment.SkipRights <= 0 {
				return errno.ParamErr.WithMessage("skip rights not enough")
			}
			// 消耗跳过权
			if _, err := tx.Enrollment.WithContext(s.ctx).
				Where(tx.Enrollment.EnrollmentID.Eq(req.EnrollmentID)).
				Update(tx.Enrollment.SkipRights, enrollment.SkipRights-1); err != nil {
				return err
			}
		}

		// 更新总积分
		if _, err := tx.Enrollment.WithContext(s.ctx).
			Where(tx.Enrollment.EnrollmentID.Eq(req.EnrollmentID)).
			Update(tx.Enrollment.TotalPoints, enrollment.TotalPoints+delta); err != nil {
			return err
		}

		// 记录积分事件
		reason := fmt.Sprintf("answer=%s", req.AnswerType.String())
		metadata := datatypes.JSONMap(nil)
		if transferTarget != nil {
			reason = fmt.Sprintf("%s,target=%s", reason, transferTarget.StudentID)
			metadata = datatypes.JSONMap{
				"target_enrollment_id": transferTarget.EnrollmentID,
				"target_student_id":    transferTarget.StudentID,
			}
		}
		event := &model.ScoreEvent{
			EventID:      uuid.NewString(),
			EnrollmentID: enrollment.EnrollmentID,
			StudentID:    enrollment.StudentID,
			ClassID:      enrollment.ClassID,
			Delta:        delta,
			Reason:       reason,
			EventType:    int32(req.EventType),
			Metadata:     metadata,
		}
		if err := tx.ScoreEvent.WithContext(s.ctx).Create(event); err != nil {
			return err
		}

		return nil
	})
}

func (s *APIService) loadRoster(classID string) ([]*common.RosterItem, error) {
	q := query.Use(s.db)
	enrollments, err := q.Enrollment.WithContext(s.ctx).
		Preload(q.Enrollment.Student).
		Where(q.Enrollment.ClassID.Eq(classID)).
		Find()
	if err != nil {
		return nil, errno.QueryStudentIDsErr
	}
	if len(enrollments) == 0 {
		return nil, errno.StudentNotFoundErr
	}

	sort.SliceStable(enrollments, func(i, j int) bool {
		if enrollments[i].CallCount == enrollments[j].CallCount {
			return enrollments[i].StudentID < enrollments[j].StudentID
		}
		return enrollments[i].CallCount < enrollments[j].CallCount
	})

	roster := make([]*common.RosterItem, 0, len(enrollments))
	for _, en := range enrollments {
		roster = append(roster, toRosterItem(en))
	}
	return roster, nil
}

func (s *APIService) afterRollCall(item common.RosterItem, req *api.RollCallRequest) error {
	q := query.Use(s.db)
	return q.Transaction(func(tx *query.Query) error {
		// 更新点名次数
		newCallCount := item.EnrollmentInfo.CallCount + 1
		if _, err := tx.Enrollment.WithContext(s.ctx).
			Where(tx.Enrollment.EnrollmentID.Eq(item.EnrollmentInfo.EnrollmentID)).
			Update(tx.Enrollment.CallCount, newCallCount); err != nil {
			return err
		}

		// 每点名2次增加一个转移权
		if newCallCount%2 == 0 {
			if _, err := tx.Enrollment.WithContext(s.ctx).
				Where(tx.Enrollment.EnrollmentID.Eq(item.EnrollmentInfo.EnrollmentID)).
				UpdateSimple(tx.Enrollment.TransferRights.Add(1)); err != nil {
				return err
			}
		}

		// 记录点名记录
		record := &model.RollCallRecord{
			RecordID:     uuid.NewString(),
			ClassID:      req.ClassID,
			EnrollmentID: item.EnrollmentInfo.EnrollmentID,
			StudentID:    item.StudentInfo.StudentID,
			Mode:         int32(req.Mode),
			EventType:    int32(req.EventType),
			CreatedAt:    time.Now(),
		}

		if err := tx.RollCallRecord.WithContext(s.ctx).Create(record); err != nil {
			return err
		}

		return nil
	})
}

func (s *APIService) processTransfer(tx *query.Query, src *model.Enrollment, req *api.SolveRollCallRequest) (*model.Enrollment, error) {
	targetID := strings.TrimSpace(req.GetTargetEnrollmentID())
	if targetID == "" {
		return nil, errno.ParamErr.WithMessage("target_enrollment_id required for transfer")
	}

	// 查找目标选课记录
	target, err := tx.Enrollment.WithContext(s.ctx).
		Where(tx.Enrollment.EnrollmentID.Eq(targetID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.StudentNotFoundErr
		}
		return nil, err
	}

	// 验证
	if target.ClassID != src.ClassID {
		return nil, errno.ParamErr.WithMessage("target must be in the same class")
	}
	if target.EnrollmentID == src.EnrollmentID {
		return nil, errno.ParamErr.WithMessage("cannot transfer to self")
	}
	if src.TransferRights <= 0 {
		return nil, errno.TransferRightsNotEnough
	}

	// 消耗源学生的转移权
	if _, err := tx.Enrollment.WithContext(s.ctx).
		Where(tx.Enrollment.EnrollmentID.Eq(src.EnrollmentID)).
		Update(tx.Enrollment.TransferRights, src.TransferRights-1); err != nil {
		return nil, err
	}

	// 更新目标学生的点名次数
	nextCount := target.CallCount + 1
	if _, err := tx.Enrollment.WithContext(s.ctx).
		Where(tx.Enrollment.EnrollmentID.Eq(target.EnrollmentID)).
		Update(tx.Enrollment.CallCount, nextCount); err != nil {
		return nil, err
	}

	// 如果点名次数为偶数，增加转移权
	if nextCount%2 == 0 {
		if _, err := tx.Enrollment.WithContext(s.ctx).
			Where(tx.Enrollment.EnrollmentID.Eq(target.EnrollmentID)).
			UpdateSimple(tx.Enrollment.TransferRights.Add(1)); err != nil {
			return nil, err
		}
	}

	return target, nil
}

func toRosterItem(en *model.Enrollment) *common.RosterItem {
	major := en.Student.Major
	var majorPtr *string
	if major != "" {
		majorPtr = &major
	}
	return &common.RosterItem{
		StudentInfo: &common.Student{
			StudentID: en.Student.StudentID,
			Name:      en.Student.Name,
			Major:     majorPtr,
		},
		EnrollmentInfo: &common.Enrollment{
			EnrollmentID:   en.EnrollmentID,
			StudentID:      en.StudentID,
			ClassID:        en.ClassID,
			TotalPoints:    en.TotalPoints,
			CallCount:      en.CallCount,
			TransferRights: en.TransferRights,
			SkipRights:     en.SkipRights,
		},
	}
}

func calculateScoreDelta(req *api.SolveRollCallRequest) float64 {
	base := 0.0
	switch req.AnswerType {
	case common.AnswerType_NORMAL:
		// 正常回答：到达课堂+1，回答问题自定义0-3分
		base = 1.0 // 到达课堂的基础分
		if req.CustomScore != nil {
			// 加上回答问题的分数（可以是0-3，也可以是-1表示不准确重复问题）
			base += *req.CustomScore
		}
	case common.AnswerType_HELP:
		// 请求帮助：到达+1，准确重复问题+0.5
		base = 1.5
	case common.AnswerType_SKIP:
		// 跳过：有跳过权则不扣分（到达+1）
		base = 1.0
	case common.AnswerType_TRANSFER:
		// 转移：有转移权则不扣分（到达+1）
		base = 1.0
	default:
		base = 0
	}

	// 只有正常回答和请求帮助才应用事件加成
	if req.AnswerType == common.AnswerType_NORMAL || req.AnswerType == common.AnswerType_HELP {
		switch req.EventType {
		case common.RandomEventType_Double_Point:
			base *= 2
		case common.RandomEventType_CRAZY_THURSDAY:
			base *= 1.5
		}
	}

	return math.Round(base*100) / 100
}

// LeaderboardItem 排行榜项
type LeaderboardItem struct {
	Rank        int     `json:"rank"`
	StudentID   string  `json:"student_id"`
	Name        string  `json:"name"`
	Major       *string `json:"major,omitempty"`
	TotalPoints float64 `json:"total_points"`
	CallCount   int64   `json:"call_count"`
}

// GetLeaderboard 获取班级积分排行榜
func (s *APIService) GetLeaderboard(classID string, top int) ([]LeaderboardItem, error) {
	q := query.Use(s.db)

	// 验证班级存在
	_, err := q.Class.WithContext(s.ctx).
		Where(q.Class.ClassID.Eq(classID)).
		First()
	if err != nil {
		return nil, errno.ServiceErr
	}

	// 查询enrollments并按total_points降序排序
	enrollments, err := q.Enrollment.WithContext(s.ctx).
		Where(q.Enrollment.ClassID.Eq(classID)).
		Order(q.Enrollment.TotalPoints.Desc(), q.Enrollment.CallCount.Desc()).
		Preload(q.Enrollment.Student).
		Find()
	if err != nil {
		return nil, errno.ServiceErr
	}

	// 限制返回数量
	limit := len(enrollments)
	if top > 0 && top < limit {
		limit = top
	}

	result := make([]LeaderboardItem, 0, limit)
	for i := 0; i < limit; i++ {
		en := enrollments[i]
		var major *string
		if en.Student.Major != "" {
			major = &en.Student.Major
		}
		result = append(result, LeaderboardItem{
			Rank:        i + 1,
			StudentID:   en.StudentID,
			Name:        en.Student.Name,
			Major:       major,
			TotalPoints: en.TotalPoints,
			CallCount:   en.CallCount,
		})
	}

	return result, nil
}

// ClassStats 班级统计信息
type ClassStats struct {
	TotalStudents      int              `json:"total_students"`
	TotalCalls         int64            `json:"total_calls"`
	AveragePoints      float64          `json:"average_points"`
	PointsDistribution map[string]int   `json:"points_distribution"` // 积分区间分布
	CallFrequency      map[string]int64 `json:"call_frequency"`      // 点名次数分布
}

// GetClassStats 获取班级统计信息
func (s *APIService) GetClassStats(classID string) (*ClassStats, error) {
	q := query.Use(s.db)

	// 验证班级存在
	_, err := q.Class.WithContext(s.ctx).
		Where(q.Class.ClassID.Eq(classID)).
		First()
	if err != nil {
		return nil, errno.ServiceErr
	}

	// 查询所有enrollments
	enrollments, err := q.Enrollment.WithContext(s.ctx).
		Where(q.Enrollment.ClassID.Eq(classID)).
		Find()
	if err != nil {
		return nil, errno.ServiceErr
	}

	if len(enrollments) == 0 {
		return &ClassStats{
			TotalStudents:      0,
			TotalCalls:         0,
			AveragePoints:      0,
			PointsDistribution: make(map[string]int),
			CallFrequency:      make(map[string]int64),
		}, nil
	}

	stats := &ClassStats{
		TotalStudents:      len(enrollments),
		PointsDistribution: make(map[string]int),
		CallFrequency:      make(map[string]int64),
	}

	var totalPoints float64
	for _, en := range enrollments {
		stats.TotalCalls += en.CallCount
		totalPoints += en.TotalPoints

		// 积分区间分布 (0-10, 10-20, 20-30, ...)
		bucket := int(en.TotalPoints/10) * 10
		bucketKey := fmt.Sprintf("%d-%d", bucket, bucket+10)
		stats.PointsDistribution[bucketKey]++

		// 点名次数分布
		callKey := fmt.Sprintf("%d", en.CallCount)
		stats.CallFrequency[callKey]++
	}

	stats.AveragePoints = math.Round(totalPoints/float64(len(enrollments))*100) / 100

	return stats, nil
}
