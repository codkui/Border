package main

import (
	"NSnow/public"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

func main() {

	// var a int64 = 8
	// b := public.Int64ToBytes(a)
	// c := public.BytesToInt64(b)
	// fmt.Println(c)
	// return
	conn, err := net.Dial("tcp", "devel.skylinuav.com:8096")
	//conn, err := net.Dial("tcp", "127.0.0.1:8096")
	demo := map[string]interface{}{"GGA": "$GPGGA,000001,3112.518576,N,12127.901251,E,1,8,1,0,M,-32,M,3,0*4B"}
	encodeData, _ := json.Marshal(demo)
	respData := []byte{}
	apiByte := public.Int32ToBytes(int32(5))
	indexByte := public.Int32ToBytes(int32(3))
	lenByte := public.Int32ToBytes(int32(len(encodeData)))
	typeByte := public.Int32ToBytes(int32(2))
	dataIndex := public.Int32ToBytes(int32(0))
	dataLast := public.Int32ToBytes(int32(1))
	//fmt.Println(lenByte)
	//fmt.Println(public.BytesToInt32(lenByte))

	respData = append(respData, apiByte...)
	respData = append(respData, indexByte...)
	respData = append(respData, lenByte...)
	respData = append(respData, typeByte...)
	respData = append(respData, dataIndex...)
	respData = append(respData, dataLast...)
	respData = append(respData, encodeData...)
	// var b interface{}
	// _ = msgpack.Unmarshal(encodeData, &b)

	// fmt.Println(b)
	// return
	if err != nil {
		fmt.Println("连接错误")
		fmt.Println(err)
		return
	}
	a := 1
	defer conn.Close()
	go rev(conn)
	for {
		a++
		//if a%1000 == 0 {
		fmt.Println("send")
		fmt.Println(a)
		fmt.Println(respData)
		//}
		time.Sleep(1000000000)

		conn.Write(respData)

	}
}

func rev(conn net.Conn) {
	buff := make([]byte, 24)
	b := 1
	var dataLen int32
	for {
		_, err := conn.Read(buff)

		if err != nil {
			fmt.Println(err)
			return
		}
		b++
		if b%1000 == 0 {
			fmt.Println("rev")
			fmt.Println(b)
		}
		dataLen = public.BytesToInt32(buff[8:12])
		//fmt.Println(buff[8:12])
		fmt.Println(dataLen)
		if dataLen > 0 {
			if dataLen>int32(2048*1024){
				dataLen=int32(2048*1024)
			}
			relData := make([]byte, int(dataLen))
			n, err := conn.Read(relData)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if int32(n) != dataLen {

				fmt.Println("读取长度有误")
				fmt.Println(n)
				//continue
			}
			var decodeData interface{}
			_ = json.Unmarshal(relData[0:n], &decodeData)
			fmt.Println("接收到数据长度：" + strconv.Itoa(int(dataLen)))
			fmt.Println("接收到数据真实长度：" + strconv.Itoa(n))
			fmt.Println(relData[:n])
			fmt.Println(decodeData)
		}

		//fmt.Println(buff[0:n])
	}
}
