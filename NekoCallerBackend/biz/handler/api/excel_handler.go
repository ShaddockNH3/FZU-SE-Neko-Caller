package api

import (
	"context"
	"fmt"

	"FZUSENekoCaller/biz/model/common"
	"FZUSENekoCaller/biz/service"
	"FZUSENekoCaller/pkg/constants"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// ImportClassDataFromExcel 从Excel导入班级数据
// @router /v1/import/excel [POST]
func ImportClassDataFromExcel(ctx context.Context, c *app.RequestContext) {
	resp := new(common.BaseResponse)

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		resp.Code = constants.CodeFailed
		resp.Message = "获取上传文件失败: " + err.Error()
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	// 获取班级名称
	className := c.PostForm("class_name")
	if className == "" {
		resp.Code = constants.CodeFailed
		resp.Message = "班级名称不能为空"
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	// 打开Excel文件
	src, err := file.Open()
	if err != nil {
		resp.Code = constants.CodeFailed
		resp.Message = "打开文件失败: " + err.Error()
		c.JSON(consts.StatusInternalServerError, resp)
		return
	}
	defer src.Close()

	// 调用ExcelService解析Excel
	excelService := service.NewExcelService(ctx)
	importReq, err := excelService.ParseExcelToImportRequest(src, className)
	if err != nil {
		resp.Code = constants.CodeFailed
		resp.Message = err.Error()
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	// 调用服务层导入
	importService := service.NewImportService(ctx)
	if err := importService.ImportClassData(importReq); err != nil {
		resp.Code = constants.CodeFailed
		resp.Message = err.Error()
		c.JSON(consts.StatusInternalServerError, resp)
		return
	}

	resp.Code = constants.CodeSuccess
	resp.Message = fmt.Sprintf("成功导入班级 %s，共 %d 名学生", className, len(importReq.Students))

	c.JSON(consts.StatusOK, resp)
}

// ExportClassRoster 导出班级花名册为Excel
// @router /v1/classes/:class_id/export [GET]
func ExportClassRoster(ctx context.Context, c *app.RequestContext) {
	resp := new(common.BaseResponse)

	classID := c.Param("class_id")
	if classID == "" {
		resp.Code = constants.CodeFailed
		resp.Message = "班级ID不能为空"
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	// 获取班级信息
	apiService := service.NewAPIService(ctx)
	class, err := apiService.GetClass(classID)
	if err != nil {
		resp.Code = constants.CodeFailed
		resp.Message = "班级不存在: " + err.Error()
		c.JSON(consts.StatusNotFound, resp)
		return
	}

	// 获取花名册
	roster, err := apiService.GetClassRoster(classID)
	if err != nil {
		resp.Code = constants.CodeFailed
		resp.Message = "获取花名册失败: " + err.Error()
		c.JSON(consts.StatusInternalServerError, resp)
		return
	}

	// 调用ExcelService生成Excel
	excelService := service.NewExcelService(ctx)
	f, err := excelService.GenerateRosterExcel(class.ClassName, roster)
	if err != nil {
		resp.Code = constants.CodeFailed
		resp.Message = "生成Excel文件失败: " + err.Error()
		c.JSON(consts.StatusInternalServerError, resp)
		return
	}
	defer f.Close()

	// 生成文件缓冲区
	buffer, err := f.WriteToBuffer()
	if err != nil {
		resp.Code = constants.CodeFailed
		resp.Message = "生成Excel文件失败: " + err.Error()
		c.JSON(consts.StatusInternalServerError, resp)
		return
	}

	// 设置响应头
	filename := fmt.Sprintf("%s-积分详单.xlsx", class.ClassName)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", filename))
	c.Data(consts.StatusOK, "application/octet-stream", buffer.Bytes())
}
