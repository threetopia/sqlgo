package sqlgo

import (
	"fmt"
	"testing"
)

const genericQuery string = "SELECT t.column_one AS columnOne, t.column_two AS columnTwo, t.column_three AS columnThree, t.column_no_alias FROM table AS t"

func TestGenericQueryPrependWay(t *testing.T) {
	sql := NewSQLGo().
		SQLSelect(
			SetSelect("t.column_one", "columnOne"),
			SetSelect("t.column_two", "columnTwo"),
			SetSelect("t.column_three", "columnThree"),
			SetSelect("t.column_no_alias", ""),
		).
		SQLFrom("table", "t")
	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
		t.Errorf("reuslt must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
	}
}

func TestGenericQueryChainWay(t *testing.T) {
	sql := NewSQLGo().
		SetSQLSelect("t.column_one", "columnOne").
		SetSQLSelect("t.column_two", "columnTwo").
		SetSQLSelect("t.column_three", "columnThree").
		SetSQLSelect("t.column_no_alias", "").
		SQLFrom("table", "t")
	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
		t.Errorf("reuslt must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
	}
}
func TestDelet(t *testing.T) {
	sql := NewSQLGo().
		SQLDelete("table").
		SetSQLWhere("AND", "asd", "=", "qwe")
	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}
