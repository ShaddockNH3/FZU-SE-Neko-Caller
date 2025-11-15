package service

import (
	"context"
	"errors"

	"FZUSENekoCaller/biz/dal/model"
	"FZUSENekoCaller/biz/dal/query"
	"FZUSENekoCaller/biz/model/api"
	"FZUSENekoCaller/pkg/errno"

	"github.com/google/uuid"
	"gorm.io/gen/field"
	"gorm.io/gorm"
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
			attrs := []field.AssignExpr{
				tx.Student.StudentID.Value(studentReq.StudentID),
				tx.Student.Name.Value(studentReq.Name),
			}
			if studentReq.Major != nil {
				attrs = append(attrs, tx.Student.Major.Value(*studentReq.Major))
			}

			student, err := tx.Student.WithContext(s.ctx).
				Where(tx.Student.StudentID.Eq(studentReq.StudentID)).
				Assign(tx.Student.Name.Value(studentReq.Name)).
				Attrs(attrs...).
				FirstOrCreate()

			if err != nil {
				return errno.CreateClassOrStudentErr
			}

			_, err = tx.Enrollment.WithContext(s.ctx).
				Where(tx.Enrollment.ClassID.Eq(classID), tx.Enrollment.StudentID.Eq(student.StudentID)).
				First()
			if err == nil {
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
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
