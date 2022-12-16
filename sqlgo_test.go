package sqlgo

import (
	"testing"
)

const genericQuery string = "SELECT t.column_one AS columnOne, t.column_two AS columnTwo, t.column_three AS columnThree, t.column_no_alias, (SELECT sub_query_key AS sub_query_value FROM sub_table AS t WHERE t.column_one ILIKE ANY ($1) AND t.column_two ILIKE ANY ($2) AND t.column_three=$3) AS sq FROM table AS t WHERE t.column_one ILIKE ANY ($1) AND t.column_two ILIKE ANY ($2) AND t.column_three=$4"

// func TestGenericQueryPrependWay(t *testing.T) {
// 	sql := NewSQLGo().
// 		SQLSelect(
// 			SetSQLSelect("t.column_one", "columnOne"),
// 			SetSQLSelect("t.column_two", "columnTwo"),
// 			SetSQLSelect("t.column_three", "columnThree"),
// 			SetSQLSelect("t.column_no_alias", ""),
// 			SetSQLSelect(
// 				NewSQLGo().SQLSelect(
// 					SetSQLSelect("sub_query_key", "sub_query_value"),
// 				).SQLFrom("sub_table", "t").
// 					SQLWhere(
// 						SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1}),
// 						SetSQLWhere("AND", "t.column_two", "ILIKE ANY", []int{1, 2}),
// 						SetSQLWhere("AND", "t.column_three", "=", 3),
// 						SetSQLWhere("AND", "t.column_four", "=", "empat"),
// 					), "sq",
// 			),
// 		).
// 		SQLFrom("table", "t").
// 		SQLJoin(SetSQLJoin("INNER", "join_table", "jt", SetSQLJoinWhere("AND", "jt.id", "=", "t.id"), SetSQLWhere("AND", "jt.id", "ILIKE ANY", []string{"coba_satu", "coba_dua"}))).
// 		SQLWhere(
// 			SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1}),
// 			SetSQLWhere("AND", "t.column_two", "ILIKE ANY", []int{1, 2}),
// 			SetSQLWhere("AND", "t.column_five", "=", "lima"),
// 			SetSQLWhere("AND", "t.column_six", "=", 6),
// 			SetSQLWhere("AND", "t.column_seven", "=", 1234567),
// 			SetSQLWhere("AND", "t.column_four", "=", "empat"),
// 		)
// 	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
// 		t.Errorf("result must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
// 		t.Logf("sqlParam: %s", sql.GetSQLGoParameter().GetSQLParameterList())
// 	}
// }

// func TestGenericQueryPipeline(t *testing.T) {
// 	sql := NewSQLGo().
// 		SetSQLSelect("t.column_one", "columnOne").
// 		SetSQLSelect("t.column_two", "columnTwo").
// 		SetSQLSelect("t.column_three", "columnThree").
// 		SetSQLSelect("t.column_no_alias", "").
// 		SetSQLSelect(
// 			NewSQLGo().
// 				SetSQLSelect("sub_query_key", "sub_query_value").
// 				SQLFrom("sub_table", "t").
// 				SQLWhere(
// 					SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3}),
// 					SetSQLWhere("AND", "t.column_two", "ILIKE ANY", []int{1, 2, 3, 4}),
// 					SetSQLWhere("AND", "t.column_three", "=", 1234),
// 				), "sq",
// 		).
// 		SQLFrom("table", "t").
// 		// SetSQLJoin("INNER", "join_table", "jt", SetSQLJoinWhere("AND", "jt.id", "=", "t.id")).
// 		SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3})
// 	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
// 		t.Logf("sqlStr: %s", sqlStr)
// 		t.Errorf("reuslt must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
// 	}
// }

// func TestInsert(t *testing.T) {
// 	sql := NewSQLGo().
// 		SQLInsert("table",
// 			SetSQLInsertColumn("col1", "col2", "col3"),
// 			SetSQLInsertValue("val1-1", "val1-2", "val1-3"),
// 			SetSQLInsertValue("val2-1", "val2-2", "val2-3"),
// 		)
// 	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
// 		t.Errorf("result must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
// 		t.Logf("sqlParam: %s", sql.GetSQLGoParameter().GetSQLParameterList())
// 	}
// }

// func TestUpdate(t *testing.T) {
// 	sql := NewSQLGo().
// 		SQLUpdate("table",
// 			SetSQLUpdate("col1", "satu"),
// 			SetSQLUpdate("col2", "dua"),
// 		)
// 	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
// 		t.Errorf("result must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
// 		t.Logf("sqlParam: %s", sql.GetSQLGoParameter().GetSQLParameterList())
// 	}
// }

// func TestDelete(t *testing.T) {
// 	sql := NewSQLGo().
// 		SQLDelete("table").SQLWhere(SetSQLWhere("AND", "col1", "=", "val1"))

// 	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
// 		t.Errorf("result must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
// 		t.Logf("sqlParam: %s", sql.GetSQLGoParameter().GetSQLParameterList())
// 	}
// }

func TestValues(t *testing.T) {
	sql := NewSQLGo().
		SQLSelect(
			SetSQLSelect("col1", "col1"),
			SetSQLSelect("col2", "col2"),
			SetSQLSelect("col3", "col3"),
		).
		SQLFrom(NewSQLGo().SQLValues(
			SetSQLValuesValue("val1-1", "val1-2", "val1-3"),
			SetSQLValuesValue("val2-1", "val2-2", "val2-3"),
			SetSQLValuesValue("val1-1", "val1-2", "val1-3"),
		), "test", "col1", "col2", "col3")

	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
		t.Logf("sqlParam: %s", sql.GetSQLGoParameter().GetSQLParameterList())
	}
}
