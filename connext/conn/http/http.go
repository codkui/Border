package http

import (
	"NSnow/public"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

//maxConIndex 通信协议内部的唯一id号
var maxConIndex int64 = 1

//ConOutputChan 通讯协议的id对应连接的通道对应表
var ConOutputChan map[string]chan interface{}

//ConInputChan 连接的缓冲队列
var ConInputChan chan public.UseType

//OutType 请求的返回部分统一结构
type OutType struct {
	Code int16
	Data [1]interface{}
	Msg  string
}

//DataType 请求的header部分统一结构
type DataType map[string]string

//LoadType 请求的参数部分统一结构
type LoadType map[string]interface{}

func init() {
	ConOutputChan = make(map[string]chan interface{})
	ConInputChan = make(chan public.UseType, 3000)
}

/*Start 启动服务，传递启动IPV4与端口，不传递则创建随机自由端口，函数会返回成功的创建端口，失败则返回-1
 */
func Start(ip string, port int) (int16, error) {

	if port == 0 {
		port = rand.Intn(55535)
		port += 10000
	}
	http.HandleFunc("/", request)
	fmt.Println("运行在" + ip + ":" + strconv.Itoa(port))
	err := http.ListenAndServe(ip+":"+strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println("建立http服务失败")
		os.Exit(1)
	}
	return 90, nil
}

/*Request 接受请求，返回唯一标识与反序列的数据
 */

func request(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		data, _ := testOutFile()
		res.Write([]byte(data))
		return
	}
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var outStr []byte
	if req.URL.Path == "/response" {
		data, err := loadResponseData(req)
		if err != nil {
			outStr = outPut(false, 400, "数据无法解析")
		} else {
			ConOutputChan[data.Index] <- data
			outStr = outPut(true, 200, "")
		}

		res.Write(outStr)
	}
	fmt.Println("接受到请求" + req.URL.Path)

	loadD, err := loadRequestData(req)
	if err != nil {
		//fmt.Println(err)
		outStr = outPut(false, 400, "数据无法解析")
		res.Write(outStr)
	} else {
		index := strconv.FormatInt(maxConIndex, 10)
		maxConIndex++
		c := make(chan interface{}, 1)
		ConOutputChan[index] = c
		//outStr = outPut(loadD, 200, "")
		//通用数据结构体构建
		headerData, _ := loadRequesetHeader(req)
		token := ""
		token = headerData["Token"]
		if _, ok := loadD["Token"]; ok {
			if v, ok := loadD["Token"].(string); ok {
				token = v
			}
			//token = loadD["Token"]
		}
		useData := public.UseType{
			API:    req.URL.Path,
			Header: headerData,
			Data:   loadD,
			Method: req.Method,
			Token:  token,
			Index:  index,
			Con:    "HTTP",
		}
		fmt.Println("压入通道")
		ConInputChan <- useData
		fmt.Println("压入通道")
		newUserData := <-c
		delete(ConOutputChan, index)
		close(c)
		outStr = outPut(newUserData, 200, "")
		res.Write(outStr)
	}
	res.Write([]byte(""))
}

func loadResponseData(req *http.Request) (public.UseType, error) {
	dataDemo := public.UseType{}
	resData, err := loadRequestData(req)
	if err != nil {
		return dataDemo, err
	}
	// switch v := resData["Con"].(type) {
	// case string:
	// 	var s string
	// 	s = v
	// 	dataDemo.Con = s
	// 	break
	// }
	// t := reflect.TypeOf(dataDemo)
	// //v := reflect.ValueOf(dataDemo)

	// for i := 0; i < t.NumField(); i++ {
	// 	key := t.Field(i).Name
	// 	if value, ok := resData[key].(string); ok && key != "Header" {
	// 		dataDemo.Con = value
	// 	}

	//}
	if v, ok := resData["Con"].(string); ok {
		dataDemo.Con = v
	}
	if v, ok := resData["Index"].(string); ok {
		dataDemo.Index = v
	}
	if v, ok := resData["Server"].(string); ok {
		dataDemo.Server = v
	}
	if v, ok := resData["API"].(string); ok {
		dataDemo.API = v
	}
	if v, ok := resData["Token"].(string); ok {
		dataDemo.Token = v
	}
	if v, ok := resData["Header"].(map[string]string); ok {
		dataDemo.Header = v
	}
	if v, ok := resData["Data"].(map[string]interface{}); ok {
		dataDemo.Data = v
	}
	if v, ok := resData["Method"].(string); ok {
		dataDemo.Method = v
	}

	return dataDemo, nil
}

//解析header头
func loadRequesetHeader(req *http.Request) (DataType, error) {
	demoData := make(DataType)
	for k, v := range req.Header {
		demoData[k] = v[0]
	}
	return demoData, nil
}

/*loadRequestData 解析http所有类型的参数
 */
func loadRequestData(req *http.Request) (map[string]interface{}, error) {
	//fmt.Println("解析次数")
	data := make(LoadType)
	switch req.Method {
	case "GET":
		reqData := req.URL.Query()
		for k, v := range reqData {
			data[k] = v
		}
		//fmt.Println(reqData)
		break
	case "POST":
		//queryForm, _ := url.ParseQuery(req.URL.RawQuery)

		switch req.Header["Content-Type"][0] {
		case "application/json; charset=UTF-8":
			jsonData, _ := ioutil.ReadAll(req.Body)
			var reqData map[string]interface{}
			json.Unmarshal(jsonData, &reqData)
			for k, v := range reqData {
				data[k] = v
			}
			break
		default:
			req.ParseForm()
			for k, v := range req.Form {
				if len(v) == 1 {
					data[k] = v[0]
				} else {
					data[k] = v
				}
			}
			break
		}
		break
	case "PUT":
		break
	case "DELETE":
		break
	case "HEAD":
		break
	default:
		break
	}

	return data, nil
}

/*Response 返回给请求数据，传递唯一标识与序列化的数据
 */
func Response() (bool, error) {
	return true, nil
}

func testOutFile() (string, error) {
	data, err := ioutil.ReadFile("test.html")
	//fmt.Println(data)
	if err != nil {
		fmt.Println(err)
	}
	return string(data), err
}

func outPut(data interface{}, code int16, msg string) []byte {
	if v, ok := data.(string); ok {
		var s string
		s = v
		return []byte(s)
	}
	if code == 0 {
		code = 200
	}
	outData := new(OutType)
	outData.Code = code
	outData.Data[0] = data
	outData.Msg = msg

	outStr, err := json.Marshal(outData)
	if err != nil {
		fmt.Println(err)
	}
	return outStr
}
