package sqlgo

import (
	"fmt"
	"strings"
)

type SQLGoJoin interface {
	SQLJoin(values ...sqlGoJoinValue) SQLGoJoin
	SetSQLJoin(joinType string, table sqlGoTable, alias sqlGoAlias, sqlWhere ...sqlGoWhereValue) SQLGoJoin

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoJoin
	SQLGoMandatory
}

type (
	sqlGoJoin struct {
		values         []sqlGoJoinValue
		sqlGoParameter SQLGoParameter
	}

	sqlGoJoinValue struct {
		joinType string
		table    sqlGoTable
		alias    sqlGoAlias
		sqlWhere []sqlGoWhereValue
	}
)

func NewSQLGoJoin() SQLGoJoin {
	return &sqlGoJoin{}
}

func SetSQLJoinWhere(whereType string, whereColumn string, operator string, value interface{}) sqlGoWhereValue {
	return SetSQLWhereNotParam(whereType, whereColumn, operator, value)
}

func SetSQLJoin(joinType string, table sqlGoTable, alias sqlGoAlias, sqlWhere ...sqlGoWhereValue) sqlGoJoinValue {
	return sqlGoJoinValue{
		joinType: joinType,
		table:    table,
		alias:    alias,
		sqlWhere: sqlWhere,
	}
}

func (s *sqlGoJoin) SQLJoin(values ...sqlGoJoinValue) SQLGoJoin {
	s.values = append(s.values, values...)
	return s
}

func (s *sqlGoJoin) SetSQLJoin(joinType string, table sqlGoTable, alias sqlGoAlias, sqlWhere ...sqlGoWhereValue) SQLGoJoin {
	s.values = append(s.values, SetSQLJoin(joinType, table, alias, sqlWhere...))
	return s
}

func (s *sqlGoJoin) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoJoin {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoJoin) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGoJoin) BuildSQL() string {
	var sql string
	if len(s.values) < 1 {
		return sql
	}

	for i, v := range s.values {
		if i > 0 {
			sql = fmt.Sprintf("%s ", sql)
		}

		sqlWhere := NewSQLGo().SQLWhere(v.sqlWhere...)
		switch vType := v.table.(type) {
		case SQLGo:
			vType.SetSQLGoParameter(s.GetSQLGoParameter())
			s.SetSQLGoParameter(vType.GetSQLGoParameter())
			sqlWhere.SetSQLGoParameter(s.GetSQLGoParameter())
			s.SetSQLGoParameter(sqlWhere.GetSQLGoParameter())
			sql = fmt.Sprintf("%s %s JOIN (%s) AS %s%s",
				sql,
				strings.ToUpper(v.joinType),
				vType.BuildSQL(),
				v.alias,
				sqlWhere.BuildSQL(),
			)
		default:
			sqlWhere.SetSQLGoParameter(s.GetSQLGoParameter())
			s.SetSQLGoParameter(sqlWhere.GetSQLGoParameter())
			sql = fmt.Sprintf("%s %s JOIN %s AS %s%s",
				sql,
				strings.ToUpper(v.joinType),
				vType,
				v.alias,
				strings.ReplaceAll(sqlWhere.BuildSQL(), " WHERE ", " ON "),
			)
		}
	}

	return sql
}
