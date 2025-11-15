package service

import (
	"FZUSENekoCaller/biz/model/api"
	"FZUSENekoCaller/biz/model/common"
	"context"
)

type APIService struct {
	ctx context.Context
}

func NewAPIService(ctx context.Context) *APIService {
	return &APIService{ctx: ctx}
}

func (s *APIService) GetClass(classID string) (*common.Class, error) {
	return &common.Class{}, nil
}

func (s *APIService) ListClasses() ([]*common.Class, error) {
	return []*common.Class{}, nil
}

func (s *APIService) DeleteClass(classID string) error {
	return nil
}

func (s *APIService) GetStudent(studentID string) (*common.Student, error) {
	return &common.Student{}, nil
}

func (s *APIService) ListAllStudents() ([]*common.Student, error) {
	return []*common.Student{}, nil
}

func (s *APIService) DeleteStudent(studentID string) error {
	return nil
}

func (s *APIService) GetClassRoster(classID string) ([]*common.RosterItem, error) {
	return []*common.RosterItem{}, nil
}

func (s *APIService) RemoveStudentFromClass(studentID string, classID string) error {
	return nil
}

func (s *APIService) RollCall(req *api.RollCallRequest) (*api.RollCallResponse, error) {
	return &api.RollCallResponse{}, nil
}
