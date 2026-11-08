package constants

import (
	"log"
	"os"
	"time"

	// 内嵌时区数据库：容器里常见只有 UTC（alpine / scratch），
	// 没这个的话 LoadLocation("Asia/Shanghai") 会失败，事件条件就悄悄退回 UTC。
	_ "time/tzdata"
)

// AppLocation 判断随机事件条件时按哪个时区读「现在几点」。
//
// 默认 Asia/Shanghai：1024 福报看秒、质数孤独看分钟、疯狂星期四看星期几，
// 这些在学生眼里都是「北京时间」。部署到别的时区用 APP_TIMEZONE 覆盖。
var AppLocation = loadLocation(envStr("APP_TIMEZONE", "Asia/Shanghai"))

func loadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		log.Printf("[CONFIG] 时区 %q 不可用（%v），回退到服务器本地时区 %s", name, err, time.Local)
		return time.Local
	}
	return loc
}

func envStr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
