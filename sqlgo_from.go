package sqlgo

import "fmt"

type SQLGoFrom interface {
	SQLFrom(table sqlGoTable, alias sqlGoAlias, columns ...sqlGoFromColumn) SQLGoFrom
	SetSQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGoFrom
	SetSQLFromColumn(columns ...sqlGoFromColumn) SQLGoFrom

	SetSQLGoSchema(schema SQLGoSchema) SQLGoFrom
	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoFrom
	SQLGoBase
}

type (
	sqlGoFrom struct {
		table          sqlGoTable
		alias          sqlGoAlias
		columns        sqlGoFromColumnSlice
		sqlGoSchema    SQLGoSchema
		sqlGoParameter SQLGoParameter
	}

	sqlGoFromColumn      string
	sqlGoFromColumnSlice []sqlGoFromColumn
)

func NewSQLGoFrom() SQLGoFrom {
	return new(sqlGoFrom)
}

func (s *sqlGoFrom) SQLFrom(table sqlGoTable, alias sqlGoAlias, columns ...sqlGoFromColumn) SQLGoFrom {
	s.SetSQLFrom(table, alias)
	s.SetSQLFromColumn(columns...)
	return s
}

func (s *sqlGoFrom) SetSQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGoFrom {
	s.table = table
	s.alias = alias
	return s
}

func (s *sqlGoFrom) SetSQLFromColumn(columns ...sqlGoFromColumn) SQLGoFrom {
	s.columns = append(s.columns, columns...)
	return s
}

func (s *sqlGoFrom) SetSQLGoSchema(sqlGoSchema SQLGoSchema) SQLGoFrom {
	s.sqlGoSchema = sqlGoSchema
	return s
}

func (s *sqlGoFrom) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoFrom {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoFrom) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGoFrom) BuildSQL() string {
	var sql string
	if s.table == nil {
		return sql
	}

	sql = "FROM "
	switch vType := s.table.(type) {
	case SQLGo:
		vType.SetSQLGoParameter(s.GetSQLGoParameter())
		sql = fmt.Sprintf("%s(%s)", sql, vType.BuildSQL())
		s.SetSQLGoParameter(vType.GetSQLGoParameter())
	default:
		sql = fmt.Sprintf("%s%s%s", sql, s.sqlGoSchema.BuildSQL(), vType)
	}

	if len(s.columns) > 0 {
		sql = fmt.Sprintf("%s %s(", sql, s.alias)
		for i, v := range s.columns {
			if i > 0 {
				sql = fmt.Sprintf("%s, ", sql)
			}
			sql = fmt.Sprintf("%s%s", sql, v)
		}
		sql = fmt.Sprintf("%s)", sql)
	} else if s.alias != "" {
		sql = fmt.Sprintf("%s AS %s", sql, s.alias)
	}
	return sql
}
