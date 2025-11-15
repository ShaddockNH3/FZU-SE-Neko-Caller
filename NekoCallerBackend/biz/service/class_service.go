package service

import (
	"FZUSENekoCaller/biz/dal/model"
	"FZUSENekoCaller/biz/dal/query"
	"FZUSENekoCaller/biz/model/api"
	"FZUSENekoCaller/pkg/errno"
	"context"

	"github.com/google/uuid"
)

type ImportService struct {
	ctx context.Context
}

func NewImportService(ctx context.Context) *ImportService {
	return &ImportService{ctx: ctx}
}

func (s *ImportService) ImportClassData(req *api.ImportDataRequest) error {
	return query.Q.Transaction(func(tx *query.Query) error {
		classID := uuid.New().String()
		newClass := &model.Class{
			ClassID:   classID,
			ClassName: req.ClassName,
		}
		if err := tx.Class.WithContext(s.ctx).Create(newClass); err != nil {
			return errno.CreateClassErr
		}

		for _, studentReq := range req.Students {
			student, err := tx.Student.WithContext(s.ctx).
				Where(tx.Student.StudentID.Eq(studentReq.StudentID)).
				Attrs(tx.Student.StudentID.Value(studentReq.StudentID),
					tx.Student.Name.Value(studentReq.Name),
					tx.Student.Major.Value(*studentReq.Major)).
				FirstOrCreate()

			if err != nil {
				return errno.CreateClassOrStudentErr
			}

			enrollmentID := uuid.New().String()
			newEnrollment := &model.Enrollment{
				EnrollmentID: enrollmentID,
				StudentID:    student.StudentID,
				ClassID:      classID,
			}
			if err := tx.Enrollment.WithContext(s.ctx).Create(newEnrollment); err != nil {
				return errno.StudentJoinClassErr
			}
		}

		return nil
	})
}
