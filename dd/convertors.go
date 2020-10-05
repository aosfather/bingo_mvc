package dd
//默认的无转换器效果
var _noConvertor=&noConvertor{}
type noConvertor struct {

}
func(this *noConvertor) Input(v interface{})interface{}{
	return v
}
func(this *noConvertor) Output(v interface{})interface{}{
	return v
}
//--------------字典转换器-----------------//
func dictConvertorFactory(name string)Convertor{
	catalog:=getDictCatalog(name)
	if catalog!=nil {
		return &dictConvertor{catalog}
	}
	return _noConvertor
}
type dictConvertor struct {
	dictcatalog *DictCatalog
}
func(this *dictConvertor) Input(v interface{})interface{}{
	if this.dictcatalog == nil {
		return v
	}
		for _, item := range this.dictcatalog.Items {
			if item.Code == v || item.Label == v {
				return item.Code
			}
		}
	return nil
}

func(this *dictConvertor) Output(v interface{})interface{}{
	if this.dictcatalog == nil {
		return v
	}
	for _, item := range this.dictcatalog.Items {
		if item.Code == v || item.Label == v {
			return item.Label
		}
	}
	return nil
}

//注册转换器
func init(){
	RegisterConvertorFactory("dict",dictConvertorFactory)
}