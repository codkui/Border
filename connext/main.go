package main

import "fmt"
import "./conn"

//import tcp "./conn/tcp"
import "NSnow/public"

type geome interface {
	area() float64
	perim() float64
}

type rect struct {
	width, height float64
}

func (r rect) area() float64 {
	return r.width * r.height
}
func (r rect) perim() float64 {
	return r.width
}
func main() {
	fmt.Println("start running!")
	re := rect{3, 5}
	var _ geome = rect{}
	fmt.Println(re.area())

	go conn.Start(8095, 8096)
	//go tcp.Test()

	test()
	fmt.Println("启动完毕")

}

func test() {
	a := public.UseType{}
	for {
		// select {
		// case a = <-conn.ConInputChan:
		// case a = <-tcp.InputChan:
		// }
		a = <-conn.InputChan
		fmt.Println("开始处理请求" + a.API)
		//conn.ConOutputChan[a.Index] <- a
		go doAct(a)
	}
	//a = public.UseType{}
}

func doAct(data public.UseType) {

	resData, err := public.GetAnswer("127.0.0.1:8088", "HTTP", "/Request", data)
	if err != nil {
		//conn.ConOutputChan[data.Index] <- "请求错误"
		fmt.Println("接收到请求处理失败")
		conn.Output(data, "{\"data\":\"请求错误\"}")
		fmt.Println("接收到请求处理失败")
		fmt.Println(data.API)
		return
	}
	fmt.Println("请求处理完毕")
	fmt.Println(resData)
	conn.Output(data, resData)
	// switch data.Con {
	// case "HTTP":
	// 	conn.ConOutputChan[data.Index] <- resData
	// 	break
	// case "TCP":
	// 	tcp.OutChanList[data.Index] <- resData
	// 	break
	// default:
	// 	break
	// }

}
