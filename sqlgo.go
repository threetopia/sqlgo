package sqlgo

import "fmt"

type SQLGo struct {
	sqlSelect  *SQLGoSelect
	sqlFrom    *SQLGoFrom
	sqlWhere   *SQLGoWhere
	params     []interface{}
	paramCount int
}

func NewSQLGo() *SQLGo {
	return &SQLGo{
		sqlSelect: NewSQLGoSelect(),
		sqlFrom:   NewSQLGoFrom(),
		sqlWhere:  NewSQLGoWhere(),
	}
}

func (sg *SQLGo) SQLSelect(values ...sqlGoSelectValues) *SQLGo {
	sg.sqlSelect.SQLSelect(values...)
	return sg
}

func (sg *SQLGo) SQLFrom(table interface{}, alias string) *SQLGo {
	sg.sqlFrom.SQLFrom(table, alias)
	return sg
}

func (sg *SQLGo) SQLWhere(values ...sqlGoWhereValues) *SQLGo {
	sg.sqlWhere.SQLWhere(values...)
	return sg
}

func (sg *SQLGo) BuildSQL() string {
	sql := ""
	if sqlSelect := sg.sqlSelect.SetParamsCount(sg.GetParamsCount()).BuildSQL(); sqlSelect != "" {
		sql = fmt.Sprintf("%s%s", sql, sqlSelect)
		sg.SetParams(sg.sqlSelect.GetParams()...)
		sg.SetParamsCount(sg.sqlSelect.GetParamsCount())
	}
	if sqlFrom := sg.sqlFrom.SetParamsCount(sg.GetParamsCount()).BuildSQL(); sqlFrom != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlFrom)
		sg.SetParams(sg.sqlFrom.GetParams()...)
		sg.SetParamsCount(sg.sqlFrom.GetParamsCount())
	}
	if sqlWhere := sg.sqlWhere.SetParamsCount(sg.GetParamsCount()).BuildSQL(); sqlWhere != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlWhere)
		sg.SetParams(sg.sqlWhere.GetParams()...)
		sg.SetParamsCount(sg.sqlWhere.GetParamsCount())
	}
	return sql
}

func (sg *SQLGo) SetParams(params ...interface{}) *SQLGo {
	if len(params) < 1 {
		return sg
	}
	sg.params = append(sg.params, params...)
	return sg
}

func (sg *SQLGo) GetParams() []interface{} {
	return sg.params
}

func (sg *SQLGo) SetParamsCount(paramsCount int) *SQLGo {
	sg.paramCount = paramsCount
	return sg
}

func (sg *SQLGo) GetParamsCount() int {
	return sg.paramCount
}
