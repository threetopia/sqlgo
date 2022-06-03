package sqlgo

import "fmt"

type SQLGo struct {
	sqlSelect  *SQLGoSelect
	sqlFrom    *SQLGoFrom
	sqlJoin    *SQLGoJoin
	sqlWhere   *SQLGoWhere
	params     []interface{}
	paramCount int
}

func NewSQLGo() *SQLGo {
	return &SQLGo{
		sqlSelect: NewSQLGoSelect(),
		sqlFrom:   NewSQLGoFrom(),
		sqlJoin:   NewSQLGoJoin(),
		sqlWhere:  NewSQLGoWhere(),
	}
}

func (sg *SQLGo) SQLSelect(values ...SqlGoSelectValue) *SQLGo {
	sg.sqlSelect.SQLSelect(values...)
	return sg
}

func (sg *SQLGo) SetSQLSelect(value interface{}, alias string) *SQLGo {
	sg.sqlSelect.SetSQLSelect(value, alias)
	return sg
}

func (sg *SQLGo) SQLFrom(table interface{}, alias string) *SQLGo {
	sg.sqlFrom.SQLFrom(table, alias)
	return sg
}

func (sg *SQLGo) SQLJoin(values ...SQLGoJoinValue) *SQLGo {
	sg.sqlJoin.SQLJoin(values...)
	return sg
}

func (sg *SQLGo) SetSQLJoin(joinType string, table string, alias string, sqlWhere ...SqlGoWhereValue) *SQLGo {
	sg.sqlJoin.SetSQLJoin(joinType, table, alias, sqlWhere...)
	return sg
}

func (sg *SQLGo) SetSQLFrom(table interface{}, alias string) *SQLGo {
	return sg.SQLFrom(table, alias)
}

func (sg *SQLGo) SQLWhere(values ...SqlGoWhereValue) *SQLGo {
	sg.sqlWhere.SQLWhere(values...)
	return sg
}

func (sg *SQLGo) SetSQLWhere(whereType string, whereColumn string, operator string, value interface{}) *SQLGo {
	sg.sqlWhere.SetSQLWhere(whereType, whereColumn, operator, value)
	return sg
}

func (sg *SQLGo) SetJoinScope() *SQLGo {
	sg.sqlWhere.setJoinScope()
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
	if sqlJoin := sg.sqlJoin.SetParamsCount(sg.GetParamsCount()).BuildSQL(); sqlJoin != "" {
		sql = fmt.Sprintf("%s %s", sql, sqlJoin)
		sg.SetParams(sg.sqlJoin.GetParams()...)
		sg.SetParamsCount(sg.sqlJoin.GetParamsCount())
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
