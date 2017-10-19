package tcp

import (
	"NSnow/public"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/vmihailenco/msgpack"
)

/*
系统占用码
1 已接收
2 请求重发
3 心跳包
4
*/

/*
GOOS=linux GOARCH=amd64 go build 交叉编译
*/
//InputChan 接收的队列
var InputChan chan public.UseType

//OutputChan 输出的队列
var OutputChan map[string]chan interface{}

var maxConnIndex int64 = 1

var resCode int32 = 1
var zeroLen int32 = 0
var restCode int32 = 2
var dataType int32 = 1
var dataIndex int32 = 0
var dataLast int32 = 1
var maxFDdataLen int = 102400000

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
		fmt.Println("收到接入客户")
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
	tempResData := make(map[string]map[int][]byte)
	tempResHeader := make(map[string]map[string]string)
	tempResIndex := make(map[string]map[int]bool)
	OutputChan[index] = c
	
	allBuff := make([]byte, 0)
	headerDemo := make(map[string]string)
	var dataLen int32
	a := 1
	buff := make([]byte, 24)
	useData := public.UseType{Index: index, Con: "TCP"}
	//a := 1
	defer func() {
		if err:=recover();err!=nil{
            fmt.Println(err) // 这里的err其实就是panic传入的内容，55
        }
		delete(OutputChan, index)
		close(c)
		conn.Close()
	}()
	go ResponseData(conn, c)
	for {
		buff = make([]byte, 24)
		i, err := conn.Read(buff)
		allBuff = make([]byte, 0)

		if err != nil {
			fmt.Println(err)
			return
		}
		if i != 24 {
			fmt.Println(i)
			continue
		}
		//continue
		a++
		if a%1000 == 0 {
			fmt.Println(a)
		}
		APICode := public.BytesToInt32(buff[0:4])

		localIndex := public.BytesToInt32(buff[4:8])
		dataLen = public.BytesToInt32(buff[8:12])
		dataType = public.BytesToInt32(buff[12:16])
		dataIndex = public.BytesToInt32(buff[16:20])
		dataLast = public.BytesToInt32(buff[20:24])

		//return
		if dataLen > 0 {
			if dataLen>int32(2048*1024){
				dataLen=int32(2048*1024)
			}
			relData := make([]byte, int(dataLen))
			n, err := conn.Read(relData)
			fmt.Println(relData)
			if err != nil {
				fmt.Println("conn closed")
				return
			}

			fmt.Println("已接收到数据")

			if int32(n) != dataLen {
				continue
			}

			fmt.Println("接收到数据长度：" + strconv.Itoa(n))
			//fmt.Println(decodeData)
			// allBuff = append(allBuff, public.Int32ToBytes(resCode)...)
			// allBuff = append(allBuff, buff[4:8]...)
			// allBuff = append(allBuff, public.Int32ToBytes(zeroLen)...)
			// allBuff = append(allBuff, buff[12:24]...)
			// //fmt.Println(allBuff)
			// conn.Write(allBuff)
			// var canReadData interface{}
			// msgpack.Unmarshal(relData, &canReadData)
			// useDataData := make(map[string][1]interface{})
			// if data, ok := canReadData.(map[string]interface{}); ok {
			// 	for k, v := range data {
			// 		tempData := [1]interface{}{}
			// 		tempData[0] = v
			// 		useDataData[k] = tempData
			// 	}
			// }

			// useData.Data = useDataData

			//headerDemo["length"] = strconv.Itoa(int(dataLen))
			headerDemo["APICode"] = strconv.Itoa(int(APICode))
			headerDemo["localIndex"] = strconv.Itoa(int(localIndex))
			dataIndexInt := int(dataIndex)
			headerDemo["dataLast"] = strconv.Itoa(int(dataLast))
			headerDemo["encode"] = strconv.Itoa(int(dataType))

			//thisHeader, _ := tempResHeader[headerDemo["localIndex"]]["dataIndex"].(map[int]bool)
			if _, ok := tempResData[headerDemo["localIndex"]]; ok {
				thisData := tempResData[headerDemo["localIndex"]]
				if _, ok := thisData[dataIndexInt]; ok == false {
					thisData[dataIndexInt] = relData
					tempResIndex[headerDemo["localIndex"]][dataIndexInt] = true
				}
			} else {
				//num, _ := strconv.Atoi(headerDemo["dataLast"])
				dataIndexList := make(map[int]bool)
				//headerDemo["dataIndex"] = make(map[int]bool, 5)
				tempResHeader[headerDemo["localIndex"]] = headerDemo
				tempResIndex[headerDemo["localIndex"]] = dataIndexList
				tempResData[headerDemo["localIndex"]] = make(map[int][]byte)
				tempResData[headerDemo["localIndex"]][dataIndexInt] = relData
				tempResIndex[headerDemo["localIndex"]][dataIndexInt] = true
			}

			//对分包进行检测，小于当前的序列如果没有收到，则通知对方重发，如果接受完毕，数据押入通道 如果超时，则直接返回失败
			revFlag := true
			for a := 0; a < dataIndexInt; a++ {
				if tempResIndex[headerDemo["localIndex"]][a] == false {
					revFlag = false
					allBuff = make([]byte, 0)
					allBuff = append(allBuff, public.Int32ToBytes(restCode)...)
					allBuff = append(allBuff, buff[4:8]...)
					allBuff = append(allBuff, public.Int32ToBytes(zeroLen)...)
					allBuff = append(allBuff, buff[12:16]...)
					allBuff = append(allBuff, public.Int32ToBytes(int32(a))...)
					allBuff = append(allBuff, buff[20:24]...)
					//fmt.Println(allBuff)
					conn.Write(allBuff)
				}
			}

			maxDataIndex, _ := strconv.Atoi(tempResHeader[headerDemo["localIndex"]]["dataLast"])
			if revFlag == true && dataIndexInt == maxDataIndex-1 {
				allBuffs := make([]byte, 0)
				for i = 0; i < maxDataIndex; i++ {
					allBuffs = append(allBuffs, tempResData[headerDemo["localIndex"]][i]...)
				}
				//allBuffsLen := len(allBuffs)

				var decodeData interface{}
				switch int(dataType) {
				case 1:

					_ = msgpack.Unmarshal(allBuffs, &decodeData)
					break
				case 2:
					_ = json.Unmarshal(allBuffs, &decodeData)
					break
				default:
					continue
					break
				}
				fmt.Println(decodeData)
				//这里进行单数据转数组的操作
				mapv,ok:=decodeData.(map[string]interface{})
				if !ok{
					continue
				}
				sample:=[1]interface{}{}
				for k,v :=range mapv{
					
					sample[0]=v
					mapv[k]=sample
				}
				fmt.Println(sample)
				fmt.Println(mapv)
				useData.API = headerDemo["APICode"]
				useData.Header = headerDemo
				useData.Data = mapv
				InputChan <- useData
			}

			//conn.Write(allBuff)

		} else {
			switch APICode {
			case restCode:
				restAsk := public.RestAskType{localIndex, dataIndex, dataLast}

				c <- restAsk
				break
			case resCode:
				break
			default:
				break
			}
		}
		// n, err := conn.Read(buf)

		// if err != nil {
		// 	fmt.Println("conn closed")
		// 	return
		// }
		// var decodeData interface{}
		// _ = msgpack.Unmarshal(buf, &decodeData)

		// fmt.Println(decodeData)
		// conn.Write(buf[0:n])

		//InputChan
		//fmt.Println("recv msg:", buf[0:n])
		//conn.Write(buf[0:n])
		//fmt.Println("recv msg:", string(buf[0:n]))
		//fmt.Println(buf[0:n])

	}
}
func loadMsg() {}

func ResponseData(conn net.Conn, c chan interface{}) {
	var respDataDemo []byte
	tempDataS := make(map[int32]map[int][]byte)
	tempHeaderS := make(map[int32][]byte)
	for {
		respDataDemo = make([]byte, 0)
		data, isClose := <-c
		fmt.Println("开始返回数据")
		if !isClose {
			fmt.Println("通道已关闭")
			return
		}
		dataTCP, ok := data.(public.TCPType)
		if ok != true {
			//如果是重发请求，进行处理
			switch v := data.(type) {
			case public.RestAskType:

				sendHeader, ok := tempHeaderS[v.Index]
				dataAllNum := len(tempDataS[v.Index])
				if ok == false {
					continue
				}
				i := int(v.MsgIndex)
				if int(i) > dataAllNum {
					continue
				}
				sliceLen := public.Int32ToBytes(int32(len(tempDataS[v.Index][i])))
				for n := 0; n < 4; n++ {
					sendHeader[8+n] = sliceLen[n]
				}
				sliceIndex := public.Int32ToBytes(int32(i))
				for n := 0; n < 4; n++ {
					sendHeader[16+n] = sliceIndex[n]
				}

				sendHeader = append(sendHeader, tempDataS[v.Index][i]...)

				fmt.Println(sendHeader)
				//restData, _ := msgpack.Marshal(dataTCP)
				conn.Write(sendHeader)
				break
			default:
				break
			}
			continue
		} else {

			//发送数据，拉入缓存并进行顺序发送
			var jsonStr interface{}
			if v, ok := dataTCP.Data.(string); ok {
				err := json.Unmarshal([]byte(v), &jsonStr)
				fmt.Println(jsonStr)
				if err != nil {
					jsonStr = dataTCP.Data
					// fmt.Println("json解析错误")
					// fmt.Println(err)
					// continue
				}
				//这里序列化方式写死了，后期需要改写成自动适配访问的序列化方式
				dataTCPData, _ := json.Marshal(jsonStr)
				tempDataOne := make(map[int][]byte)
				dataAllNum := int(len(dataTCPData) / maxFDdataLen)
				if len(dataTCPData)%maxFDdataLen > 0 {
					dataAllNum++
				}
				for i := 0; i < dataAllNum-1; i++ {
					tempDataOne[i] = dataTCPData[i*maxFDdataLen : (i+1)*maxFDdataLen]
				}
				tempDataOne[dataAllNum-1] = dataTCPData[(dataAllNum-1)*maxFDdataLen:]
				tempDataS[dataTCP.Index] = tempDataOne
				respDataDemo = append(respDataDemo, public.Int32ToBytes(dataTCP.API)...)
				respDataDemo = append(respDataDemo, public.Int32ToBytes(dataTCP.Index)...)
				respDataDemo = append(respDataDemo, public.Int32ToBytes(int32(len(dataTCPData)))...)
				respDataDemo = append(respDataDemo, public.Int32ToBytes(int32(2))...)
				respDataDemo = append(respDataDemo, public.Int32ToBytes(int32(0))...)
				respDataDemo = append(respDataDemo, public.Int32ToBytes(int32(dataAllNum))...)
				tempHeaderS[dataTCP.API] = respDataDemo

				sendHeader := respDataDemo
				for i := 0; i < dataAllNum; i++ {
					sliceLen := public.Int32ToBytes(int32(len(tempDataOne[i])))
					for n := 0; n < 4; n++ {
						sendHeader[8+n] = sliceLen[n]
					}
					sliceIndex := public.Int32ToBytes(int32(i))
					for n := 0; n < 4; n++ {
						sendHeader[16+n] = sliceIndex[n]
					}

					sendHeader = append(sendHeader, tempDataOne[i]...)
				}

				fmt.Println(sendHeader)
				//restData, _ := msgpack.Marshal(dataTCP)
				conn.Write(sendHeader)
			}
		}
		fmt.Println("通道正常")

	}
}
