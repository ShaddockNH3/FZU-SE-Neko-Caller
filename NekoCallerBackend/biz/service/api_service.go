package service

import (
	"FZUSENekoCaller/biz/dal/model"
	"FZUSENekoCaller/biz/dal/query"
	"FZUSENekoCaller/biz/model/api"
	"FZUSENekoCaller/biz/model/common"
	"FZUSENekoCaller/pkg/errno"
	"FZUSENekoCaller/pkg/utils"
	"context"
)

type APIService struct {
	ctx context.Context
}

func NewAPIService(ctx context.Context) *APIService {
	return &APIService{ctx: ctx}
}

func (s *APIService) GetClass(classID string) (*common.Class, error) {
	var class *model.Class

	q := query.Q

	class, err := q.Class.WithContext(s.ctx).
		Where(q.Class.ClassID.
			Eq(classID)).
		First()
	if err != nil {
		return &common.Class{}, errno.QueryStudentIDsErr
	}

	var student_ids []string
	err = q.Enrollment.WithContext(s.ctx).
		Where(q.Enrollment.ClassID.Eq(classID)).
		Pluck(q.Enrollment.StudentID, &student_ids)
	if err != nil {
		return &common.Class{}, errno.QueryStudentIDsErr
	}

	commonClass := &common.Class{
		ClassID:    class.ClassID,
		ClassName:  class.ClassName,
		StudentIds: student_ids,
	}

	return commonClass, nil
}

func (s *APIService) ListClasses() ([]*common.Class, error) {
	return []*common.Class{}, nil
}

func (s *APIService) DeleteClass(classID string) error {
	return nil
}

func (s *APIService) GetStudent(studentID string) (*common.Student, error) {
	var student *model.Student

	q := query.Q

	student, err := q.Student.WithContext(s.ctx).
		Where(q.Student.StudentID.
			Eq(studentID)).
		First()
	if err != nil {
		return &common.Student{}, errno.QueryStudentIDsErr
	}

	commonStudent := &common.Student{
		StudentID: student.StudentID,
		Name:      student.Name,
		Major:     &student.Major,
	}

	return commonStudent, nil
}

func (s *APIService) ListAllStudents() ([]*common.Student, error) {
	return []*common.Student{}, nil
}

func (s *APIService) DeleteStudent(studentID string) error {
	return nil
}

func (s *APIService) GetClassRoster(classID string) ([]*common.RosterItem, error) {
	var enrollments []*model.Enrollment

	q := query.Q
	enrollments, err := q.Enrollment.WithContext(s.ctx).
		Preload(q.Enrollment.Student).
		Where(q.Enrollment.ClassID.Eq(classID)).
		Find()
	if err != nil {
		return []*common.RosterItem{}, errno.QueryStudentIDsErr
	}

	roster := make([]*common.RosterItem, 0, len(enrollments))
	for _, en := range enrollments {
		item := &common.RosterItem{
			StudentInfo: &common.Student{
				StudentID: en.Student.StudentID,
				Name:      en.Student.Name,
				Major:     &en.Student.Major,
			},
			EnrollmentInfo: &common.Enrollment{
				EnrollmentID:   en.EnrollmentID,
				StudentID:      en.StudentID,
				ClassID:        en.ClassID,
				TotalPoints:    en.TotalPoints,
				CallCount:      en.CallCount,
				TransferRights: en.TransferRights,
			},
		}
		roster = append(roster, item)
	}

	return roster, nil
}

func (s *APIService) RemoveStudentFromClass(studentID string, classID string) error {
	return nil
}

func (s *APIService) RollCall(req *api.RollCallRequest) (common.RosterItem, error) {
	var enrollments []*model.Enrollment

	q := query.Q
	enrollments, err := q.Enrollment.WithContext(s.ctx).
		Preload(q.Enrollment.Student).
		Where(q.Enrollment.ClassID.Eq(req.ClassID)).
		Find()
	if err != nil {
		return common.RosterItem{}, errno.QueryStudentIDsErr
	}

	roster := make([]*common.RosterItem, 0, len(enrollments))
	for _, en := range enrollments {
		item := &common.RosterItem{
			StudentInfo: &common.Student{
				StudentID: en.Student.StudentID,
				Name:      en.Student.Name,
				Major:     &en.Student.Major,
			},
			EnrollmentInfo: &common.Enrollment{
				EnrollmentID:   en.EnrollmentID,
				StudentID:      en.StudentID,
				ClassID:        en.ClassID,
				TotalPoints:    en.TotalPoints,
				CallCount:      en.CallCount,
				TransferRights: en.TransferRights,
			},
		}
		roster = append(roster, item)
	}

	currentRoster, err := utils.Roll(roster, req.Mode, req.EventType)

	if err != nil {
		return common.RosterItem{}, err
	}

	return currentRoster, nil
}
