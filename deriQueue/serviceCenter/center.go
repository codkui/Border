package serviceCenter

import (
	"NSnow/public"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var maxServiceID int64 = 1

//ServiceLog 服务记录总表
var ServiceLog map[string]ServiceInfo

//ServiceList 服务总表
var ServiceList map[string][]Service

//ServiceHistory 关闭的服务历史记录
var ServiceHistory map[string]ServiceInfo

//Service 服务数据，记录服务名，服务ID，地址，url地址解析正则
type Service struct {
	Name    string //服务名
	ID      string
	Address string
	Conn    string
	URLReg  string

	Status bool
}

//ServiceInfo 记录每个服务的关键信息，连接数，最高连接数，分钟内平均响应时间，分钟内连接记录
type ServiceInfo struct {
	LineNum      int32
	MaxLineNum   int32
	ResponseTime int32
	Log          map[string]interface{}
}

//ResType 解析的访问中转数据结构
type ResType struct {
	Address string
	Conn    string
	URL     string
	data    public.UseType
}

type ApiCodeType struct {
	Service string
	Api     string
}

//AnalyticRes 解析请求数据，取得服务address,允许访问类型,URL,数据体
func AnalyticRes(data public.UseType) (string, string, string, public.UseType, error) {
	//接口序列号访问方式转换

	apiCode, _ := strconv.Atoi(data.API)
	if apiCode > 0 {
		codeList, err := ioutil.ReadFile("apiCode.json")
		if err == nil {
			var codeListData map[string]ApiCodeType
			er := json.Unmarshal(codeList, &codeListData)
			if er != nil {
				return "", "", "", data, errors.New("api格式有误")
			}
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			fmt.Println(codeListData)
			index := r.Intn(len(ServiceList[codeListData[data.API].Service]))
			thisService := ServiceList[codeListData[data.API].Service][index]
			return thisService.Address, thisService.Conn, codeListData[data.API].Api, data, nil
		}
		return "", "", "", data, errors.New("api格式有误")
	}
	//路由拆分
	aguments := strings.Split(data.API, "/")
	if len(aguments) <= 1 {
		return "", "", "", data, errors.New("api格式有误")
	}
	aguments = aguments[1:]
	fmt.Println(aguments)
	service := aguments[0]
	fmt.Println(ServiceList[service])
	//路由解析
	if _, had := ServiceList[service]; had == false || len(ServiceList[service]) == 0 {
		return "", "", "", data, errors.New("服务不存在")
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(ServiceList[service]))
	thisService := ServiceList[service][index]
	//可能以后需要在这里作UID认证

	//根据服务正则解析到服务地址
	reg := thisService.URLReg
	agumNum := strings.Count(reg, "{?")
	if agumNum > len(aguments) {
		return "", "", "", data, errors.New("api参数不足")
	}
	for n := 0; n < len(aguments); n++ {
		reg = strings.Replace(reg, "{?"+strconv.Itoa(n)+"}", aguments[n], -1)
		//fmt.Println(reg)
	}
	return thisService.Address, thisService.Conn, reg, data, nil

}

//DoService 代理请求服务，取得返回数据
func DoService(data public.UseType) (string, error) {
	address, connType, url, data, err := AnalyticRes(data)
	if err != nil {
		return "", err
	}
	fmt.Println(address + url)
	fmt.Println(connType)
	resData, err := public.GetAnswer(address, connType, url, data)
	if err != nil {
		return "", err
	}
	return resData, nil
}

func init() {
	ServiceList = make(map[string][]Service)

	// selfService := Service{Name: "mapping",
	// 	ID:      strconv.FormatInt(maxServiceID, 10),
	// 	Address: "devel.skylinuav.com",
	// 	URLReg:  "/work/Public/{?0}/?service={?2}.{?3}",
	// 	Conn:    "HTTP",
	// }
	apiConf, err := ioutil.ReadFile("apiConf.json")
	//fmt.Println(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	var confData []Service
	err1 := json.Unmarshal(apiConf, &confData)
	if err1 != nil {
		fmt.Println(string(apiConf))
		fmt.Println(err1)
		return
	}
	fmt.Println(confData)
	idStr := ""
	//if value, ok := confData.([]Service); ok {
	for _, v := range confData {
		maxServiceID++
		idStr = strconv.FormatInt(maxServiceID, 10)
		v.ID = idStr
		ServiceList[v.Name] = []Service{v}
	}
	//}

	// ServiceList["mapping"] = []Service{selfService}
	// ServiceList["gga"] = []Service{ggaService}
	// ServiceList["file"] = []Service{fileService}
	//ServiceList["mapping"][]=selfService
	//printSlice(ServiceList["mapping"])
	//append(ServiceList["mapping"], selfService)
	fmt.Println(ServiceList)
}
