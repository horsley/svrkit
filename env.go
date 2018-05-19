package svrkit

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

//SwitchPwd 切换工作目录到程序文件所在目录，防止路径问题
func SwitchPwd() {
	var wd string

	if runtime.GOOS == "linux" {
		if exe, err := os.Readlink("/proc/self/exe"); err != nil {
			log.Println("switchPwd:read exe path err:", err)
			os.Exit(1)
		} else {
			wd = filepath.Dir(exe)

		}
	} else {
		wd = filepath.Dir(os.Args[0])
	}

	if err := os.Chdir(wd); err != nil {
		log.Println("switchPwd:chdir to path:", wd, " err:", err)
		os.Exit(1)
	}
}
