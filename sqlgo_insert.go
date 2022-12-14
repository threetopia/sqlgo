package sqlgo

import "fmt"

type SQLGoInsert interface {
	SQLInsert(table sqlGoTable, columns sqlGoInsertColumnList, values ...sqlGoInsertValueSlice) SQLGoInsert
	SetSQLInsert(table sqlGoTable) SQLGoInsert
	SetSQLInsertColumn(columns ...sqlGoInsertColumn) SQLGoInsert
	SetSQLInsertValue(values ...sqlGoInsertValueSlice) SQLGoInsert

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoInsert
	SQLGoMandatory
}

type (
	sqlGoInsert struct {
		table          sqlGoTable
		columns        []sqlGoInsertColumn
		values         []sqlGoInsertValueSlice
		sqlGoParameter SQLGoParameter
	}

	sqlGoInsertColumn     string
	sqlGoInsertColumnList []sqlGoInsertColumn
	sqlGoInsertValue      interface{}
	sqlGoInsertValueSlice []sqlGoInsertValue
)

func NewSQLGoInsert() SQLGoInsert {
	return &sqlGoInsert{}
}

func SetSQLInsertColumn(columns ...sqlGoInsertColumn) sqlGoInsertColumnList {
	return columns
}

func SetSQLInsertValue(values ...sqlGoInsertValue) sqlGoInsertValueSlice {
	return values
}

func (s *sqlGoInsert) SQLInsert(table sqlGoTable, columns sqlGoInsertColumnList, values ...sqlGoInsertValueSlice) SQLGoInsert {
	s.setSQLInsertTable(table)
	s.SetSQLInsertColumn(columns...)
	s.SetSQLInsertValue(values...)
	return s
}

func (s *sqlGoInsert) SetSQLInsert(table sqlGoTable) SQLGoInsert {
	s.setSQLInsertTable(table)
	return s
}

func (s *sqlGoInsert) setSQLInsertTable(table sqlGoTable) SQLGoInsert {
	s.table = table
	return s
}

func (s *sqlGoInsert) SetSQLInsertColumn(columns ...sqlGoInsertColumn) SQLGoInsert {
	s.columns = append(s.columns, columns...)
	return s
}

func (s *sqlGoInsert) SetSQLInsertValue(values ...sqlGoInsertValueSlice) SQLGoInsert {
	s.values = append(s.values, values...)
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
	if len(s.columns) < 1 {
		return ""
	}
	sql := fmt.Sprintf("INSERT INTO %s (", s.table)
	for i, v := range s.columns {
		if i > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}
		sql = fmt.Sprintf("%s%s", sql, v)
	}

	sql = fmt.Sprintf("%s)", sql)
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
			s.sqlGoParameter.SetSQLParameter(vValue)
			sql = fmt.Sprintf("%s%s", sql, s.sqlGoParameter.GetSQLParameterSign(vValue))
		}
		sql = fmt.Sprintf("%s)", sql)
	}
	return sql
}
