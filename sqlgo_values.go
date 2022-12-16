package sqlgo

import "fmt"

type SQLGoValues interface {
	SQLValues(values ...sqlGoValuesValueSlice) SQLGoValues
	SetSQLValuesValue(values ...sqlGoValuesValueSlice) SQLGoValues

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoValues
	SQLGoMandatory
}

type (
	sqlGoValues struct {
		values         []sqlGoValuesValueSlice
		sqlGoParameter SQLGoParameter
	}

	sqlGoValuesValue      interface{}
	sqlGoValuesValueSlice []sqlGoValuesValue
)

func NewSQLGoValues() SQLGoValues {
	return &sqlGoValues{}
}

func SetSQLValuesValue(values ...sqlGoValuesValue) sqlGoValuesValueSlice {
	return values
}

func (s *sqlGoValues) SQLValues(values ...sqlGoValuesValueSlice) SQLGoValues {
	s.SetSQLValuesValue(values...)
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
	if len(s.values) < 1 {
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

	return sql
}
