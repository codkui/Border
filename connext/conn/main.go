package conn

import (
	"NSnow/public"
	"fmt"
	"strconv"

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
		fmt.Println("接收到综合接受通道数据")
		InputChan <- a
	}
}

//Output 返回数据，传入原始数据和返回的数据
func Output(useData public.UseType, data interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	switch useData.Con {
	case "HTTP":
		http.ConOutputChan[useData.Index] <- data
		break
	case "TCP":
		demo := public.TCPType{}
		v, _ := strconv.Atoi(useData.Header["APICode"])
		demo.API = int32(v)
		v, _ = strconv.Atoi(useData.Header["localIndex"])
		demo.Index = int32(v)
		demo.Data = data
		if c, ok := tcp.OutputChan[useData.Index]; ok {
			fmt.Println("写入返回数据到tcp队列")
			c <- demo
		}
		break
	default:
		break
	}
}
