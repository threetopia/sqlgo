package sqlgo

import (
	"fmt"
	"testing"
)

func TestOutputRegular(t *testing.T) {
	sql := NewSQLGo().SQLSelect(
		SetSelect("asd", ""),
		SetSelect(NewSQLGo().SQLSelect(
			SetSelect("asd", ""),
			SetSelect("qwe", "qwe"),
		).SQLFrom("poi", "poi").SQLWhere(
			SetWhere("AND", "poi", "=", 1),
			SetWhere("AND", "poi", "=", 2),
			SetWhere("AND", "poi", "ANY", []string{"3", "3"}),
			SetWhere("AND", "poi", "IN", []string{"4", "5"}),
			SetWhere("AND", "poi", "=", 2),
		), "qwe"),
		SetSelect(NewSQLGo().SQLSelect(
			SetSelect("asd", ""),
			SetSelect("qwe", "qwe"),
		).SQLFrom("poi", "poi").SQLWhere(
			SetWhere("AND", "poi", "=", 6),
		), "poi"),
	).SQLFrom(NewSQLGo().SQLSelect(
		SetSelect("asd", ""),
		SetSelect("qwe", "qwe"),
	).SQLFrom("asd", "").SQLWhere(
		SetWhere("AND", "asd", "=", "7"),
	), "asd").
		SQLJoin(SetJoin("INNER", "table2", "tb2", SetJoinWhere("AND", "test1", "=", "test2"), SetJoinWhere("AND", "test1", "=", "test2"))).
		SQLJoin(SetJoin("INNER", "table2", "tb2", SetJoinWhere("AND", "test1", "=", "test2"), SetJoinWhere("AND", "test1", "=", "test2"))).
		SQLWhere(
			SetWhere("AND", "asd", "=", "8"),
			SetWhere("AND", "qwe", "=", "9"),
		).SQLGroup("col1", "col2").SQLOffsetLimit(0, 10)
	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}

func TestOutputChain(t *testing.T) {
	sql := NewSQLGo().
		SetSQLSelect("asd", "asd").
		SetSQLSelect("qwe", "qew").
		SetSQLSelect(NewSQLGo().
			SetSQLSelect("asd", "asd").
			SetSQLSelect("qwe", "qew").
			SetSQLFrom("vbn", "vbn").
			SetSQLWhere("AND", "column1", "=", "123").
			SetSQLWhere("AND", "column1", "=", "123"), "poi").
		SetSQLFrom("table", "").
		SetSQLJoin("INNER", "table", "tbl", SetJoinWhere("ON", "asd", "=", "asdasd")).
		SetSQLJoin("LEFT", "table", "tbl", SetJoinWhere("ON", "asd", "=", "asdasd")).
		SetSQLJoin("OUTER",
			NewSQLGo().
				SetSQLSelect("asd", "asd").
				SetSQLFrom("asdTbl", "").
				SetSQLJoin("OUTER", "tbl4", "tb4", SetJoinWhere("ON", "col1", "=", "kljlj"), SetJoinWhere("AND", "col1", "=", "kljlj")).
				SetSQLJoin("OUTER", "tbl5", "tb5", SetJoinWhere("ON", "col1", "=", "kljlj"), SetJoinWhere("AND", "col1", "=", "kljlj")).
				SetSQLWhere("AND", "ert", "=", "yui").
				SetSQLWhere("AND", "ghj", "=", "jkl").
				SetSQLWhere("AND", "ghj", "IN", []int{1, 2, 3}).
				SetSQLWhere("OR", "ghj", "ANY", []int{1, 2, 3}),
			"tbl", SetWhere("AND", "asd", "=", "asdasd")).
		SetSQLWhere("AND", "column1", "=", "123").
		SetSQLWhere("AND", "column1", "=", "123").
		SetSQLWhere("AND", "column1", "=", "123").
		SetSQLGroup("col1", "col2").
		SetSQLOffset(0).
		SetSQLLimit(100)
	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}

func TestInsert(t *testing.T) {
	sql := NewSQLGo().
		SQLInsert("test", SetInsertColumns("colA", "colB"),
			SetInsertValues("colA1", "colB1"),
			SetInsertValues("colA2", "colB2"),
			SetInsertValues("colA3", "colB3"),
		)
	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}

func TestInsertChain(t *testing.T) {
	sql := NewSQLGo().
		SetSQLInsert("test").
		SetSQLInsertColumn("colA", "colB").
		SetSQLInsertValue("colA1", "colB1").
		SetSQLInsertValue("colA2", "colB2").
		SetSQLInsertValue("colA3", "colB3")
	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}

func TestUpdate(t *testing.T) {
	sql := NewSQLGo().
		SQLUpdate("table",
			SetUpdate("asdasd", "asdsad"),
			SetUpdate("qwerty", "qwerty"),
		).SQLWhere(SetWhere("AND", "asd", "=", 123), SetWhere("AND", "qwe", "=", 456))
	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}

func TestUpdateChain(t *testing.T) {
	sql := NewSQLGo().
		SetSQLUpdate("table").
		SetSQLUpdateValue("asdasd", "asdasd").
		SetSQLUpdateValue("qwerty", "qwerty").
		SetSQLWhere("AND", "asd", "=", 123).
		SetSQLWhere("AND", "qwe", "=", 456)
	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}
