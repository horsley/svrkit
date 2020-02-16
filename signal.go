package svrkit

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

//SignalHandler 描述一个信号处理项
type SignalHandler struct {
	Cmd  string
	Desc string
	//收到信号之后的处理函数，返回一个结果提示文本
	Action func(args ...string) string
}

//SignalManager 信号处理器
type SignalManager struct {
	//Listen 监听地址
	Listen string
	//信号触发的参数 flag
	ArgFlag string
	items   map[string]SignalHandler
}

//AddHandler 注册处理项
func (s *SignalManager) AddHandler(handlers ...SignalHandler) {
	if s.items == nil {
		s.items = make(map[string]SignalHandler)
	}
	for _, v := range handlers {
		s.items[v.Cmd] = v
	}
}

//IsSignalSender 会检测 os.Args, 如果发现ArgFlag 则会把下一参数作为 signal 发送到 Listen 地址
//如果没有下一参数，会打印已注册的 signal 说明信息，返回 true，如果没有找到ArgFlag，返回 false
// 示例代码
// signal := svrkit.SignalManager{
// 	Listen:  "127.0.0.1:18909",
// 	ArgFlag: "-s",
// }
// signal.AddHandler(svrkit.SignalHandler{
// 	Cmd:  "reload",
// 	Desc: "reload cache",
// 	Action: func() string {
// 		log.Println("received signal reload")
// 		//logic here
// 		return "reload ok"
// 	},
// })
// if signal.IsSignalSender() {
// 	os.Exit(0)
// } else {
// 	go signal.ListenSignal()
// }
func (s *SignalManager) IsSignalSender() bool {
	for i, v := range os.Args {
		if v == s.ArgFlag {
			if i != len(os.Args)-1 {
				signal := os.Args[i+1]
				if s.Listen != "" {
					resp, err := http.PostForm("http://"+s.Listen, url.Values{
						"signal": []string{signal},
						"args":   os.Args[i+2:],
					})
					if err != nil {
						fmt.Println("send signal err:", err)
					} else {
						defer resp.Body.Close()
						respText, _ := ioutil.ReadAll(resp.Body)
						fmt.Println(string(respText))
					}
				}
			} else { //flag 之后无其他参数，输出帮助
				fmt.Println(s.helpText())
			}
			return true
		}
	}
	return false
}

//ListenSignal 开始在 Listen 地址上监听 signal，请用一个独立携程做这个事情
func (s *SignalManager) ListenSignal() {
	server := &http.Server{Addr: s.Listen, Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			return
		}
		signal := r.FormValue("signal")
		for k, v := range s.items {
			if k == signal {
				rw.Write([]byte(v.Action(r.Form["args"]...)))
				return
			}
		}
		//not found
		rw.Write([]byte(fmt.Sprintf("signal '%s' not defined\n\n%s", signal, s.helpText())))
	})}
	log.Fatalln("signal listen err:", server.ListenAndServe())
}

func (s *SignalManager) helpText() string {
	result := fmt.Sprintln(os.Args[0], s.ArgFlag, "[signal]")
	result += fmt.Sprintln("Available signals:")
	for k, v := range s.items {
		result += fmt.Sprintln("    ", k, ":", v.Desc)
	}
	return result
}
