package tcp

import "net"
import "fmt"
import "NSnow/public"
import "strconv"
import "encoding/json"

//InputChan 接收的队列
var InputChan chan public.UseType

//OutputChan 输出的队列
var OutputChan map[string]chan interface{}

var maxConnIndex int64 = 1

func init() {
	InputChan = make(chan public.UseType, 3000)
	OutputChan = make(map[string]chan interface{}, 5)
}

//Start 启动服务
func Start(ip string, port int) {
	ser, err := net.Listen("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("tcp 建立，端口 " + strconv.Itoa(port))
	for {
		newConn, err := ser.Accept()
		if err != nil {
			continue
		}

		go recvConnMsg(newConn)
	}

}

func recvConnMsg(conn net.Conn) {
	//  var buf [50]byte
	maxConnIndex++
	index := strconv.FormatInt(maxConnIndex, 10)
	c := make(chan interface{})
	OutputChan[index] = c
	buf := make([]byte, 50)
	useData := public.UseType{Index: index, Con: "TCP"}
	//a := 1
	defer func() {
		conn.Close()
		delete(OutputChan, index)
		close(c)
	}()
	go ResponseData(conn, c)
	for {
		n, err := conn.Read(buf)

		if err != nil {
			fmt.Println("conn closed")
			return
		}
		useData.Data = buf[:n]
		InputChan <- useData
		//InputChan
		//fmt.Println("recv msg:", buf[0:n])
		//conn.Write(buf[0:n])
		//fmt.Println("recv msg:", string(buf[0:n]))
		//fmt.Println(buf[0:n])

	}
}
func loadMsg() {}
func ResponseData(conn net.Conn, c chan interface{}) {
	for {
		data, isClose := <-c
		fmt.Println("开始返回数据")
		if !isClose {
			fmt.Println("通道已关闭")
			return
		}
		fmt.Println("通道正常")
		jsonStr, err := json.Marshal(data)
		fmt.Println(jsonStr)
		if err != nil {
			fmt.Println(err)
			continue
		}
		conn.Write(jsonStr)
	}
}
