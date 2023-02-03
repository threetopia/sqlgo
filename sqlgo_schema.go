package sqlgo

import "fmt"

type SQLGoSchema interface {
	SQLSchema(schema string) SQLGoSchema
	SetSQLSchema(schema string) SQLGoSchema

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoSchema
	SQLGoMandatory
}

type sqlGoSchema struct {
	schema         string
	sqlGoParameter SQLGoParameter
}

func NewSQLGoSchema() SQLGoSchema {
	return new(sqlGoSchema)
}

func (s *sqlGoSchema) SetSQLSchema(schema string) SQLGoSchema {
	return s.SQLSchema(schema)
}

func (s *sqlGoSchema) SQLSchema(schema string) SQLGoSchema {
	if s.schema == "" {
		s.schema = schema
	}
	return s
}

func (s *sqlGoSchema) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoSchema {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoSchema) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGoSchema) BuildSQL() string {
	var sql string
	if s.schema == "" {
		return sql
	}

	sql = fmt.Sprintf("%s.", s.schema)
	return sql
}
