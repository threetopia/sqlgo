package sqlgo

type SQLGoInsert interface {
	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoInsert
	SQLGoMandatory
}

type sqlGoInsert struct {
	sqlGoParameter SQLGoParameter
}

func NewSQLGoInsert() SQLGoInsert {
	return &sqlGoInsert{}
}

func (s *sqlGoInsert) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGoInsert {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGoInsert) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGoInsert) BuildSQL() string {
	return ""
}
