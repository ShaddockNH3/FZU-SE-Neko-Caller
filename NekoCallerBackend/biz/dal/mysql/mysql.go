package mysql

import (
	"fmt"
	"log"
	"os"
	"sync"

	"FZUSENekoCaller/biz/dal/model"
	"FZUSENekoCaller/biz/dal/query"
	"FZUSENekoCaller/pkg/constants"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	once sync.Once
	db   *gorm.DB
)

// Init establishes the global GORM connection and wires it into gorm-gen query helpers.
func Init() *gorm.DB {
	once.Do(func() {
		dsn := buildDSN()
		gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Silent),
			NamingStrategy:                           schema.NamingStrategy{SingularTable: false},
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			log.Fatalf("failed to connect mysql: %v", err)
		}

		if err := gormDB.AutoMigrate(
			&model.Student{},
			&model.Class{},
			&model.Enrollment{},
			&model.RollCallRecord{},
			&model.ScoreEvent{},
		); err != nil {
			log.Fatalf("auto migrate failed: %v", err)
		}

		query.SetDefault(gormDB)
		db = gormDB
	})

	return db
}

// GetDB returns the initialized *gorm.DB instance.
func GetDB() *gorm.DB {
	if db == nil {
		return Init()
	}
	return db
}

func buildDSN() string {
	user := getEnv("MYSQL_USER", constants.MySQLUser)
	password := getEnv("MYSQL_PASSWORD", constants.MySQLPassword)
	host := getEnv("MYSQL_HOST", constants.MySQLHost)
	port := getEnv("MYSQL_PORT", constants.MySQLPort)
	database := getEnv("MYSQL_DATABASE", constants.MySQLDatabase)

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		database,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
