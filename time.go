package svrkit

import (
	"time"
)

//TimeSleepLoop 睡一段时间循环,参数2决定是否跳过第一次睡眠在调用的时候立马 action，action返回用于终止循环
func TimeSleepLoop(duration time.Duration, sleepFirst bool, action func() (shouldBreak bool)) {
	for {
		if sleepFirst {
			time.Sleep(duration)
		}
		sleepFirst = true

		shouldBreak := action()
		if shouldBreak {
			break
		}
	}
}
