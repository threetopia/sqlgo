package sqlgo

import (
	"fmt"
)

type SQLGoSelect struct {
	values     []sqlGoSelectValues
	params     []interface{}
	paramCount int
}

type sqlGoSelectValues struct {
	Alias string
	Value interface{}
}

func NewSQLGoSelect() *SQLGoSelect {
	return &SQLGoSelect{}
}

func SetSelect(value interface{}, alias string) sqlGoSelectValues {
	return sqlGoSelectValues{
		Value: value,
		Alias: alias,
	}
}

func (ss *SQLGoSelect) SQLSelect(values ...sqlGoSelectValues) *SQLGoSelect {
	ss.values = values
	return ss
}

func (ss *SQLGoSelect) BuildSQL() string {
	if len(ss.values) < 1 {
		return ""
	}

	sql := "SELECT "
	for i, v := range ss.values {
		if i > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}

		alias := ""
		if v.Alias != "" {
			alias = fmt.Sprintf(" AS %s", v.Alias)
		}

		switch vType := v.Value.(type) {
		case *SQLGo:
			sql = fmt.Sprintf("%s(%s)%s", sql, vType.BuildSQL(), alias)
		default:
			sql = fmt.Sprintf("%s%s%s", sql, vType, alias)
		}
	}

	return sql
}

func (ss *SQLGoSelect) GetParams() []interface{} {
	return ss.params
}

func (ss *SQLGoSelect) GetParamsCount() int {
	return ss.paramCount
}
