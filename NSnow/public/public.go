package public

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//UseType 请求的内部中转数据统一结构
type UseType struct {
	Con    string
	Index  string
	Server string
	API    string
	Token  string
	UID    string
	Method string
	Header map[string]string
	Data   interface{}
}

//TCPType socket数据结构，核心为各个参数
type TCPType struct {
	API   int32
	Index int32
	Len   int32
	Type  int32
	Data  interface{}
}

type RestAskType struct {
	Index    int32
	MsgIndex int32
	MsgLast  int32
}

//GetAnswer 获取响应
func GetAnswer(address string, conn string, url string, data UseType) (string, error) {
	switch conn {
	case "HTTP":
		resData, err := getData("http://"+address+url, data)
		fmt.Println("http://" + address + url)
		if err != nil {
			return "", err
		}
		return resData, nil
		break
	default:
		return "", errors.New("方法暂时不支持")
		break
	}
	return "", errors.New("异常错误，位置serviceCenter getAnswer")
}

//GetRoles 获取接口规范
func GetRoles(address string, conn string, url string) (string, error) {
	switch conn {
	case "HTTP":
		resData, err := getData("http://"+address+url, UseType{Data: map[string]interface{}{"Roles": true}})
		fmt.Println("http://" + address + url)
		if err != nil {
			return "", err
		}
		return resData, nil
		break
	default:
		return "", errors.New("方法暂时不支持")
		break
	}
	return "", errors.New("异常错误，位置serviceCenter getAnswer")
}

func getData(url string, data UseType) (string, error) {
	var jsonData []byte
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	sendDataIo := strings.NewReader(string(jsonData))
	//fmt.Println(string(jsonData))
	//defer sendDataIo
	resData, err := http.Post(url, "application/json; charset=UTF-8", sendDataIo)
	if err != nil {
		return "访问失败", err
	}
	// if resData.StatusCode != 200 {
	// 	return "", errors.New("服务异常，返回状态码为" + string(resData.StatusCode))
	// }

	body, _ := ioutil.ReadAll(resData.Body)
	resData.Body.Close()
	return string(body), err
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func Int32ToBytes(i int32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

func BytesToInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}
