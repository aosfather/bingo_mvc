package sqltemplate

import "testing"

type test1 struct {
	My string `Table:"my_test" Option:"auto"`
	Test2
	My0 string
}

type Test2 struct {
	My2 string
	My3 string `Option:"pk"`
}

func TestSqlTemplate_CreateInserSql(t *testing.T) {
	sql := SqlTemplate{}
	t1 := test1{}
	t1.My = "12"
	t1.My3 = "34"
	t1.My0 = "end"
	str, v, _ := sql.CreateInserSql(&t1)
	t.Log(str, v)
	t.Log(v)

	str, v, _ = sql.CreateUpdateSql(&t1, "My0")
	t.Log(str, v)
	t.Log(v)
}

func TestSqlTemplate_CreateQuerySql(t *testing.T) {
	sql := SqlTemplate{}
	t1 := test1{}
	t1.My = "12"
	t1.My3 = "34"
	t1.My0 = "end"
	str, v, _ := sql.CreateQuerySql(&t1)
	t.Log(str, v)
	t.Log(v)
}

func TestSqlTemplate_CreateDeleteSql(t *testing.T) {
	sql := SqlTemplate{}
	t1 := test1{}
	t1.My = "12"
	t1.My3 = "34"
	t1.My0 = "end"
	str, v, _ := sql.CreateDeleteSql(&t1)
	t.Log(str, v)
	t.Log(v)
}
