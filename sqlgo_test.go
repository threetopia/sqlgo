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
		), "poi"),
	).SQLFrom(NewSQLGo().SQLSelect(
		SetSelect("asd", ""),
		SetSelect("qwe", "qwe"),
	), "asd")

	fmt.Println(sql.BuildSQL(), sql.GetParams())
}
