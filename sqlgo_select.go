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

		switch vType := v.Value.(type) {
		// case string:
		// 	sql = fmt.Sprintf("%s%s", sql, vType)
		case *SQLGo:
			sql = fmt.Sprintf("%s(%s)", sql, vType.SetParamsCount(ss.GetParamsCount()).BuildSQL())
			ss.SetParams(vType.GetParams()...)
			ss.SetParamsCount(vType.GetParamsCount())
		default:
			sql = fmt.Sprintf("%s%s", sql, vType)
		}

		if v.Alias != "" {
			sql = fmt.Sprintf("%s AS %s", sql, v.Alias)
		}
	}

	return sql
}

func (ss *SQLGoSelect) SetParams(params ...interface{}) *SQLGoSelect {
	if len(params) < 1 {
		return ss
	}
	ss.params = append(ss.params, params...)
	return ss
}

func (ss *SQLGoSelect) GetParams() []interface{} {
	return ss.params
}

func (ss *SQLGoSelect) SetParamsCount(paramsCount int) *SQLGoSelect {
	ss.paramCount = paramsCount
	return ss
}

func (ss *SQLGoSelect) GetParamsCount() int {
	return ss.paramCount
}
