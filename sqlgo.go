package sqlgo

import "fmt"

type SQLGo struct {
	sqlSelect  *SQLGoSelect
	sqlFrom    *SQLGoFrom
	params     []interface{}
	paramCount int
}

func NewSQLGo() *SQLGo {
	return &SQLGo{
		sqlSelect: NewSQLGoSelect(),
		sqlFrom:   NewSQLGoFrom(),
	}
}

func (sg SQLGo) SQLSelect(values ...sqlGoSelectValues) SQLGo {
	sg.sqlSelect.SQLSelect(values...)
	return sg
}

func (sg SQLGo) SQLFrom(table interface{}, alias string) SQLGo {
	sg.sqlFrom.SQLFrom(table, alias)
	return sg
}

func (sg SQLGo) BuildSQL() string {
	sql := ""
	if sqlSelect := sg.sqlSelect.SetParams(sg.GetParams()).SetParamsCount(sg.GetParamsCount()).BuildSQL(); sqlSelect != "" {
		sql = fmt.Sprintf("%s%s", sql, sqlSelect)
	}
	if sqlFrom := sg.sqlFrom.SetParams(sg.GetParams()).SetParamsCount(sg.GetParamsCount()).BuildSQL(); sqlFrom != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlFrom)
	}
	return sql
}

func (sg SQLGo) SetParams(params []interface{}) SQLGo {
	sg.params = params
	return sg
}

func (sg SQLGo) GetParams() []interface{} {
	return sg.params
}

func (sg SQLGo) SetParamsCount(paramsCount int) SQLGo {
	sg.paramCount = paramsCount
	return sg
}

func (sg SQLGo) GetParamsCount() int {
	return sg.paramCount
}
