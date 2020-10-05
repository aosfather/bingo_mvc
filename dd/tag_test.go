package dd

import "testing"

type testStruct struct {
	A string `DD:"Test"`
}
func TestValidate(t *testing.T) {
	t.Log("test")
	domain:=DomainType{"Test","测试",T_string,10,false,nil,""}
	a:=testStruct{}
	a.A=""
	AddDomainType(domain)
	a.A=""
	t.Log(Validate(&a))
	a.A="123"
	t.Log(Validate(&a))
	a.A="1234567890123"
	t.Log(Validate(&a))

	//validate by de
	de:=DataElement{Code:"TestD",DataType: "Test"}
	AddDataElement(de)
	t.Log(ValidateByDataType("A",a.A,"TestD"))
	a.A="1234"
	t.Log(ValidateByDataType("A",a.A,"TestD"))
}
