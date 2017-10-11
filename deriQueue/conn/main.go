package conn

import (
	"NSnow/public"
	"fmt"

	"./http"
	"./tcp"
)

//InputChan 收到的请求队列
var InputChan chan public.UseType

type conn interface {
	Respone() bool
	Reuquest() bool
	Start(string, int16) bool
}

func init() {
	InputChan = make(chan public.UseType, 3000)
	
}

//Start 启动服务 传递 http端口，tcp端口，udp端口
func Start(httpPort int, tcpPort int) {
	go http.Start("", httpPort)
	go tcp.Start("", tcpPort)
	inputLoad()
	
}

func inputLoad() {
	fmt.Println("监听启动")
	a := public.UseType{}
	for {
		select {
		case a = <-http.ConInputChan:
		case a = <-tcp.InputChan:
		}
		fmt.Println("接收到总统到数据")
		InputChan <- a
	}
}

//Output 返回数据，传入原始数据和返回的数据
func Output(useData public.UseType, data interface{}) {
	switch useData.Con {
	case "HTTP":
		http.ConOutputChan[useData.Index] <- data
		break
	case "TCP":
		tcp.OutputChan[useData.Index] <- data
		break
	default:
		break
	}
}
