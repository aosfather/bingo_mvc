package dd
var errors map[int]string=make(map[int]string)
var _default_errors map[int]string=make(map[int]string)
func GetErrorMsg(code int) string {
	if msg,ok:= errors[code];ok{
		return msg
	}
	return GetDefaultErrorMsg(code)
}

func SetErrorMsg(code int,msg string){
	if msg!=""{
		errors[code]=msg
	}
}

func GetDefaultErrorMsg(code int) string{
	return _default_errors[code]
}

func init(){
	//初始化默认错误消息
	_default_errors[0]="success!"
	//基本类型校验
	_default_errors[101]="exceeds the limit"
	_default_errors[102]="can not be null"


}
