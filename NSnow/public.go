package public

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
