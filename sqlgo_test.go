package sqlgo

import (
	"testing"
)

// const genericQuery string = "SELECT t.column_one AS columnOne, t.column_two AS columnTwo, t.column_three AS columnThree, t.column_no_alias FROM table AS t INNER JOIN join_table1 AS jt1 ON jt1.id=t.id INNER JOIN join_table2 AS jt2 ON jt2.id=t.id WHERE t.column_one ILIKE ANY ($1)"

// func TestGenericQueryPrependWay(t *testing.T) {
// 	sql := NewSQLGo().
// 		SQLSelect(
// 			SetSQLSelect("t.column_one", "columnOne"),
// 			SetSQLSelect("t.column_two", "columnTwo"),
// 			SetSQLSelect("t.column_three", "columnThree"),
// 			SetSQLSelect("t.column_no_alias", ""),
// 		).
// 		SQLFrom("table", "t").
// 		SQLJoin(
// 			SetSQLJoin("INNER", "join_table1", "jt1", SetSQLJoinWhere("AND", "jt1.id", "=", "t.id")),
// 			SetSQLJoin("INNER", "join_table2", "jt2", SetSQLJoinWhere("AND", "jt2.id", "=", "t.id")),
// 		).
// 		SQLWhere(
// 			SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3}),
// 		)
// 	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
// 		t.Errorf("result must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
// 	}
// }

// func TestGenericQueryPipeline(t *testing.T) {
// 	sql := NewSQLGo().
// 		SetSQLSelect("t.column_one", "columnOne").
// 		SetSQLSelect("t.column_two", "columnTwo").
// 		SetSQLSelect("t.column_three", "columnThree").
// 		SetSQLSelect("t.column_no_alias", "").
// 		SQLFrom("table", "t").
// 		SetSQLJoin("INNER", "join_table1", "jt1", SetSQLJoinWhere("AND", "jt1.id", "=", "t.id")).
// 		SetSQLJoin("INNER", "join_table2", "jt2", SetSQLJoinWhere("AND", "jt2.id", "=", "t.id")).
// 		SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3})
// 	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
// 		t.Errorf("reuslt must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
// 	}
// }

// const deleteQuery string = "DELETE FROM table WHERE column_one=$1"

// func TestDelete(t *testing.T) {
// 	sql := NewSQLGo().
// 		SQLDelete("table").
// 		SetSQLWhere("AND", "column_one", "=", "value_one")
// 	if sqlStr := sql.BuildSQL(); sqlStr != deleteQuery {
// 		t.Errorf("result must be (%s) BuildSQL give (%s)", deleteQuery, sqlStr)
// 	}
// }

const insertQuery string = "INSERT INTO table (col1, col2, col3) VALUES ($1, $2, $3), ($1, $2, $3), ($1, $2, $3)"

func TestInsert(t *testing.T) {
	// sql := NewSQLGo().
	// 	SQLInsert("table", SetSQLInsertColumn("col1", "col2", "col3"),
	// 		SetSQLInsertValue("val1", "val2", "val3"),
	// 		SetSQLInsertValue("val1", "val2", "val3"),
	// 		SetSQLInsertValue("val1", "val2", "val3"),
	// 	)
	sql := NewSQLGo().
		SetSQLInsert("table").
		SetSQLInsertColumn("col1", "col2", "col3").
		SetSQLInsertValue("val1", "val2", "val3").
		SetSQLInsertValue("val1", "val2", "val3").
		SetSQLInsertValue("val1", "val2", "val3")
	if sqlStr := sql.BuildSQL(); sqlStr != insertQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", insertQuery, sqlStr)
	}
}

const updateQuery string = "UPDATE table SET col1=$1, col2=$2 WHERE col3=$3"

func TestUpdate(t *testing.T) {
	// sql := NewSQLGo().
	// 	SQLUpdate("table",
	// 		SetSQLUpdateValue("col1", "val1"),
	// 		SetSQLUpdateValue("col2", "val2"),
	// 	).
	// 	SQLWhere(SetSQLWhere("AND", "col3", "=", "val3"))
	sql := NewSQLGo().
		SetSQLUpdate("table").
		SetSQLUpdateValue("col1", "val1").
		SetSQLUpdateValue("col2", "val2").
		SetSQLWhere("AND", "col3", "=", "val3")
	if sqlStr := sql.BuildSQL(); sqlStr != updateQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", updateQuery, sqlStr)
	}
}

// func TestWhereINClause(t *testing.T) {
// 	sql := NewSQLGo().
// 		SetSQLSelect("t.col1", "col1").
// 		SetSQLSelect("t.col2", "col2").
// 		SetSQLFrom("table", "t").
// 		SetSQLWhere("AND", "asd", "IN", []string{"satu", "satu", "dua", "tiga", "empat", "satu"})
// 	fmt.Println(sql.BuildSQL(), sql.GetSQLGoParameter().GetSQLParameter())
// }

// func TestWhereAnyClause(t *testing.T) {
// 	sql := NewSQLGo().
// 		SetSQLSelect("t.col1", "col1").
// 		SetSQLSelect("t.col2", "col2").
// 		SetSQLFrom("table", "t").
// 		SetSQLWhere("AND", "asd", "ANY", []string{"satu", "satu", "dua", "tiga", "empat", "satu"})
// 	fmt.Println(sql.BuildSQL(), sql.GetSQLGoParameter().GetSQLParameter())
// }

// func TestOffsetLimit(t *testing.T) {
// 	sql := NewSQLGo().
// 		SetSQLSelect("t.column_one", "columnOne").
// 		SetSQLSelect("t.column_two", "columnTwo").
// 		SetSQLSelect("t.column_three", "columnThree").
// 		SetSQLSelect("t.column_no_alias", "").
// 		SQLFrom("table", "t").
// 		SQLPageLimit(1, 10)
// 	fmt.Println(sql.BuildSQL(), sql.GetSQLGoParameter().GetSQLParameter())
// }

// func TestALL(t *testing.T) {
// 	sql := NewSQLGo()
// 	sql.SetSQLSelect("u.id", "id")
// 	sql.SetSQLSelect("u.full_name", "full_name")
// 	sql.SetSQLSelect("u.id_card", "id_card")
// 	sql.SetSQLSelect("u.country_code", "country_code")
// 	sql.SetSQLSelect("u.search_meta", "search_meta")
// 	sql.SetSQLSelect("u.data_hash", "data_hash")
// 	sql.SetSQLSelect("u.deleted", "deleted")
// 	sql.SetSQLSelect("u.created_at", "created_at")
// 	sql.SetSQLSelect("u.updated_at", "updated_at")
// 	sql.SetSQLFrom(`"user"`, "u")
// 	sql.SetSQLJoin("INNER", "search_meta_view", "smv", SetSQLJoinWhere("AND", "smv.id", "=", "u.id"))
// 	sql.SetSQLWhere("AND", "u.id", "ANY", []string{"asd"})
// 	fmt.Println(sql.BuildSQL(), sql.GetSQLGoParameter().GetSQLParameter())
// }
