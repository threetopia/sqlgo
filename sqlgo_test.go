package sqlgo

import (
	"fmt"
	"testing"
)

func TestOutput(t *testing.T) {
	// sql := NewSQLBuilder().SQLSelect(SQLValue{
	// 	"ascol1": "col1",
	// 	"ascol2": "col2",
	// 	"asSubQ": NewSQLBuilder().SQLSelect(SQLValue{
	// 		"ascol3": "col3",
	// 		"ascol4": "col4",
	// 	}),
	// }).SQLFrom(SQLValue{"asTable": "table"})

	sql := NewSQLBuilder().
		SQLSelect(
			SetSelect("column1", "alias1"),
			SetSelect("column2", "alias2"),
			SetSelect(
				NewSQLBuilder().SQLSelect(
					SetSelect("column3", "alias3"),
				).
					SQLFrom(SetFrom("table3", "alias3")).
					SQLWhere(SetWhere("AND", "test1", "=", "1"), SetWhere("AND", "test2", "=", "2")),
				"alias3"),
			SetSelect("column4", "alias4"),
		).
		SQLFrom(SetFrom("table5", "alias5")).
		SQLJoin(
			SetJoin("LEFT", "joinTable", "jt",
				SetWhere("AND", "jt.id", "=", "alias3.id"),
				SetWhere("AND", "jt.id", "=", "alias2.id"),
			),
			SetJoin("INNER", "joinTable", "jt",
				SetWhere("ON", "jt.id", "=", "alias3.id"),
				SetWhere("AND", "jt.id", "=", "alias2.id"),
			),
		).
		SQLWhere(
			SetWhere("AND", "alias1", "=", "12"),
			SetWhere("AND", "alias1", "ANY", []string{"12", "12", "12"}),
			SetWhere("AND", "alias1", "IN", []string{"12", "12", "12"}),
			SetWhere("AND", "alias1", "IN",
				NewSQLBuilder().SQLSelect(
					SetSelect("asd", "asdAlias"),
				).SQLFrom(SetFrom("testTable", "tt")).
					SQLWhere(
						SetWhere("AND", "tt.id", "=", "valuTable"),
					),
			),
		)

	fmt.Println(sql.BuildSQL(), sql.GetParams())
}
