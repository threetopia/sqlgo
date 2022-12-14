package sqlgo

import "fmt"

type SQLGoValues interface {
	SQLValues(alias sqlGoAlias, columns sqlGoValuesColumnSlice, values ...sqlGoValuesValueSlice) SQLGoValues
	SetSQLValues(alias sqlGoAlias) SQLGoValues
	SetSQLValuesColumn(columns ...sqlGoValuesColumn) SQLGoValues
	SetSQLValuesValue(values ...sqlGoValuesValueSlice) SQLGoValues

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoValues
	SQLGoMandatory
}

type (
	sqlGoValues struct {
		alias          sqlGoAlias
		columns        sqlGoValuesColumnSlice
		values         []sqlGoValuesValueSlice
		sqlGoParameter SQLGoParameter
	}

	sqlGoValuesColumn      string
	sqlGoValuesColumnSlice []sqlGoValuesColumn
	sqlGoValuesValue       interface{}
	sqlGoValuesValueSlice  []sqlGoValuesValue
)

func NewSQLGoValues() SQLGoValues {
	return &sqlGoValues{}
}

func SetSQLValuesColumn(columns ...sqlGoValuesColumn) sqlGoValuesColumnSlice {
	return columns
}

func SetSQLValuesValue(values ...sqlGoValuesValue) sqlGoValuesValueSlice {
	return values
}

func (s *sqlGoValues) SQLValues(alias sqlGoAlias, columns sqlGoValuesColumnSlice, values ...sqlGoValuesValueSlice) SQLGoValues {
	s.alias = alias
	s.SetSQLValuesColumn(columns...)
	s.SetSQLValuesValue(values...)
	return s
}

func (s *sqlGoValues) SetSQLValues(alias sqlGoAlias) SQLGoValues {
	s.alias = alias
	return s
}

func (s *sqlGoValues) SetSQLValuesColumn(columns ...sqlGoValuesColumn) SQLGoValues {
	s.columns = append(s.columns, columns...)
	return s
}

func (s *sqlGoValues) SetSQLValuesValue(values ...sqlGoValuesValueSlice) SQLGoValues {
	s.values = append(s.values, values...)
	return s
}

func (s *sqlGoValues) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoValues {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoValues) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGoValues) BuildSQL() string {
	var sql string
	if len(s.columns) < 1 {
		return sql
	}

	sql = "VALUES "
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
	sql = fmt.Sprintf("%s) %s(", sql, s.alias)
	for i, v := range s.columns {
		if i > 0 {
			sql = fmt.Sprintf("%s, ", sql)
		}
		sql = fmt.Sprintf("%s%s", sql, v)
	}

	sql = fmt.Sprintf("%s", sql)
	return sql
}
