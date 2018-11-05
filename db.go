package svrkit

import (
	"database/sql"
	"log"
	"time"
)

//DBKeepaliveDefaultInterval 默认的keepalive间隔 3小时
var DBKeepaliveDefaultInterval = 3 * 60 * 60 * time.Second

// DBKeepalive 定时ping db保持链接激活，mysql 默认配置空闲链接断开超时为 28800s 等于8小时
// 空闲被断开后的第一次请求会错误（之后db库会自动重连）,为了规避这种情况我们需要以小于db链接断开的超时间隔
// 去做一个ping操作保活，请用独立协程调用本方法 go svrkit.DBKeepalive(db, xxx)
func DBKeepalive(db *sql.DB, interval time.Duration) {
	if interval.Nanoseconds() == 0 {
		interval = DBKeepaliveDefaultInterval
	}

	for {
		time.Sleep(interval)
		err := db.Ping()
		if err != nil {
			log.Println("DBKeepalive ping err:", err)
		}
	}
}
