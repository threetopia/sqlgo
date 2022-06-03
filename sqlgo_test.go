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
		), "qwe"),
		SetSelect(NewSQLGo().SQLSelect(
			SetSelect("asd", ""),
			SetSelect("qwe", "qwe"),
		).SQLFrom("poi", "poi").SQLWhere(
			SetWhere("AND", "poi", "=", 1),
		), "poi"),
	).SQLFrom(NewSQLGo().SQLSelect(
		SetSelect("asd", ""),
		SetSelect("qwe", "qwe"),
	).SQLFrom("asd", "").SQLWhere(
		SetWhere("AND", "asd", "=", "2"),
	), "asd").SQLWhere(
		SetWhere("AND", "asd", "=", "3"),
		SetWhere("AND", "qwe", "=", "4"),
	)

	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}
