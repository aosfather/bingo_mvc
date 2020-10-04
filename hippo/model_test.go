package hippo

import "testing"

func TestFType_ToString(t *testing.T) {
	var ft FType = STRING
	t.Log(ft.ToString())
}

func TestAuthColumn_Validate(t *testing.T) {
	f := AuthField{Code: "F1", Label: "测试字段", Type: STRING}
	col := AuthColumn{Field: &f, Type: "S", Value: "1", Negation: false}
	t.Log(col.Validate("2"))
	t.Log(col.Validate("1"))
}

func TestAuthRow(t *testing.T) {
	f := AuthField{Code: "F1", Label: "测试字段", Type: STRING}
	col := AuthColumn{Field: &f, Type: "S", Value: "1", Negation: false}
	f2 := AuthField{Code: "F2", Label: "测试字段2", Type: STRING}
	col2 := AuthColumn{Field: &f2, Type: "S", Value: "1", Negation: false}
	row := make(AuthRow)
	row["F1"] = &col
	row["F2"] = &col2
	p := make(map[string]interface{})
	p["F1"] = "12"
	p["F2"] = "1"
	t.Log("row:", row.HasPermition(p))
	p["F1"] = "1"
	t.Log("row:", row.HasPermition(p))
}

func TestAuthRow2(t *testing.T) {
	row := make(AuthRow)
	f := AuthField{Code: "F1", Label: "测试字段", Type: STRING}
	row.Add(&f, "S", "1", false)
	f2 := AuthField{Code: "F2", Label: "测试字段2", Type: STRING}
	row.Add(&f2, "S", "1", false)
	p := make(map[string]interface{})
	p["F1"] = "12"
	p["F2"] = "1"
	t.Log("row:", row.HasPermition(p))
	p["F1"] = "1"
	t.Log("row:", row.HasPermition(p))
}

func TestAuthCollection_HasPermition(t *testing.T) {
	//正向检查
	f := AuthField{Code: "F1", Label: "测试字段", Type: STRING}
	col := AuthColumn{Field: &f, Type: "S", Value: "1", Negation: false}
	f2 := AuthField{Code: "F2", Label: "测试字段2", Type: STRING}
	col2 := AuthColumn{Field: &f2, Type: "S", Value: "1", Negation: false}
	row := make(AuthRow)
	row["F1"] = &col
	row["F2"] = &col2

	f3 := AuthField{Code: "F3", Label: "测试字段3", Type: STRING}
	col3 := AuthColumn{Field: &f3, Type: "S", Value: "1", Negation: false}
	vrow := make(AuthRow)
	vrow["F3"] = &col3

	p := make(map[string]interface{})
	p["F1"] = "12"
	p["F2"] = "1"
	collection := AuthCollection{Table: "", rows: []*AuthRow{&row}, vetoRows: []*AuthRow{&vrow}}
	t.Log(collection.HasPermition(p))
	p["F1"] = "1"
	t.Log(collection.HasPermition(p))
	p["F3"] = "22"
	t.Log(collection.HasPermition(p))
	p["F3"] = "1"
	t.Log(collection.HasPermition(p))
}
