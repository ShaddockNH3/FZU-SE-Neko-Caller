//go:build ignore

package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"

	"FZUSENekoCaller/biz/dal/model"
	"FZUSENekoCaller/pkg/constants"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		constants.MySQLUser,
		constants.MySQLPassword,
		constants.MySQLHost,
		constants.MySQLPort,
		constants.MySQLDatabase,
	)

	gormdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./biz/dal/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(gormdb)

	g.ApplyBasic(
		model.Student{},
		model.Class{},
		model.Enrollment{},
		model.RollCallRecord{},
		model.ScoreEvent{},
	)

	g.Execute()
}
