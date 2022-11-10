package sqlgo

import (
	"fmt"
	"testing"
)

const genericQuery string = "SELECT t.column_one AS columnOne, t.column_two AS columnTwo, t.column_three AS columnThree, t.column_no_alias FROM table AS t INNER JOIN join_table AS jt ON jt.id=t.id WHERE t.column_one ILIKE ANY ($1)"

func TestGenericQueryPrependWay(t *testing.T) {
	sql := NewSQLGo().
		SQLSelect(
			SetSQLSelect("t.column_one", "columnOne"),
			SetSQLSelect("t.column_two", "columnTwo"),
			SetSQLSelect("t.column_three", "columnThree"),
			SetSQLSelect("t.column_no_alias", ""),
		).
		SQLFrom("table", "t").
		SQLJoin(SetSQLJoin("INNER", "join_table", "jt", SetSQLJoinWhere("AND", "jt.id", "=", "t.id"))).
		SQLWhere(
			SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3}),
		)
	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
	}
}

func TestGenericQueryPipeline(t *testing.T) {
	sql := NewSQLGo().
		SetSQLSelect("t.column_one", "columnOne").
		SetSQLSelect("t.column_two", "columnTwo").
		SetSQLSelect("t.column_three", "columnThree").
		SetSQLSelect("t.column_no_alias", "").
		SQLFrom("table", "t").
		SetSQLJoin("INNER", "join_table", "jt", SetSQLJoinWhere("AND", "jt.id", "=", "t.id")).
		SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3})
	if sqlStr := sql.BuildSQL(); sqlStr != genericQuery {
		t.Errorf("reuslt must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
	}
}

const deleteQuery string = "DELETE FROM table WHERE column_one=$2"

func TestDelete(t *testing.T) {
	sql := NewSQLGo().
		SQLDelete("table").
		SetSQLWhere("AND", "column_one", "=", "value_one")
	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
	if sqlStr := sql.BuildSQL(); sqlStr != deleteQuery {
		t.Errorf("reuslt must be (%s) BuildSQL give (%s)", genericQuery, sqlStr)
	}
}

func TestWhereINClause(t *testing.T) {
	sql := NewSQLGo().
		SQLDelete("table").
		SetSQLWhere("AND", "asd", "IN", []string{"satu", "satu", "dua", "tiga", "empat", "satu"})
	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}

func TestWhereAnyClause(t *testing.T) {
	sql := NewSQLGo().
		SQLDelete("table").
		SetSQLWhere("AND", "asd", "ANY", []string{"satu", "satu", "dua", "tiga", "empat", "satu"})
	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}
func TestOffsetLimit(t *testing.T) {
	sql := NewSQLGo().
		SetSQLSelect("t.column_one", "columnOne").
		SetSQLSelect("t.column_two", "columnTwo").
		SetSQLSelect("t.column_three", "columnThree").
		SetSQLSelect("t.column_no_alias", "").
		SQLFrom("table", "t").
		SQLPageLimit(1, 8)

	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}

func TestALL(t *testing.T) {
	sql := NewSQLGo()
	sql.SetSQLSelect("u.id", "id")
	sql.SetSQLSelect("u.full_name", "full_name")
	sql.SetSQLSelect("u.id_card", "id_card")
	sql.SetSQLSelect("u.country_code", "country_code")
	sql.SetSQLSelect("u.search_meta", "search_meta")
	sql.SetSQLSelect("u.data_hash", "data_hash")
	sql.SetSQLSelect("u.deleted", "deleted")
	sql.SetSQLSelect("u.created_at", "created_at")
	sql.SetSQLSelect("u.updated_at", "updated_at")
	sql.SetSQLFrom(`"user"`, "u")
	sql.SetSQLJoin("INNER", "search_meta_view", "smv", SetSQLJoinWhere("AND", "smv.id", "=", "u.id"))

	sql.SetSQLWhere("AND", "u.id", "ANY", []string{"asd"})
	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}
