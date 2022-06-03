package sqlgo

import (
	"fmt"
	"testing"
)

func TestOutput(t *testing.T) {
	sql := NewSQLGo().SQLSelect(
		SetSelect("asd", ""),
		SetSelect("qwe", "qwe"),
		SetSelect(NewSQLGo().SQLSelect(
			SetSelect("asd", ""),
			SetSelect("qwe", "qwe"),
		).SQLFrom("poi", "poi").SQLWhere(
			SetWhere("AND", "poi", "=", "123"),
		), "poi"),
	).SQLFrom(NewSQLGo().SQLSelect(
		SetSelect("asd", ""),
		SetSelect("qwe", "qwe"),
	).SQLFrom("asd", "").SQLWhere(
		SetWhere("AND", "asd", "=", "asd3"),
	), "asd").SQLWhere(
		SetWhere("AND", "asd", "=", "asd"),
		SetWhere("AND", "qwe", "=", "qwe"),
	)

	fmt.Println(sql.BuildSQL(), sql.GetParams(), sql.GetParamsCount())
}
