package sqlgo

import "testing"

const genericQuery string = "SELECT t.column_one AS columnOne, t.column_two AS columnTwo, t.column_three AS columnThree, t.column_no_alias"

func TestGenericQueryPrependWay(t *testing.T) {
	sql := NewSQLGo().
		SQLSelect(
			SetSQLSelect("t.column_one", "columnOne"),
			SetSQLSelect("t.column_two", "columnTwo"),
			SetSQLSelect("t.column_three", "columnThree"),
			SetSQLSelect("t.column_no_alias", ""),
			SetSQLSelect("sq", NewSQLGo().SQLFrom("sub_table", "sq").SetSQLSelect("sub_query_key", "sub_query_value")),
		).SQLFrom("table", "t")
	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
		t.Logf("sqlStr: %s", sqlStr)
		t.Errorf("result must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
	}
}
