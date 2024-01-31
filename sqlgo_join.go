package sqlgo

import (
	"fmt"
	"strings"
)

type SQLGoJoin interface {
	SQLJoin(values ...sqlGoJoinValue) SQLGoJoin
	SetSQLJoin(joinType string, table sqlGoTable, alias sqlGoAlias, sqlWhere ...sqlGoWhereValue) SQLGoJoin

	SetSQLGoSchema(schema SQLGoSchema) SQLGoJoin
	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoJoin
	SQLGoBase
}

type (
	sqlGoJoin struct {
		valueSlice     sqlGoJoinValueSlice
		sqlGoSchema    SQLGoSchema
		sqlGoParameter SQLGoParameter
	}

	sqlGoJoinValue struct {
		joinType string
		table    sqlGoTable
		alias    sqlGoAlias
		sqlWhere sqlGoWhereValueSlice
	}

	sqlGoJoinValueSlice []sqlGoJoinValue
)

func NewSQLGoJoin() SQLGoJoin {
	return new(sqlGoJoin)
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
	s.valueSlice = append(s.valueSlice, values...)
	return s
}

func (s *sqlGoJoin) SetSQLJoin(joinType string, table sqlGoTable, alias sqlGoAlias, sqlWhere ...sqlGoWhereValue) SQLGoJoin {
	s.valueSlice = append(s.valueSlice, SetSQLJoin(joinType, table, alias, sqlWhere...))
	return s
}

func (s *sqlGoJoin) SetSQLGoSchema(sqlGoSchema SQLGoSchema) SQLGoJoin {
	s.sqlGoSchema = sqlGoSchema
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
	if len(s.valueSlice) < 1 {
		return sql
	}

	for i, v := range s.valueSlice {
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
			sql = fmt.Sprintf("%s%s JOIN (%s) AS %s %s",
				sql,
				strings.ToUpper(v.joinType),
				vType.BuildSQL(),
				v.alias,
				sqlWhere.BuildSQL(),
			)
		default:
			sqlWhere.SetSQLGoParameter(s.GetSQLGoParameter())
			s.SetSQLGoParameter(sqlWhere.GetSQLGoParameter())
			sql = fmt.Sprintf("%s%s JOIN %s%s AS %s %s",
				sql,
				strings.ToUpper(v.joinType),
				s.sqlGoSchema.BuildSQL(),
				vType,
				v.alias,
				strings.ReplaceAll(sqlWhere.BuildSQL(), "WHERE", "ON"),
			)
		}
	}

	return sql
}
