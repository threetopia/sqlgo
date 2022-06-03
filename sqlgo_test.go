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
		SQLJoin(SetJoin("INNER", "table2", "tb2", SetWhere("AND", "test1", "=", "test2"), SetWhere("AND", "test1", "=", "test2"))).
		SQLJoin(SetJoin("INNER", "table2", "tb2", SetWhere("AND", "test1", "=", "test2"), SetWhere("AND", "test1", "=", "test2"))).
		SQLWhere(
			SetWhere("AND", "asd", "=", "8"),
			SetWhere("AND", "qwe", "=", "9"),
		)

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
		SetSQLJoin("INNER", "table", "tbl", SetWhere("AND", "asd", "=", "asdasd")).
		SetSQLJoin("LEFT", "table", "tbl", SetWhere("AND", "asd", "=", "asdasd")).
		SetSQLWhere("AND", "column1", "=", "123").
		SetSQLWhere("AND", "column1", "=", "123")
	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}
