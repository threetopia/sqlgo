package sqlgo

import (
	"fmt"
)

type SQLGoInsert interface {
	SQLInsert(table sqlGoTable, columns sqlGoInsertColumnSlice, values ...sqlGoInsertValueSlice) SQLGoInsert
	SetSQLInsert(table sqlGoTable) SQLGoInsert
	SetSQLInsertColumn(columns ...sqlGoInsertColumn) SQLGoInsert
	SetSQLInsertValue(values ...sqlGoInsertValue) SQLGoInsert
	SetSQLInsertValueSlice(values ...sqlGoInsertValueSlice) SQLGoInsert

	SetSQLGoSchema(schema SQLGoSchema) SQLGoInsert
	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoInsert
	SQLGoBase
}

type (
	sqlGoInsert struct {
		table          sqlGoTable
		columns        sqlGoInsertColumnSlice
		values         []sqlGoInsertValueSlice
		sqlGoSchema    SQLGoSchema
		sqlGoParameter SQLGoParameter
	}

	sqlGoInsertColumn      string
	sqlGoInsertColumnSlice []sqlGoInsertColumn
	sqlGoInsertValue       sqlGoValue
	sqlGoInsertValueSlice  []sqlGoInsertValue
	sqlGoInsertToTsVector  struct {
		lang  string
		value sqlGoValue
	}
)

func NewSQLGoInsert() SQLGoInsert {
	return new(sqlGoInsert)
}

func SetSQLInsertColumn(columns ...sqlGoInsertColumn) sqlGoInsertColumnSlice {
	return columns
}

func SetSQLInsertValue(values ...sqlGoInsertValue) sqlGoInsertValueSlice {
	return values
}

func SetSQLInsertToTsVector(lang string, value sqlGoInsertValue) sqlGoInsertToTsVector {
	return sqlGoInsertToTsVector{
		lang:  lang,
		value: value,
	}
}

func (s *sqlGoInsert) SQLInsert(table sqlGoTable, columns sqlGoInsertColumnSlice, values ...sqlGoInsertValueSlice) SQLGoInsert {
	s.SetSQLInsert(table)
	s.SetSQLInsertColumn(columns...)
	s.SetSQLInsertValueSlice(values...)
	return s
}

func (s *sqlGoInsert) SetSQLInsert(table sqlGoTable) SQLGoInsert {
	s.table = table
	return s
}

func (s *sqlGoInsert) SetSQLInsertColumn(columns ...sqlGoInsertColumn) SQLGoInsert {
	s.columns = append(s.columns, columns...)
	return s
}

func (s *sqlGoInsert) SetSQLInsertValueSlice(values ...sqlGoInsertValueSlice) SQLGoInsert {
	s.values = append(s.values, values...)
	return s
}

func (s *sqlGoInsert) SetSQLInsertValue(values ...sqlGoInsertValue) SQLGoInsert {
	s.SetSQLInsertValueSlice(values)
	return s
}

func (s *sqlGoInsert) SetSQLGoSchema(sqlGoSchema SQLGoSchema) SQLGoInsert {
	s.sqlGoSchema = sqlGoSchema
	return s
}

func (s *sqlGoInsert) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoInsert {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoInsert) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGoInsert) BuildSQL() string {
	var sql string
	if len(s.columns) < 1 {
		return sql
	}

	sql = fmt.Sprintf("INSERT INTO %s%s (", s.sqlGoSchema.BuildSQL(), s.table)
	for i, v := range s.columns {
		if i > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}
		sql = fmt.Sprintf("%s%s", sql, v)
	}
	sql = fmt.Sprintf("%s)", sql)

	if len(s.values) < 1 {
		return sql
	}
	sql = fmt.Sprintf("%s VALUES ", sql)
	for iValues, vValues := range s.values {
		if iValues > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}

		sql = fmt.Sprintf("%s(", sql)
		for iValue, vValue := range vValues {
			if iValue > 0 {
				sql = fmt.Sprintf("%s, ", sql)
			}

			switch vType := vValue.(type) {
			case sqlGoInsertToTsVector:
				s.sqlGoParameter.SetSQLParameter(vType.value)
				sql = fmt.Sprintf("%sto_tsvector('%s', %s)", sql, vType.lang, s.sqlGoParameter.GetSQLParameterSign(vType.value))
			default:
				s.sqlGoParameter.SetSQLParameter(vType)
				sql = fmt.Sprintf("%s%s", sql, s.sqlGoParameter.GetSQLParameterSign(vType))
			}
		}
		sql = fmt.Sprintf("%s)", sql)
	}
	return sql
}
