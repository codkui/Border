package main

import (
	"fmt"

	"./conn"
	"./queue"
	sv "./serviceCenter"
	//"archive/zip"
)

//Que 服务代理中心核心队列
var Que queue.EsQueue

type geome interface {
	Area() float64
	Perim() float64
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
	go test()
	conn.Start(8088, 8087)

}

func test() {
	for {
		a := <-conn.InputChan
		fmt.Println("从队列取得请求" + a.API)
		fmt.Println(a.Data)
		go func() {
			res, err := sv.DoService(a)
			if err != nil {
				fmt.Println(err)
				//conn.ConOutputChan[a.Index] <- "访问服务失败，服务不存在"
				conn.Output(a, "访问服务失败，服务不存在或存在异常")
				return
			}
			fmt.Println("返回数据" + res)
			//解析器，解析请求到各个服务，并保持连接，返回数据时压入返回队列
			conn.Output(a, res)
		}()

	}
}

/*代码废弃区域，仅供参考*/
//Que := queue.NewQueue(8088)
// b := UseType{}
// fmt.Println(b)
// a := map[string]interface{}{"aaa": 1, "bbb": "cc"}
// fmt.Println(a)
// str, err := msgpack.Marshal(a)
// if err != nil {
// 	fmt.Println(err)
// 	return
// }
//str := gencode.Marshal(a)
// fmt.Println("start running!")
// re := rect{3, 5}
// var _ geome = rect{}
// fmt.Println(re.area())
// go test()
// conn.Start("127.0.0.1", 8080)
// fmt.Println(str)
// var testStr map[string]interface{}
// err = msgpack.Unmarshal(str, &testStr)
// fmt.Printf("%v %#v\n", err, testStr)
