package sqlgo

import "testing"

const genericQuery string = "SELECT t.column_one AS columnOne, t.column_two AS columnTwo, t.column_three AS columnThree, t.column_no_alias, (SELECT sub_query_key AS sub_query_value FROM sub_table AS t WHERE t.column_one ILIKE ANY ($1) AND t.column_two ILIKE ANY ($2) AND t.column_three=$3) AS sq FROM table AS t WHERE t.column_one ILIKE ANY ($1) AND t.column_two ILIKE ANY ($2) AND t.column_three=$4"

func TestGenericQueryPrependWay(t *testing.T) {
	sql := NewSQLGo().
		SQLSelect(
			SetSQLSelect("t.column_one", "columnOne"),
			SetSQLSelect("t.column_two", "columnTwo"),
			SetSQLSelect("t.column_three", "columnThree"),
			SetSQLSelect("t.column_no_alias", ""),
			SetSQLSelect(
				NewSQLGo().SQLSelect(
					SetSQLSelect("sub_query_key", "sub_query_value"),
				).SQLFrom("sub_table", "t").
					SQLWhere(
						SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3}),
						SetSQLWhere("AND", "t.column_two", "ILIKE ANY", []int{1, 2, 3, 4}),
						SetSQLWhere("AND", "t.column_three", "=", 1234),
						SetSQLWhere("AND", "t.column_three", "=", "abc"),
					), "sq",
			),
		).
		SQLFrom("table", "t").
		SQLWhere(
			SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3}),
			SetSQLWhere("AND", "t.column_two", "ILIKE ANY", []int{1, 2, 3, 4}),
			SetSQLWhere("AND", "t.column_three", "=", 1231),
			SetSQLWhere("AND", "t.column_three", "=", "abc"),
			SetSQLWhere("AND", "t.column_three", "=", "abcd"),
		)
	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
		t.Logf("sqlParam: %s", sql.GetSQLGoParameter().GetSQLParameterList())
	}
}

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
