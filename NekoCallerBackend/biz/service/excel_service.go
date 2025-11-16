package service

import (
	"context"
	"fmt"
	"io"

	"FZUSENekoCaller/biz/model/api"
	"FZUSENekoCaller/biz/model/common"

	"github.com/xuri/excelize/v2"
)

type ExcelService struct {
	ctx context.Context
}

func NewExcelService(ctx context.Context) *ExcelService {
	return &ExcelService{ctx: ctx}
}

// ParseExcelToImportRequest 解析Excel文件为导入请求
func (s *ExcelService) ParseExcelToImportRequest(reader io.Reader, className string) (*api.ImportDataRequest, error) {
	// 读取 Excel
	xlsx, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("读取Excel失败: %w", err)
	}
	defer xlsx.Close()

	// 获取第一个工作表
	sheets := xlsx.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("Excel文件为空")
	}

	rows, err := xlsx.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %w", err)
	}

	if len(rows) <= 1 {
		return nil, fmt.Errorf("Excel文件中没有学生数据")
	}

	// 解析学生数据 (跳过标题行)
	// 预期列: 学号 | 姓名 | 专业
	students := make([]*common.Student, 0)
	for i, row := range rows {
		if i == 0 {
			// 跳过标题行
			continue
		}
		if len(row) < 2 {
			continue
		}

		studentID := row[0]
		name := row[1]
		var major *string
		if len(row) >= 3 && row[2] != "" {
			major = &row[2]
		}

		if studentID == "" || name == "" {
			continue
		}

		students = append(students, &common.Student{
			StudentID: studentID,
			Name:      name,
			Major:     major,
		})
	}

	if len(students) == 0 {
		return nil, fmt.Errorf("未找到有效的学生数据")
	}

	return &api.ImportDataRequest{
		ClassName: className,
		Students:  students,
	}, nil
}

// GenerateRosterExcel 生成花名册Excel文件
func (s *ExcelService) GenerateRosterExcel(className string, roster []*common.RosterItem) (*excelize.File, error) {
	// 创建Excel文件
	f := excelize.NewFile()
	sheetName := "Sheet1"
	f.SetSheetName(f.GetSheetName(0), sheetName)

	// 写入标题行
	headers := []string{"学号", "姓名", "专业", "随机点名次数", "总积分"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据
	for i, item := range roster {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), item.StudentInfo.StudentID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.StudentInfo.Name)
		if item.StudentInfo.Major != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), *item.StudentInfo.Major)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "")
		}
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), item.EnrollmentInfo.CallCount)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), item.EnrollmentInfo.TotalPoints)
	}

	// 设置列宽
	f.SetColWidth(sheetName, "A", "A", 15)
	f.SetColWidth(sheetName, "B", "B", 15)
	f.SetColWidth(sheetName, "C", "C", 20)
	f.SetColWidth(sheetName, "D", "D", 15)
	f.SetColWidth(sheetName, "E", "E", 15)

	return f, nil
}
