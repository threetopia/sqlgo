package sqlgo

import "fmt"

type SQLGoFrom struct {
	table      interface{}
	alias      string
	params     []interface{}
	paramCount int
}

func NewSQLGoFrom() *SQLGoFrom {
	return &SQLGoFrom{}
}

func (sg *SQLGoFrom) SQLFrom(table interface{}, alias string) *SQLGoFrom {
	sg.table = table
	sg.alias = alias
	return sg
}

func (sf *SQLGoFrom) BuildSQL() string {
	if sf.table == nil {
		return ""
	}

	sql := "FROM "
	switch vType := sf.table.(type) {
	case *SQLGo:
		sql = fmt.Sprintf("%s(%s)", sql, vType.SetParamsCount(sf.GetParamsCount()).BuildSQL())
		sf.SetParams(vType.GetParams()...).
			SetParamsCount(vType.GetParamsCount())
	default:
		sql = fmt.Sprintf("%s%s", sql, vType)
	}

	if sf.alias != "" {
		sql = fmt.Sprintf("%s AS %s", sql, sf.alias)
	}
	return sql
}

func (sf *SQLGoFrom) SetParams(params ...interface{}) *SQLGoFrom {
	if len(params) < 1 {
		return sf
	}
	sf.params = append(sf.params, params...)
	return sf
}

func (sf *SQLGoFrom) GetParams() []interface{} {
	return sf.params
}

func (sf *SQLGoFrom) SetParamsCount(paramsCount int) *SQLGoFrom {
	sf.paramCount = paramsCount
	return sf
}

func (sf *SQLGoFrom) GetParamsCount() int {
	return sf.paramCount
}
