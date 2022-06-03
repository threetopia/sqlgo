package sqlgo

import "fmt"

type SQLGoJoin struct {
	values     []SQLGoJoinValue
	params     []interface{}
	paramCount int
}

type SQLGoJoinValue struct {
	joinType string
	table    interface{}
	alias    string
	sqlWhere []SqlGoWhereValue
}

func NewSQLGoJoin() *SQLGoJoin {
	return &SQLGoJoin{}
}

func SetJoin(joinType string, table string, alias string, sqlWhere ...SqlGoWhereValue) SQLGoJoinValue {
	return SQLGoJoinValue{
		joinType: joinType,
		table:    table,
		alias:    alias,
		sqlWhere: sqlWhere,
	}
}

func (sj *SQLGoJoin) SQLJoin(values ...SQLGoJoinValue) *SQLGoJoin {
	sj.values = append(sj.values, values...)
	return sj
}

func (sj *SQLGoJoin) SetSQLJoin(joinType string, table string, alias string, sqlWhere ...SqlGoWhereValue) *SQLGoJoin {
	sj.values = append(sj.values, SetJoin(joinType, table, alias, sqlWhere...))
	return sj
}

func (sj *SQLGoJoin) BuildSQL() string {
	if len(sj.values) < 1 {
		return ""
	}

	sql := ""
	for i, v := range sj.values {
		if i > 0 {
			sql = fmt.Sprintf("%s ", sql)
		}

		sqlWhere := NewSQLGo().SQLWhere(v.sqlWhere...).SetParamsCount(sj.GetParamsCount()).SetJoinScope()
		switch vType := v.table.(type) {
		case string:
			sql = fmt.Sprintf("%s%s JOIN %s AS %s%s", sql, v.joinType, vType, v.alias, sqlWhere.BuildSQL())
			sj.SetParams(sqlWhere.GetParams()...)
			sj.SetParamsCount(sqlWhere.GetParamsCount())
		}
	}

	return sql
}

func (sj *SQLGoJoin) SetParams(params ...interface{}) *SQLGoJoin {
	if len(params) < 1 {
		return sj
	}
	sj.params = append(sj.params, params...)
	return sj
}

func (sj *SQLGoJoin) GetParams() []interface{} {
	return sj.params
}

func (sj *SQLGoJoin) SetParamsCount(paramsCount int) *SQLGoJoin {
	sj.paramCount = paramsCount
	return sj
}

func (sj *SQLGoJoin) GetParamsCount() int {
	return sj.paramCount
}
