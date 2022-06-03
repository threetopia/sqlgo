package sqlgo

import (
	"fmt"
	"testing"
)

func TestOutput(t *testing.T) {
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
	), "asd").SQLWhere(
		SetWhere("AND", "asd", "=", "8"),
		SetWhere("AND", "qwe", "=", "9"),
	)

	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}
