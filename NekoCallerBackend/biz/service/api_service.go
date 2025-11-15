package service

import (
	"FZUSENekoCaller/biz/model/common"
	"context"
)

type APIService struct {
	ctx context.Context
}

func NewAPIService(ctx context.Context) *APIService {
	return &APIService{ctx: ctx}
}

func (s *APIService) ListClasses() ([]*common.Class, error) {
	return []*common.Class{}, nil
}

func (s *APIService) ListAllStudents() ([]*common.Student, error) {
	return []*common.Student{}, nil
}
