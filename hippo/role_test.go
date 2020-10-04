package hippo

import "testing"

func TestRole_AddRow(t *testing.T) {
	r := Role{Code: "test", Label: "角色1"}
	r.Init()
	row := make(AuthRow)
	f := AuthField{Code: "F1", Label: "测试字段", Type: STRING}
	row.Add(&f, "S", "1", false)
	f2 := AuthField{Code: "F2", Label: "测试字段2", Type: STRING}
	row.Add(&f2, "S", "1", false)
	r.AddRow("mytable1", &row, false)
	p := make(map[string]interface{})
	p["F1"] = "12"
	p["F2"] = "1"
	t.Log(r.HasPermition("mytable1", p))
	p["F1"] = "1"
	t.Log(r.HasPermition("mytable1", p))

}
